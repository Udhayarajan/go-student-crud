package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"net/http"
	"regexp"
	"student-crud/database"
	"student-crud/models"
)

var validate *validator.Validate

var rollNumberRegex = regexp.MustCompile("^[0-9]{2}[A-Za-z][0-9]{3}")

var trans ut.Translator

func InitValidator() {
	validate = validator.New()

	_ = validate.RegisterValidation("roll", func(fl validator.FieldLevel) bool {
		return rollNumberRegex.MatchString(fl.Field().String())
	})

	translator := en.New()
	trans, _ = ut.New(translator, translator).GetTranslator("en")

	validate.RegisterTranslation("roll", trans, func(ut ut.Translator) error {
		return ut.Add("roll", "The %s field must be a valid roll number like 20PXXX", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("roll", fe.Field())
		fmt.Println(t)
		return fmt.Sprintf(t, fe.Field())
	})

	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "The %s field is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return fmt.Sprintf(t, fe.Field())
	})

}

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
	if !has {
		sendJson(context.Writer, http.StatusBadRequest, toJson(gin.H{
			"error": "Invalid query. The query must contain 'RollNumber' which needed to be deleted",
		}))
		return
	}
	if !rollNumberRegex.MatchString(rollNum) {
		sendJson(context.Writer, http.StatusBadRequest, toJson(gin.H{
			"error": "Roll Number must be a valid roll number like 20PXXX",
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
	if err != nil {
		fmt.Println(err)
		sendJson(context.Writer, http.StatusBadRequest, toJson(gin.H{
			"error": "Please send the body for the request, which must contains 'Name' and 'RollNumber'",
		}))
		return nil
	}
	err = validate.Struct(student)
	var validationResult []string
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			sendJson(context.Writer, http.StatusBadRequest, toJson(gin.H{
				"error": "Please send the body for the request, which must contains 'Name' and 'RollNumber'",
			}))
			return nil
		}

		for _, err := range err.(validator.ValidationErrors) {
			validationResult = append(validationResult, fmt.Sprintf("`%s`: %s", err.Field(), err.Translate(trans)))
		}
		var stringArray []string
		_ = json.Unmarshal(toJson(validationResult), &stringArray)
		sendJson(context.Writer, http.StatusBadRequest, toJson(gin.H{
			"errors": stringArray,
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
