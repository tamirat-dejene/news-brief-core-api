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
	uuidGen      contract.IUUIDGenerator
}

func NewTopicHandler(topicUC contract.ITopicUsecase, uuidGen contract.IUUIDGenerator) *TopicHandler {
	return &TopicHandler{
		topicUsecase: topicUC,
		uuidGen:      uuidGen,
	}
}
func (h *TopicHandler) CreateTopic(c *gin.Context) {
	var topicDTO dto.TopicDTO
	if err := c.ShouldBindJSON(&topicDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	topic := entity.Topic{
		ID:   h.uuidGen.NewUUID(),
		Slug: topicDTO.Slug,
		Label: entity.BilingualField{
			EN: topicDTO.Label.EN,
			AM: topicDTO.Label.AM,
		},
		StoryCount: 0,
	}

	if err := h.topicUsecase.CreateTopic(c.Request.Context(), &topic); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create topic"})
		return
	}

	c.JSON(http.StatusCreated, dto.MessageResponse{Message: "Topic created successfully"})
}

// GetTopics handles the GET /v1/topics request.
func (h *TopicHandler) GetTopics(c *gin.Context) {
	topics, err := h.topicUsecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve topics"})
		return
	}

	response := dto.TopicsResponseDTO{
		Topics:      dto.MapTopicsToDTOs(topics),
		TotalTopics: len(topics),
	}

	c.JSON(http.StatusOK, response)
}
