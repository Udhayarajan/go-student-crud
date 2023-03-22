package database

import (
	"fmt"
	"log"
	"student-crud/models"
)

func Insert(student models.Student) *models.Student {
	err := DB.QueryRow("INSERT INTO students (name, roll_number) VALUES ($1, $2) RETURNING id, name, roll_number", student.Name, student.RollNumber).Scan(&student.Id, &student.Name, &student.RollNumber)
	if err != nil {
		return nil
	}
	return &student
}

func UpdateByRollNumber(rollNumber string, student models.Student) *models.Student {
	result, err := DB.Exec("UPDATE students set name=$1 where roll_number=$2", student.Name, rollNumber)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	m := GetStudentByRollNumber(rollNumber)
	if m == nil {
		return nil
	}
	student = *m
	affected, err := result.RowsAffected()
	if affected != 1 {
		log.Fatalf("expected to affect 1 row, but affected %d", affected)
		return nil
	}
	return &student
}

func DeleteByRollNumber(rollNumber string) (bool, *models.Student) {
	stmt, err := DB.Prepare("DELETE FROM students WHERE roll_number=$1")
	student := GetStudentByRollNumber(rollNumber)
	if err != nil || student == nil {
		return false, nil
	}
	result, err := stmt.Exec(rollNumber)
	if err != nil {
		return false, nil
	}
	affected, err := result.RowsAffected()
	if affected != 1 {
		log.Fatalf("expected to affect 1 row, but affected %d", affected)
		return false, nil
	}
	return true, student
}

func GetStudentByRollNumber(rollNumber string) *models.Student {
	stmt, err := DB.Prepare("SELECT * FROM students where roll_number=$1")
	defer stmt.Close()
	query, err := stmt.Query(rollNumber)
	defer query.Close()
	if err != nil {
		return nil
	}
	if query.Next() {
		var student models.Student
		query.Scan(&student.Id, &student.Name, &student.RollNumber)
		return &student
	}
	return nil
}

func GetAllStudents() []models.Student {
	var students []models.Student
	query, err := DB.Query("SELECT * FROM students")
	defer query.Close()
	if err != nil {
		return students
	}
	for query.Next() {
		var student models.Student
		query.Scan(&student.Id, &student.Name, &student.RollNumber)
		fmt.Println(student)
		students = append(students, student)
	}
	return students
}
