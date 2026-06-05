package service

import (
	"context"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

type ReferenceService interface {
	GetCategories(ctx context.Context) ([]model.Category, error)
	GetGenres(ctx context.Context) ([]model.Genre, error)
	GetAgeRatings(ctx context.Context) ([]model.AgeRating, error)
}

type referenceService struct {
	repo repository.ReferenceRepository
}

func NewReferenceService(repo repository.ReferenceRepository) ReferenceService {
	return &referenceService{repo: repo}
}

func (s *referenceService) GetCategories(ctx context.Context) ([]model.Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *referenceService) GetGenres(ctx context.Context) ([]model.Genre, error) {
	return s.repo.GetGenres(ctx)
}

func (s *referenceService) GetAgeRatings(ctx context.Context) ([]model.AgeRating, error) {
	return s.repo.GetAgeRatings(ctx)
}