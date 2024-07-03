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

// Получение всех пользователей
func (s *Service) GetUsers(filter models.Filter, pagination models.Pagination, sortBy, sortOrder string) ([]models.User, error) {
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

// Получение рабочей нагрузки пользователей по их ID
func (s *Service) GetUserWorkloadsByUserID(userID int, startDate, endDate time.Time) ([]models.UserWorkload, error) {
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

// Запуск задачи
func (s *Service) StartTask(userID, taskID int) error {

	inPorgress, _ := s.Repository.IsTaskInProgress(userID, taskID)
	if inPorgress {
		return models.ErrTaskAlreadyStarted
	}
	isTaskExists, _ := s.Repository.IsTaskExists(taskID)
	if !isTaskExists {
		return models.ErrTaskNotFound
	}

	err := s.Repository.StartTask(userID, taskID)
	if err != nil {
		return err
	}

	return nil
}

// Завершение задачи
func (s *Service) EndTask(userID, taskID int) error {
	inPorgress, _ := s.Repository.IsTaskInProgress(userID, taskID)
	if !inPorgress {
		return models.ErrTaskNotStarted
	}
	isTaskExists, _ := s.Repository.IsTaskExists(taskID)
	if !isTaskExists {
		return models.ErrTaskNotFound
	}
	err := s.Repository.EndTask(userID, taskID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUser(userID int) (models.User, error) {
	user, err := s.Repository.GetUser(userID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Получение всех задач
func (s *Service) GetTasks() ([]models.Task, error) {
	tasks, err := s.Repository.GetTasks()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
func (s *Service) CreateUser(user models.UserData) (int, error) {
	userExists, _ := s.Repository.UserExistsByPassportNumber(user.PassportNumber)
	if userExists {
		return 0, models.ErrUserAlreadyExists
	}

	userID, err := s.Repository.CreateUser(user)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *Service) UpdateUser(userID int, user models.User) error {
	err := s.Repository.UpdateUser(userID, user)
	if err != nil {
		return err
	}

	return nil
}
func (s *Service) DeleteUser(userID int) error {
	err := s.Repository.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}
