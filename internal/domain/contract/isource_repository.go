package contract

import (
	"context"
	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
)

// ISourceRepository defines the persistence contract for Source entities.
type ISourceRepository interface {
	CreateSource(ctx context.Context, source *entity.Source) error
	// CheckSlugExists checks if a source with the given slug exists.
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
	// GetBySlug retrieves a single source by its unique slug.
	GetBySlug(ctx context.Context, slug string) (*entity.Source, error)
	CheckURLExists(ctx context.Context, url string) (bool, error)

	// GetAll retrieves all sources from the collection.
	GetAll(ctx context.Context) ([]entity.Source, error)
}
