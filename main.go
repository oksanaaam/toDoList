package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ToDo struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var todos = []ToDo{
	{Id: "1", Title: "Make dinner", Status: "Not ready"},
	{Id: "2", Title: "Water plant", Status: "Done"},
	{Id: "3", Title: "Create Gin service", Status: "In progress"},
}

func main() {
	router := gin.Default()
	router.GET("/todos", getToDos)
	router.GET("/todos/:id", getToDosById)
	router.POST("/todos", postToDos)
	router.PUT("/todos/:id", updateToDos)
	router.DELETE("/todos/:id", deleteToDosById)

	router.Run("localhost:8080")
}

func getToDos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

func getToDosById(c *gin.Context) {
	id := c.Param("id")
	for _, todo := range todos {
		if todo.Id == id {
			c.IndentedJSON(http.StatusOK, todo)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "this id is not found"})
}

func postToDos(c *gin.Context) {
	var newTodo ToDo
	if err := c.BindJSON(&newTodo); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid data", "error": err.Error()})
		return
	}
	todos = append(todos, newTodo)
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "ToDo item was added to list"})
}

func updateToDos(c *gin.Context) {
	id := c.Param("id")
	var updatedToDo ToDo
	if err := c.BindJSON(&updatedToDo); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid data", "error": err.Error()})
	}
	for i, todo := range todos {
		if todo.Id == id {
			todos[i].Title = updatedToDo.Title
			todos[i].Status = updatedToDo.Status

			c.IndentedJSON(http.StatusOK, todos[i])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "this id is not found"})
}

func deleteToDosById(c *gin.Context) {
	id := c.Param("id")
	for i, todo := range todos {
		if todo.Id == id {
			todos = append(todos[:i], todos[i+1])
			c.IndentedJSON(http.StatusOK, gin.H{"message": "ToDo item deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "this id is not found"})
}
