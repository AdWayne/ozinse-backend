package service

import (
	"context"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

type FavoriteService interface {
	GetFavorites(ctx context.Context, userID int) ([]model.Project, error)
	AddFavorite(ctx context.Context, userID int, projectID int) error
	DeleteFavorite(ctx context.Context, userID int, projectID int) error
}

type favoriteService struct {
	favRepo     repository.FavoriteRepository
	projectRepo repository.ProjectRepository
}

func NewFavoriteService(favRepo repository.FavoriteRepository, projectRepo repository.ProjectRepository) FavoriteService {
	return &favoriteService{favRepo: favRepo, projectRepo: projectRepo}
}

func (s *favoriteService) GetFavorites(ctx context.Context, userID int) ([]model.Project, error) {
	return s.favRepo.GetFavorites(ctx, userID)
}

func (s *favoriteService) AddFavorite(ctx context.Context, userID int, projectID int) error {
	// Узнаем тип проекта (movie или series), чтобы правильно сохранить в базу
	pType, _, err := s.projectRepo.GetProjectTypeAndTable(ctx, projectID)
	if err != nil {
		return err 
	}

	return s.favRepo.AddFavorite(ctx, userID, projectID, pType)
}

func (s *favoriteService) DeleteFavorite(ctx context.Context, userID int, projectID int) error {
	pType, _, err := s.projectRepo.GetProjectTypeAndTable(ctx, projectID)
	if err != nil {
		return err
	}

	return s.favRepo.DeleteFavorite(ctx, userID, projectID, pType)
}