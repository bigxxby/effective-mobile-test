package models

import "errors"

type User struct {
	ID             int    `json:"id"`
	Surname        string `json:"surname" binding:"required"`
	Name           string `json:"name" binding:"required"`
	PassportNumber string `json:"passport_number"`
}

type UserData struct {
	PassportNumber string `json:"passport_number" binding:"required"`
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)
