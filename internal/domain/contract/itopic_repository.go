package contract

import (
	"context"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
)

// ITopicRepository defines the persistence contract for Topic entities.
type ITopicRepository interface {
	CreateTopic(ctx context.Context, topic *entity.Topic) error
	// CheckSlugExists checks if a topic with the given slug exists.
	GetTopicByID(ctx context.Context, topicID string) (*entity.Topic, error)
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
	GetAll(ctx context.Context) ([]entity.Topic, error)
	GetUserSubscribedTopics(ctx context.Context, userTopics []string) ([]*entity.Topic, error)
}
