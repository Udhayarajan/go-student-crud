package controller

import (
	"encoding/json"
	"fmt"
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

type H map[string]any

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

func GetStudents(writer http.ResponseWriter, request *http.Request) {
	students := database.GetAllStudents()
	fmt.Println(students)
	sendJson(writer, http.StatusOK, toJson(students))
}

func AddStudent(writer http.ResponseWriter, request *http.Request) {
	student := getStudentDetails(writer, request)
	if student == nil {
		return
	}
	student = database.Insert(*student)
	sendJson(writer, http.StatusCreated, toJson(student))
}

func DeleteStudent(writer http.ResponseWriter, request *http.Request) {
	var rollNum = request.URL.Query().Get("RollNumber")
	if rollNum == "" {
		sendJson(writer, http.StatusBadRequest, toJson(H{
			"error": "Invalid query. The query must contain 'RollNumber' which needed to be deleted",
		}))
		return
	}
	if !rollNumberRegex.MatchString(rollNum) {
		sendJson(writer, http.StatusBadRequest, toJson(H{
			"error": "Roll Number must be a valid roll number like 20PXXX",
		}))
		return
	}
	_, student := database.DeleteByRollNumber(rollNum)
	if student == nil {
		sendJson(writer, http.StatusNotFound, toJson(H{
			"error": "unable to find the record for the given RollNumber '" + rollNum + "'",
		}))
		return
	}
	sendJson(writer, http.StatusOK, toJson(student))
}

func UpdateStudent(writer http.ResponseWriter, request *http.Request) {
	student := getStudentDetails(writer, request)
	if student == nil {
		return
	}
	updateStudent := database.UpdateByRollNumber(student.RollNumber, *student)
	if updateStudent == nil {
		sendJson(writer, http.StatusNotFound, toJson(H{
			"error": "unable to find the record for the given RollNumber '" + student.RollNumber + "'",
		}))
		return
	}
	sendJson(writer, http.StatusOK, toJson(updateStudent))
}

func getStudentDetails(writer http.ResponseWriter, request *http.Request) *models.Student {
	student := models.Student{}
	err := json.NewDecoder(request.Body).Decode(&student)
	if err != nil {
		fmt.Println(err)
		sendJson(writer, http.StatusBadRequest, toJson(H{
			"error": "Please send the body for the request, which must contains 'Name' and 'RollNumber'",
		}))
		return nil
	}
	err = validate.Struct(student)
	var validationResult []string
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			sendJson(writer, http.StatusBadRequest, toJson(H{
				"error": "Please send the body for the request, which must contains 'Name' and 'RollNumber'",
			}))
			return nil
		}

		for _, err := range err.(validator.ValidationErrors) {
			validationResult = append(validationResult, fmt.Sprintf("`%s`: %s", err.Field(), err.Translate(trans)))
		}
		var stringArray []string
		_ = json.Unmarshal(toJson(validationResult), &stringArray)
		sendJson(writer, http.StatusBadRequest, toJson(H{
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

func sendJson(writer http.ResponseWriter, code int, json []byte) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(json)
}
