package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"toDoList/internal/handler"
	"toDoList/internal/loadbalancer"
	"toDoList/internal/service"
	"toDoList/internal/storage"
	"toDoList/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	var store storage.Storage
	var err error

	if cfg.DBType == "postgres" {
		store, err = storage.NewPostgresConnection(cfg.DBConnectionString)
	} else if cfg.DBType == "mongo" {
		store, err = storage.NewMongoConnection(cfg.MongoURI, cfg.MongoDBName, cfg.MongoCollectionName)
	} else {
		log.Fatalf("Unsupported DB type: %v", cfg.DBType)
	}

	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	defer store.Close()

	notificationChannel := make(chan string)

	todoService := service.NewTodoService(store)
	reminderService := service.NewReminderService(notificationChannel)

	// Launching reminder worker
	reminderService.StartWorker()

	go func() {
		for msg := range notificationChannel {
			log.Println("Notification: " + msg)
		}
	}()

	// List of servers to which we will send requests
	servers := []string{
		"localhost:8080",
		"localhost:8081",
		"localhost:8082",
	}

	// Create a load balancer
	lb := loadbalancer.NewLoadBalancer(servers)

	// Configure an HTTP server for the load balancer
	go func() {
		if isPortAvailable(":8085") {
			log.Println("Load Balancer is running on :8085")
			if err := http.ListenAndServe(":8085", lb); err != nil {
				log.Fatalf("Load Balancer failed: %v", err)
			}
		} else {
			log.Println("Port 8085 is already in use. Load balancer will not start.")
		}
	}()

	router := gin.Default()
	router.Use(handler.MaxConnections(150)) // limit the number of connections
	router.Use(handler.RateLimiter())       // limit the number of requests

	router.GET("/", handler.HomePage(todoService))
	router.GET("/todos", handler.GetToDos(todoService))
	router.GET("/todos/:id", handler.GetToDosById(todoService))
	router.GET("/todos/:id/image", handler.GetTodosImageById(todoService))
	router.POST("/todos", handler.PostToDos(todoService, reminderService))
	router.POST("/todos/:id/image", handler.UploadToDoImage(todoService))
	router.PUT("/todos/:id", handler.UpdateToDos(todoService))
	router.DELETE("/todos/:id", handler.DeleteToDosById(todoService))

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	// Step to handle termination signals (for graceful shutdown)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Start the main API server in a separate goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s\n", err)
		}
	}()
	log.Println("API Server is running on", cfg.ServerAddress)

	<-sigs
	log.Println("Shutdown signal received, starting graceful shutdown...")

	reminderService.StopWorker()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	<-ctx.Done()
	log.Println("Timeout of 3 seconds reached.")
	log.Println("Server gracefully stopped.")
}

// Function to check if the port is available
func isPortAvailable(port string) bool {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		return false
	}
	listen.Close()
	return true
}
