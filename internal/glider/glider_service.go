package glider

import (
	"context"
	"fmt"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/web/template/flight"
)

type GliderService struct {
	repo GliderRepository
}

func NewGliderService(repository GliderRepository) GliderService {
	return GliderService{
		repo: repository,
	}
}

func (g *GliderService) GetGliders(ctx context.Context, user storage.User) ([]flight.GliderView, error) {
	gliders, err := g.repo.GetGliders(ctx, user)

	view := []flight.GliderView{}

	for _, glider := range gliders {

		isSelected := false
		if user.DefaultGliderID.Valid {
			if glider.ID == user.DefaultGliderID.Int32 {
				isSelected = true
			}
		}
		linkToUpdate := fmt.Sprintf("/user/%d?defaultGliderId=%d", user.ID, glider.ID)
		id := int(glider.ID)

		view = append(view, flight.GliderView{
			Name:         glider.Name,
			LinkToUpdate: linkToUpdate,
			IsSelected:   isSelected,
			ID:           id,
		})
	}

	return view, err
}

func (g *GliderService) AddGlider(ctx context.Context, gliderName string, user storage.User) error {
	return g.repo.AddGlider(ctx, gliderName, user)
}
