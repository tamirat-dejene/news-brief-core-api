package usecase

import (
	"context"
	"errors"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
	"github.com/RealEskalate/G6-NewsBrief/internal/handler/http/dto"
)

const maxSubscriptionsMVP = 5

// Note the change in dependency from ISubscriptionRepository to IUserRepository.
type subscriptionUsecase struct {
	userRepo   contract.IUserRepository
	sourceRepo contract.ISourceRepository // Still needed to get source details and validate existence.
}

// The constructor now accepts IUserRepository.
func NewSubscriptionUsecase(userRepo contract.IUserRepository, sourceRepo contract.ISourceRepository) contract.ISubscriptionUsecase {
	return &subscriptionUsecase{
		userRepo:   userRepo,
		sourceRepo: sourceRepo,
	}
}

// Get now uses the IUserRepository to fetch the list of subscribed source keys.
func (uc *subscriptionUsecase) Get(ctx context.Context, userID string) (*dto.SubscriptionsResponseDTO, error) {
	// 1. Get the list of subscribed source keys directly from the user document.
	subscribedKeys, err := uc.userRepo.GetSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. Fetch the details for each subscribed source to build the rich response.
	// In a real implementation, you would make this more efficient by fetching all sources at once.
	subsDTO := make([]dto.SubscriptionDetailDTO, 0, len(subscribedKeys))
	for _, key := range subscribedKeys {
		source, err := uc.sourceRepo.GetByKey(ctx, key) // Assumes GetByKey method exists
		if err != nil || source == nil {
			// If a source was deleted but a user is still subscribed, we can just skip it.
			continue
		}

		// Here you would also need to fetch the subscription creation date.
		// For the MVP, we can omit it or use a placeholder since it's not in the user doc.
		subsDTO = append(subsDTO, dto.SubscriptionDetailDTO{
			SourceKey:  source.Key,
			SourceName: source.Name,
			// SubscribedAt: This data is no longer stored in the embedded model.
			Topics: source.Topics,
		})
	}

	return &dto.SubscriptionsResponseDTO{
		Subscriptions:      subsDTO,
		TotalSubscriptions: len(subsDTO),
		SubscriptionLimit:  maxSubscriptionsMVP,
	}, nil
}

// Create now uses the IUserRepository to add a subscription.
func (uc *subscriptionUsecase) Create(ctx context.Context, userID, sourceKey string) error {
	// 1. Check if the source exists before allowing a subscription.
	exists, err := uc.sourceRepo.Exists(ctx, sourceKey)
	if err != nil {
		return errors.New("could not validate source")
	}
	if !exists {
		return errors.New("source not found")
	}

	// 2. Check the subscription limit by getting the current count.
	currentSubs, err := uc.userRepo.GetSubscriptions(ctx, userID)
	if err != nil {
		return err
	}
	if len(currentSubs) >= maxSubscriptionsMVP {
		return errors.New("subscription limit reached")
	}

	// 3. Add the subscription directly to the user's document.
	return uc.userRepo.AddSubscription(ctx, userID, sourceKey)
}

// Delete now uses the IUserRepository to remove a subscription.
func (uc *subscriptionUsecase) Delete(ctx context.Context, userID, sourceKey string) error {
	return uc.userRepo.RemoveSubscription(ctx, userID, sourceKey)
}
