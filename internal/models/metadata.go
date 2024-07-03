package models

import "errors"

type Filter struct {
	PassportNumber string
	Surname        string
	Name           string
}
type Pagination struct {
	Page     int
	PageSize int
}

var (
	ErrInvalidID = errors.New("invalid ID")
)
