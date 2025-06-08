package service

import (
	"ServiceApi/internal/repository"
	"ServiceApi/internal/model"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SearchServices(filters map[string]string) ([]model.Service, error) {
	return s.repo.FindServices(filters)
}