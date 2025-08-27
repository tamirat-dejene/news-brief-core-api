package usecase

import (
	"context"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
)

type sourceUsecase struct {
	sourceRepo contract.ISourceRepository
}

// NewSourceUsecase creates a new instance of sourceUsecase.
func NewSourceUsecase(sourceRepo contract.ISourceRepository) contract.ISourceUsecase {
	return &sourceUsecase{
		sourceRepo: sourceRepo,
	}
}

// GetAll fetches all source definitions from the repository.
func (uc *sourceUsecase) GetAll(ctx context.Context) ([]entity.Source, error) {
	// For now, the use case is a simple pass-through.
	// Business logic like filtering out inactive sources could be added here later.
	return uc.sourceRepo.GetAll(ctx)
}
