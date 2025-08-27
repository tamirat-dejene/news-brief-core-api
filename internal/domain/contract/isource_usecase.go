package contract

import (
	"context"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
)

// ISourceUsecase defines the business logic contract for sources.
type ISourceUsecase interface {
	GetAll(ctx context.Context) ([]entity.Source, error)
}
