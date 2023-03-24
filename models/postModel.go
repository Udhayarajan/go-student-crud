package models

import _ "regexp"

type Student struct {
	Id         int
	Name       string `validate:"required"`
	RollNumber string `validate:"required,roll"`
}
