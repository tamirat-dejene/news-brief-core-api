package contract

import (
	"context"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
)

// ITopicUsecase defines the business logic contract for topics.
type ITopicUsecase interface {
	ListAll(ctx context.Context) ([]entity.Topic, error)
}
