package http

import (
	"net/http"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
	"github.com/RealEskalate/G6-NewsBrief/internal/handler/http/dto"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subUsecase contract.ISubscriptionUsecase
}

func NewSubscriptionHandler(subUC contract.ISubscriptionUsecase) *SubscriptionHandler {
	return &SubscriptionHandler{
		subUsecase: subUC,
	}
}

func (h *SubscriptionHandler) GetSubscriptions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	resp, err := h.subUsecase.Get(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *SubscriptionHandler) AddSubscription(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.AddSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.subUsecase.Create(c.Request.Context(), userID.(string), req.SourceKey)
	if err != nil {
		// Differentiate errors for better client feedback
		if err.Error() == "subscription limit reached" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *SubscriptionHandler) RemoveSubscription(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sourceKey := c.Param("source_key")
	if sourceKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_key is required"})
		return
	}

	err := h.subUsecase.Delete(c.Request.Context(), userID.(string), sourceKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unsubscribed"})
}
