package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ToDo struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var todos = []ToDo{
	{Id: 1, Title: "Make dinner", Status: "Not ready"},
	{Id: 2, Title: "Water plant", Status: "Done"},
	{Id: 3, Title: "Create Gin service", Status: "In progress"},
}

func main() {
	router := gin.Default()
	router.GET("/todos", getToDos)

	router.Run("localhost:8080")
}

func getToDos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}
