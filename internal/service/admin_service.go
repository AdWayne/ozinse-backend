package service

import (
	"context"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

type AdminService struct {
	repo *repository.AdminRepository
}

func NewAdminService(repo *repository.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) GetProjects(ctx context.Context, page, limit int, sortBy, sortOrder string) (model.ProjectsListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	items, total, err := s.repo.GetProjects(ctx, limit, offset, sortBy, sortOrder)
	if err != nil {
		return model.ProjectsListResponse{}, err
	}

	return model.ProjectsListResponse{
		Items:      items,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *AdminService) CreateProject(ctx context.Context, p *model.ProjectAdmin) error {
	return s.repo.CreateProject(ctx, p)
}

func (s *AdminService) UpdateProject(ctx context.Context, p *model.ProjectAdmin) error {
	return s.repo.UpdateProject(ctx, p)
}

func (s *AdminService) DeleteProject(ctx context.Context, id int, projType string) error {
	return s.repo.DeleteProject(ctx, id, projType)
}

func (s *AdminService) CreateSeason(ctx context.Context, seriesID int, season *model.SeasonAdmin) error {
	return s.repo.CreateSeason(ctx, seriesID, season)
}

func (s *AdminService) CreateEpisode(ctx context.Context, seasonID int, e *model.EpisodeAdmin) error {
	return s.repo.CreateEpisode(ctx, seasonID, e)
}

func (s *AdminService) UpdateFeaturedOrder(ctx context.Context, blockType string, items []model.FeaturedOrderItem) error {
	return s.repo.UpdateFeaturedOrder(ctx, blockType, items)
}

func (s *AdminService) GetUsers(ctx context.Context) ([]model.UserResponse, error) {
	return s.repo.GetUsers(ctx)
}

func (s *AdminService) GetRoles(ctx context.Context) ([]model.Role, error) {
	return s.repo.GetRoles(ctx)
}

func (s *AdminService) CreateRole(ctx context.Context, role *model.Role) error {
	return s.repo.CreateRole(ctx, role)
}

func (s *AdminService) UpdateRole(ctx context.Context, role *model.Role) error {
	return s.repo.UpdateRole(ctx, role)
}

func (s *AdminService) DeleteRole(ctx context.Context, id int) error {
	return s.repo.DeleteRole(ctx, id)
}

func (s *AdminService) AssignRole(ctx context.Context, userID int, roleID int) error {
	return s.repo.AssignRole(ctx, userID, roleID)
}

func (s *AdminService) CreateCategory(ctx context.Context, cat *model.Category) error {
	return s.repo.CreateCategory(ctx, cat)
}

func (s *AdminService) UpdateCategory(ctx context.Context, cat *model.Category) error {
	return s.repo.UpdateCategory(ctx, cat)
}

func (s *AdminService) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *AdminService) CreateGenre(ctx context.Context, g *model.Genre) error {
	return s.repo.CreateGenre(ctx, g)
}

func (s *AdminService) UpdateGenre(ctx context.Context, g *model.Genre) error {
	return s.repo.UpdateGenre(ctx, g)
}

func (s *AdminService) DeleteGenre(ctx context.Context, id int) error {
	return s.repo.DeleteGenre(ctx, id)
}

func (s *AdminService) CreateAgeRating(ctx context.Context, ar *model.AgeRating) error {
	return s.repo.CreateAgeRating(ctx, ar)
}

func (s *AdminService) UpdateAgeRating(ctx context.Context, ar *model.AgeRating) error {
	return s.repo.UpdateAgeRating(ctx, ar)
}

func (s *AdminService) DeleteAgeRating(ctx context.Context, id int) error {
	return s.repo.DeleteAgeRating(ctx, id)
}