package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"toDoList/internal/handler"
	"toDoList/internal/service"
	"toDoList/internal/storage"
	"toDoList/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := storage.NewPostgresConnection(cfg.DBConnectionString)
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	defer db.Close()

	todoService := service.NewTodoService(db)

	router := gin.Default()

	router.GET("/", handler.HomePage(todoService))
	router.GET("/todos", handler.GetToDos(todoService))
	router.GET("/todos/:id", handler.GetToDosById(todoService))
	router.POST("/todos", handler.PostToDos(todoService))
	router.PUT("/todos/:id", handler.UpdateToDos(todoService))
	router.DELETE("/todos/:id", handler.DeleteToDosById(todoService))

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s\n", err)
		}
	}()
	log.Println("Server is running on", cfg.ServerAddress)

	<-sigs
	log.Println("Shutdown signal received, starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	<-ctx.Done()
	log.Println("Timeout of 3 seconds reached.")

	log.Println("Server gracefully stopped.")
}
