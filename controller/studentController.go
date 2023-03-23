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
	sendJson(context.Writer, http.StatusOK, toJson(students))
}

func AddStudent(context *gin.Context) {
	student := getStudentDetails(context)
	if student == nil {
		return
	}
	student = database.Insert(*student)
	sendJson(context.Writer, http.StatusCreated, toJson(student))
}

func DeleteStudent(context *gin.Context) {
	var rollNum, has = context.GetQuery("RollNumber")
	if !has || rollNum == "" {
		sendJson(context.Writer, http.StatusBadRequest, toJson(gin.H{
			"error": "Invalid query. The query must contain 'RollNumber' which needed to be deleted",
		}))
		return
	}
	_, student := database.DeleteByRollNumber(rollNum)
	if student == nil {
		sendJson(context.Writer, http.StatusNotFound, toJson(gin.H{
			"error": "unable to find the record for the given RollNumber '" + rollNum + "'",
		}))
		return
	}
	sendJson(context.Writer, http.StatusOK, toJson(student))
}

func UpdateStudent(context *gin.Context) {
	student := getStudentDetails(context)
	if student == nil {
		return
	}
	updateStudent := database.UpdateByRollNumber(student.RollNumber, *student)
	if updateStudent == nil {
		sendJson(context.Writer, http.StatusNotFound, toJson(gin.H{
			"error": "unable to find the record for the given RollNumber '" + student.RollNumber + "'",
		}))
		return
	}
	sendJson(context.Writer, http.StatusOK, toJson(updateStudent))
}

func getStudentDetails(context *gin.Context) *models.Student {
	student := models.Student{}
	err := context.BindJSON(&student)
	if err != nil || student.Name == "" || student.RollNumber == "" {
		fmt.Println(err)
		sendJson(context.Writer, http.StatusBadRequest, toJson(gin.H{
			"error": "Please send the body for the request, which must contains 'Name' and 'RollNumber'",
		}))
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
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(json)
}
