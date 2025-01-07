package handler

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

func SSENotificationHandler(notificationChannel chan string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set heasers for SSE
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Writer.Flush()

		// Sent notification to the client through SSE, if exist in the channel
		for msg := range notificationChannel {
			// Sent notification in format SSE
			fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
			c.Writer.Flush()

			// Close if client closed connection
			if c.Writer.CloseNotify() != nil {
				return
			}
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
		// Add the path to the image
		if todo.ImagePath != "" {
			todo.ImagePath = "/images/" + filepath.Base(todo.ImagePath)
		}
		c.JSON(http.StatusOK, todo)
	}
}

func GetTodosImageById(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		todo, err := todoService.GetTodoImageById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
			return
		}

		// return error if no image
		if todo.ImagePath == "" {
			c.JSON(http.StatusNotFound, gin.H{"message": "No image attached to this todo"})
			return
		}

		// checking if file is absent
		imagePath := todo.ImagePath
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Image not found"})
			return
		}

		// sent file to the user
		c.File(imagePath)
	}
}

func PostToDos(todoService service.TodoService, reminderService *service.ReminderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newTodo model.ToDo
		if err := c.BindJSON(&newTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect data", "error": err.Error()})
			return
		}

		if newTodo.ReminderTime != "" {
			duration, err := time.ParseDuration(newTodo.ReminderTime)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid time format", "error": err.Error()})
				return
			}

			reminderTime := time.Now().Add(duration)

			reminderService.AddReminder(service.Reminder{
				ID:           newTodo.ID,
				ReminderTime: reminderTime,
				TaskName:     newTodo.Title,
			})
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

func SaveImage(file *multipart.FileHeader) (string, error) {
	// update path to div with files
	dir := "./uploads/images"
	err := os.MkdirAll(dir, os.ModePerm) // create div if not exist
	if err != nil {
		log.Println("Error creating directory:", err)
		return "", err
	}
	log.Println("Directory created or exists:", dir)

	// generate unique file name
	filename := fmt.Sprintf("%s_%s", time.Now().Format("20060102150405"), file.Filename)
	filepath := filepath.Join(dir, filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// copy data to new file
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func UploadToDoImage(todoService service.TodoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		file, _ := c.FormFile("image")
		if file == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "No image uploaded"})
			return
		}

		// save file on server
		imagePath, err := SaveImage(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save image", "error": err.Error()})
			return
		}

		// update ToDo in db with new file path
		err = todoService.UpdateTodoImage(id, imagePath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully", "image_path": imagePath})
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
