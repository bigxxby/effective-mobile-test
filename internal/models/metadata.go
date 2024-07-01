package models

type Filter struct {
	PassportNumber string
	Surname        string
	Name           string
}
type Pagination struct {
	Page     int
	PageSize int
}
