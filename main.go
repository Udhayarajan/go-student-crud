package main

import (
	"github.com/gin-gonic/gin"
	"student-crud/controller"
	"student-crud/database"
)

func init() {
	database.ConnectToDatabase()
}

func main() {
	r := gin.Default()
	r.GET("/", controller.GetStudents)

	r.POST("/add", controller.AddStudent)

	r.DELETE("/remove", controller.DeleteStudent)

	r.PUT("/update", controller.UpdateStudent)
	err := r.Run()
	if err != nil {
		return
	}
}
