package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"student-crud/database"
	"student-crud/models"
)

func GetStudents(context *gin.Context) {
	students := database.GetAllStudents()
	fmt.Println(students)
	sendJson(context.Writer, 200, toJson(students))
}

func AddStudent(context *gin.Context) {
	student := getStudentDetails(context)
	if student == nil {
		context.Status(400)
		return
	}
	student = database.Insert(*student)
	sendJson(context.Writer, http.StatusCreated, toJson(student))
}

func DeleteStudent(context *gin.Context) {
	var rollNum, _ = context.GetQuery("RollNumber")
	_, student := database.DeleteByRollNumber(rollNum)
	sendJson(context.Writer, http.StatusOK, toJson(student))
}

func UpdateStudent(context *gin.Context) {
	student := getStudentDetails(context)
	student = database.UpdateByRollNumber(student.RollNumber, *student)
	sendJson(context.Writer, http.StatusOK, toJson(student))
}

func getStudentDetails(context *gin.Context) *models.Student {
	student := models.Student{}
	err := context.BindJSON(&student)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &student
}
func toJson(obj any) []byte {
	marshal, _ := json.Marshal(obj)
	fmt.Println(marshal)
	return marshal
}

func sendJson(writer gin.ResponseWriter, code int, json []byte) {
	writer.WriteHeader(code)
	writer.Write(json)
}
