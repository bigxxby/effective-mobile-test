package models

import "errors"

type UserWorkload struct {
	TaskID       int
	TaskName     string
	TotalHours   int
	TotalMinutes int
}

var (
	ErrStartDateAfterEndDate = errors.New("start date is after end date")
	ErrStartDateInFuture     = errors.New("start date is in the future")
)
