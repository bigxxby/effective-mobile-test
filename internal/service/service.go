package service

import (
	"log"
	"time"

	"github.com/bigxxby/effective-mobile-test/internal/models"
	"github.com/bigxxby/effective-mobile-test/internal/repository"
)

type Service struct {
	Repository *repository.Repository
}

func New(repository repository.Repository) Service {
	return Service{
		Repository: &repository,
	}
}
func (s *Service) GetUsers(filter models.Filter, pagination models.Pagination, sortBy, sortOrder string) ([]models.User, error) {
	// Validate sortBy parameter
	validSortColumns := map[string]bool{
		"id":              true,
		"passport_number": true,
		"surname":         true,
		"name":            true,
	}

	if !validSortColumns[sortBy] {
		sortBy = "id"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	if pagination.Page <= 0 {
		pagination.Page = 1
	}

	if pagination.PageSize <= 0 {
		pagination.PageSize = 10
	}

	users, err := s.Repository.GetUsers(filter, pagination, sortBy, sortOrder)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return users, nil
}
func (s *Service) GetUserWorkloadsByUserID(userID int, startDate, endDate time.Time) ([]models.UserWorkload, error) {
	// Validate startDate and endDate
	if startDate.After(endDate) {
		return nil, models.ErrStartDateAfterEndDate
	}
	if time.Since(startDate) < 0 {
		return nil, models.ErrStartDateInFuture
	}

	userWorkloads, err := s.Repository.GetUserWorkloadsByUserID(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return userWorkloads, nil
}
