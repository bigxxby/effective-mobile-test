package service

import "github.com/bigxxby/effective-mobile-test/internal/repository"

type Service struct {
	Repository *repository.Repository
}

func New(repository repository.Repository) Service {
	return Service{
		Repository: &repository,
	}
}
func (s *Service) GetUsers() ([]string, error) {
	return s.Repository.GetUsers()
}
