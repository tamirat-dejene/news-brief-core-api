package http

import (
	"net/http"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
	"github.com/RealEskalate/G6-NewsBrief/internal/handler/http/dto"
	"github.com/gin-gonic/gin"
)

type SourceHandler struct {
	sourceUsecase contract.ISourceUsecase
}

func NewSourceHandler(sourceUC contract.ISourceUsecase) *SourceHandler {
	return &SourceHandler{
		sourceUsecase: sourceUC,
	}
}

// GetSources handles the GET /v1/sources request.
func (h *SourceHandler) GetSources(c *gin.Context) {
	sources, err := h.sourceUsecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sources"})
		return
	}

	response := dto.SourcesResponseDTO{
		Sources:      mapSourcesToDTOs(sources),
		TotalSources: len(sources),
	}

	c.JSON(http.StatusOK, response)
}

// mapSourcesToDTOs converts a slice of source entities to a slice of DTOs.
func mapSourcesToDTOs(sources []entity.Source) []dto.SourceDTO {
	sourceDTOs := make([]dto.SourceDTO, len(sources))
	for i, source := range sources {
		sourceDTOs[i] = dto.SourceDTO{
			Key:              source.Key,
			Name:             source.Name,
			Description:      source.Description,
			URL:              source.URL,
			LogoURL:          source.LogoURL,
			Languages:        source.Languages,
			Topics:           source.Topics,
			ReliabilityScore: source.ReliabilityScore,
		}
	}
	return sourceDTOs
}
