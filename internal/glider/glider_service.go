package glider

import (
	"context"

	"github.com/AurelienS/cigare/internal/storage"
)

// GliderService provides operations for managing gliders.
type GliderService struct {
	Repo GliderRepository
}

// NewGliderService creates a new instance of GliderService with the necessary dependencies.
func NewGliderService(repository GliderRepository) *GliderService {
	return &GliderService{
		Repo: repository,
	}
}

// GetGliders retrieves gliders accessible to the specified user.
func (g *GliderService) GetGliders(ctx context.Context, user storage.User) ([]storage.Glider, error) {
	gliders, err := g.Repo.GetGliders(ctx, user)
	if err != nil {
		return nil, err // Consider wrapping this error to add more context.
	}
	return gliders, nil
}
