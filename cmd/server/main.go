package main

import (
	"log"
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

	todoService := service.NewTodoService(db)

	router := gin.Default()

	router.GET("/todos", handler.GetToDos(todoService))
	router.GET("/todos/:id", handler.GetToDosById(todoService))
	router.POST("/todos", handler.PostToDos(todoService))
	router.PUT("/todos/:id", handler.UpdateToDos(todoService))
	router.DELETE("/todos/:id", handler.DeleteToDosById(todoService))

	router.Run(cfg.ServerAddress)
}
