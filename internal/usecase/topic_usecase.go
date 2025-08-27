package usecase

import (
	"context"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
)

type topicUsecase struct {
	topicRepo contract.ITopicRepository
}

// NewTopicUsecase creates a new instance of topicUsecase.
func NewTopicUsecase(topicRepo contract.ITopicRepository) contract.ITopicUsecase {
	return &topicUsecase{
		topicRepo: topicRepo,
	}
}

// ListAll fetches all topic definitions.
func (uc *topicUsecase) ListAll(ctx context.Context) ([]entity.Topic, error) {
	// For now, the use case is a simple pass-through.
	// Business logic like filtering or sorting could be added here later.
	return uc.topicRepo.GetAll(ctx)
}
