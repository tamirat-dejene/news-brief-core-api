package dto

import (
	"time"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
)

// UserResponse is the DTO for a user.
type UserResponse struct {
	ID          string         `json:"id"`
	Username    string         `json:"username"`
	Fullname    string         `json:"fullname"`
	Email       string         `json:"email"`
	Role        string         `json:"role"`
	AvatarURL   *string        `json:"avatar_url"`
	CreatedAt   string         `json:"created_at"`
	Preferences PreferencesDTO `json:"preferences"`
}

// LoginResponse is the DTO for a successful login.
type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

// converts an entity.User to a UserResponse DTO.
func ToUserResponse(user entity.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Fullname:  user.Fullname,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		Preferences: PreferencesDTO{
			Topics:            user.Preferences.Topics,
			SubscribedSources: user.Preferences.SubscribedSources,
			Notifications: NotificationsDTO{ // This assumes your entity.Preferences has this nested struct
				DailyBrief:   user.Preferences.Notifications.DailyBrief,
				BreakingNews: user.Preferences.Notifications.BreakingNews,
			},
		},
	}
}

// MessageResponse is a generic response for success/error messages.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse is a response for errors.
type ErrorResponse struct {
	Error string `json:"error"`
}

// BilingualFieldDTO for API responses.
type BilingualFieldDTO struct {
	EN string `json:"en"`
	AM string `json:"am"`
}

// SubscriptionDetailDTO now matches the API spec.
type SubscriptionDetailDTO struct {
	SourceSlug   string    `json:"source_slug"`
	SourceName   string    `json:"source_name"`
	SubscribedAt time.Time `json:"subscribed_at"`
	Topics       []string  `json:"topics"` // Per-subscription topics can be a future enhancement
}

// SubscriptionsResponseDTO now matches the API spec.
type SubscriptionsResponseDTO struct {
	Subscriptions      []SubscriptionDetailDTO `json:"subscriptions"`
	TotalSubscriptions int                     `json:"total_subscriptions"`
	SubscriptionLimit  int                     `json:"subscription_limit"`
}

// PreferencesDTO matches the nested preferences object in the API spec responses.
type PreferencesDTO struct {
	Lang              string           `json:"lang"`
	Topics            []string         `json:"topics"`             // This field is now correctly included
	SubscribedSources []string         `json:"subscribed_sources"` // This field is now correctly included
	BriefType         string           `json:"brief_type"`
	DataSaver         bool             `json:"data_saver"`
	Notifications     NotificationsDTO `json:"notifications"`
}

// SourceDTO represents a single source in an API response.
type SourceDTO struct {
	Slug             string   `json:"slug"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	URL              string   `json:"url"`
	LogoURL          string   `json:"logo_url"`
	Languages        string   `json:"languages"`
	Topics           []string `json:"topics"`
	ReliabilityScore float64  `json:"reliability_score"`
}

// SourcesResponseDTO is the response for the GET /v1/sources endpoint.
type SourcesResponseDTO struct {
	Sources      []SourceDTO `json:"sources"`
	TotalSources int         `json:"total_sources"`
}

// MapSourcesToDTOs converts a slice of source entities to a slice of DTOs.
func MapSourcesToDTOs(sources []entity.Source) []SourceDTO {
	sourceDTOs := make([]SourceDTO, len(sources))
	for i, source := range sources {
		sourceDTOs[i] = SourceDTO{
			Slug:             source.Slug,
			Name:             source.Name,
			Description:      source.Description,
			URL:              source.URL,
			LogoURL:          source.LogoURL,
			Languages:        string(source.Languages),
			Topics:           source.Topics,
			ReliabilityScore: source.ReliabilityScore,
		}
	}
	return sourceDTOs
}

// NotificationsDTO defines the nested notifications object in API responses.
type NotificationsDTO struct {
	DailyBrief   bool `json:"daily_brief"`
	BreakingNews bool `json:"breaking_news"`
}

// topics
// TopicDTO represents a single topic in the API response.
type TopicDTO struct {
	Slug       string            `json:"slug"`
	TopicName  string            `json:"topic_name"`
	Label      BilingualFieldDTO `json:"label"`
	StoryCount int               `json:"story_count"`
}

// TopicsResponseDTO is the response for the GET /v1/topics endpoint.
type TopicsResponseDTO struct {
	Topics      []TopicDTO `json:"topics"`
	TotalTopics int        `json:"total_topics"`
}

func MapTopicsToDTOs(topics []entity.Topic) []TopicDTO {
	topicDTOs := make([]TopicDTO, len(topics))
	for i, topic := range topics {
		topicDTOs[i] = TopicDTO{
			Slug: topic.Slug,
			Label: BilingualFieldDTO{
				EN: topic.Label.EN,
				AM: topic.Label.AM,
			},
			StoryCount: topic.StoryCount,
		}
	}
	return topicDTOs
}
