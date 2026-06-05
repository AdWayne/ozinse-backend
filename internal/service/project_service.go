package service

import (
	"context"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

type ProjectService interface {
	GetAll(ctx context.Context, filter model.ProjectFilter) ([]model.Project, error)
	GetFeatured(ctx context.Context) (map[string][]model.Project, error)
	GetByID(ctx context.Context, id int) (*model.ProjectDetail, error)
	GetSeasons(ctx context.Context, seriesID int) ([]model.Season, error)
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

// 1. ПОЛУЧЕНИЕ ВСЕХ ПРОЕКТОВ С ФИЛЬТРАЦИЕЙ
func (s *projectService) GetAll(ctx context.Context, filter model.ProjectFilter) ([]model.Project, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Size < 1 || filter.Size > 100 {
		filter.Size = 10
	}
	
	return s.repo.GetProjects(ctx, filter)
}

// 2. ПОДБОРКИ ДЛЯ ГЛАВНОГО ЭКРАНА (Трендтегілер, Жаңа жобалар)
func (s *projectService) GetFeatured(ctx context.Context) (map[string][]model.Project, error) {
	featured := make(map[string][]model.Project)

	// Подборка 1: "Жаңа жобалар" (Просто последние добавленные 6 проектов)
	newProjects, err := s.repo.GetProjects(ctx, model.ProjectFilter{Page: 1, Size: 6})
	if err != nil {
		return nil, err
	}
	featured["Жаңа жобалар"] = newProjects

	trendingProjects, err := s.repo.GetProjects(ctx, model.ProjectFilter{Page: 1, Size: 6, AgeRatingID: 2}) 
	if err != nil {
		// Если тренды упали, отдаем хотя бы новые, чтобы главный экран мобилки не ломался
		featured["Трендтегілер"] = []model.Project{}
	} else {
		featured["Tremdtegiler"] = trendingProjects
	}

	return featured, nil
}

// 3. ДЕТАЛЬНАЯ ИНФОРМАЦИЯ О ПРОЕКТЕ
func (s *projectService) GetByID(ctx context.Context, id int) (*model.ProjectDetail, error) {
	// Сначала узнаем через базу: это фильм или сериал, и в какой таблице он лежит
	_, tableName, err := s.repo.GetProjectTypeAndTable(ctx, id)
	if err != nil {
		return nil, err
	}

	// Вытаскиваем полные данные проекта
	return s.repo.GetProjectByID(ctx, id, tableName)
}

// 4. СПИСОК СЕЗОНОВ И СЕРИЙ ДЛЯ СЕРИАЛА
func (s *projectService) GetSeasons(ctx context.Context, seriesID int) ([]model.Season, error) {
	// Сначала проверим, существует ли этот сериал и действительно ли это сериал
	pType, _, err := s.repo.GetProjectTypeAndTable(ctx, seriesID)
	if err != nil || pType != "series" {
		return nil, err
	}

	return s.repo.GetSeasonsWithEpisodes(ctx, seriesID)
}