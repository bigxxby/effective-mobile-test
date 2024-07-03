package models

// ResponseUsersList описывает структуру ответа на запрос получения пользователей.
type ResponseUsersList struct {
	Users []User `json:"users"`
}

// ErrorResponse описывает структуру ответа на ошибку.
type ErrorResponse struct {
	Error string `json:"error"`
}

type OKresponse struct {
	Message string `json:"message"`
}
type UserUpdate struct {
	Surname string `json:"surname" binding:"required"`
	Name    string `json:"name" binding:"required"`
}
type ResponseTasksList struct {
	Tasks []Task `json:"tasks"`
}
