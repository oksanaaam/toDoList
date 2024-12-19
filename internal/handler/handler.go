package handler

import (
	"fmt"
	"net/http"
	"time"
	"toDoList/internal/model"
	"toDoList/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiter   = make(map[string]*rate.Limiter) // map limit for every user
	rateLimit = rate.Every(1 * time.Second)    // Limit: 1 request per second
	burst     = 5                              // max amount requsts
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// cheking if exist limit for this IP
		if _, exists := limiter[ip]; !exists {
			limiter[ip] = rate.NewLimiter(rateLimit, burst)
		}

		l := limiter[ip]

		// cheking if allow we do a request
		if !l.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func MaxConnections(limit int) gin.HandlerFunc {
	sem := make(chan struct{}, limit)
	release := func() { <-sem }
	return func(c *gin.Context) {
		select {
		case sem <- struct{}{}: // acquire before request
			defer release() // release after request
			c.Next()
		default:
			c.AbortWithError(http.StatusServiceUnavailable,
				fmt.Errorf("too many connections. limit %v", limit)) // send 503 and stop the chain
		}
	}
}

func HomePage(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome To To-Do Server")
	}
}

func GetToDos(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := todoService.GetAllTodos()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "There are no any todos", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todos)
	}
}

func GetToDosById(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		todo, err := todoService.GetTodoById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "todo not found"})
			return
		}
		c.JSON(http.StatusOK, todo)
	}
}

func PostToDos(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newTodo model.ToDo
		if err := c.BindJSON(&newTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect data", "error": err.Error()})
			return
		}
		err := todoService.AddTodo(newTodo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add todo", "error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "todo added"})
	}
}

func UpdateToDos(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var updatedTodo model.ToDo
		if err := c.BindJSON(&updatedTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect data", "error": err.Error()})
			return
		}
		err := todoService.UpdateTodo(id, updatedTodo)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "todo not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "todo updated"})
	}
}

func DeleteToDosById(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := todoService.DeleteTodo(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "todo not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "todo deleted"})
	}
}
