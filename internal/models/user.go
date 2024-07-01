package models

type User struct {
	ID             int    `json:"id"`
	PassportNumber string `json:"passport_number"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
}
