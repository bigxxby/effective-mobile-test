package models

import "errors"

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	ErrTaskNotStarted        = errors.New("start date is after end date")
	ErrTaskAlreadyStarted    = errors.New("task already started")
	ErrStartDateAfterEndDate = errors.New("start date is after end date")
	ErrStartDateInFuture     = errors.New("start date is in the future")
	ErrEndDateInFuture       = errors.New("end date is in the future")
	ErrTaskNotEnded          = errors.New("task not ended")
	ErrTaskNotFound          = errors.New("task not found")
)
