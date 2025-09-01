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
	uuidGen       contract.IUUIDGenerator
}

func NewSourceHandler(sourceUC contract.ISourceUsecase, uuidGen contract.IUUIDGenerator) *SourceHandler {
	return &SourceHandler{
		sourceUsecase: sourceUC,
		uuidGen:       uuidGen,
	}
}

// CreateSource handles the POST /v1
func (h *SourceHandler) CreateSource(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Forbidden: Admins only"})
		return
	}

	var sourceDTO dto.SourceDTO
	if err := c.ShouldBindJSON(&sourceDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	source := entity.Source{
		ID:               h.uuidGen.NewUUID(),
		Slug:             sourceDTO.Slug,
		Name:             sourceDTO.Name,
		Description:      sourceDTO.Description,
		URL:              sourceDTO.URL,
		LogoURL:          sourceDTO.LogoURL,
		Languages:        entity.SetLanguageType(sourceDTO.Languages),
		Topics:           sourceDTO.Topics,
		ReliabilityScore: sourceDTO.ReliabilityScore,
	}

	if err := h.sourceUsecase.CreateSource(c.Request.Context(), &source); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create source"})
		return
	}

	c.JSON(http.StatusCreated, dto.MessageResponse{Message: "Source created successfully"})
}

// GetSources handles the GET /v1/sources request.
func (h *SourceHandler) GetSources(c *gin.Context) {
	sources, err := h.sourceUsecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sources"})
		return
	}

	response := dto.SourcesResponseDTO{
		Sources:      dto.MapSourcesToDTOs(sources),
		TotalSources: len(sources),
	}

	c.JSON(http.StatusOK, response)
}
