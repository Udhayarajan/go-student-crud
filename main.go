package main

import (
	"net/http"
	"student-crud/controller"
	"student-crud/database"
)

func init() {
	database.ConnectToDatabase()
}

func main() {
	controller.InitValidator()

	http.HandleFunc("/", controller.GetStudents)

	http.HandleFunc("/add", controller.AddStudent)

	http.HandleFunc("/remove", controller.DeleteStudent)

	http.HandleFunc("/update", controller.UpdateStudent)
	err := http.ListenAndServe("localhost:8080", nil)

	if err != nil {
		return
	}
}
