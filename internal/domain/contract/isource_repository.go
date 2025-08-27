package contract

import (
	"context"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
)

// ISourceRepository defines the persistence contract for Source entities.
type ISourceRepository interface {
	// Exists checks if a source with the given key exists.
	Exists(ctx context.Context, key string) (bool, error)

	// GetByKey retrieves a single source by its unique key.
	GetByKey(ctx context.Context, key string) (*entity.Source, error)

	// GetAll retrieves all sources from the collection.
	GetAll(ctx context.Context) ([]entity.Source, error)
}
