package http

import (
	"net/http"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
	"github.com/RealEskalate/G6-NewsBrief/internal/handler/http/dto"

	"github.com/gin-gonic/gin"
)

type TopicHandler struct {
	topicUsecase contract.ITopicUsecase
}

func NewTopicHandler(topicUC contract.ITopicUsecase) *TopicHandler {
	return &TopicHandler{
		topicUsecase: topicUC,
	}
}

// GetTopics handles the GET /v1/topics request.
func (h *TopicHandler) GetTopics(c *gin.Context) {
	topics, err := h.topicUsecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve topics"})
		return
	}

	response := dto.TopicsResponseDTO{
		Topics:      mapTopicsToDTOs(topics),
		TotalTopics: len(topics),
	}

	c.JSON(http.StatusOK, response)
}

func mapTopicsToDTOs(topics []entity.Topic) []dto.TopicDTO {
	topicDTOs := make([]dto.TopicDTO, len(topics))
	for i, topic := range topics {
		topicDTOs[i] = dto.TopicDTO{
			Key: topic.Key,
			Label: dto.BilingualFieldDTO{
				EN: topic.Label.EN,
				AM: topic.Label.AM,
			},
			Description: dto.BilingualFieldDTO{
				EN: topic.Description.EN,
				AM: topic.Description.AM,
			},
			ImageURL:   topic.ImageURL,
			StoryCount: topic.StoryCount,
		}
	}
	return topicDTOs
}
