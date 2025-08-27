package contract

import (
	"context"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
)

// ITopicRepository defines the persistence contract for Topic entities.
type ITopicRepository interface {
	GetAll(ctx context.Context) ([]entity.Topic, error)
}
