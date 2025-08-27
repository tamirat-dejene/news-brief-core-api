package usecase

import (
	"context"
	"errors"

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
func (uc *topicUsecase) CreateTopic(ctx context.Context, topic *entity.Topic) error {
	if topic == nil {
		return errors.New("topic cannot be nil")
	}
	if topic.Slug == "" {
		return errors.New("topic slug cannot be empty")
	}
	exists, err := uc.topicRepo.CheckSlugExists(ctx, topic.Slug)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("topic with slug already exists")
	}
	return uc.topicRepo.CreateTopic(ctx, topic)
}

// ListAll fetches all topic definitions.
func (uc *topicUsecase) ListAll(ctx context.Context) ([]entity.Topic, error) {
	// For now, the use case is a simple pass-through.
	// Business logic like filtering or sorting could be added here later.
	return uc.topicRepo.GetAll(ctx)
}
