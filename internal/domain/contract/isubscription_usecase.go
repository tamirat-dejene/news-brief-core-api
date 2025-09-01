package contract

import (
	"context"

	"github.com/RealEskalate/G6-NewsBrief/internal/handler/http/dto"
)

// ISubscriptionUsecase defines the business logic for user subscriptions.
type ISubscriptionUsecase interface {
	Get(ctx context.Context, userID string) (*dto.SubscriptionsResponseDTO, error)
	Create(ctx context.Context, userID, sourceKey string) error
	Delete(ctx context.Context, userID, sourceKey string) error
}
