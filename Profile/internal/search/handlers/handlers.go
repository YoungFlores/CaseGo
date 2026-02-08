package handlers

import (
	"net/http"

	"github.com/YoungFlores/Case_Go/Profile/internal/search/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/search/service"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	service *service.SearchService
}

func NewSearchHandler(service *service.SearchService) *SearchHandler {
	return &SearchHandler{
		service: service,
	}
}

func (h *SearchHandler) GetProfilesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.SearchDTO
	var helpers dto.SearchHelpersDTO

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid search params"})
		return
	}
	if err := c.ShouldBindQuery(&helpers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination params"})
		return
	}

	res, err := h.service.SearchProfileService(ctx, req, helpers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}

func (h *SearchHandler) GetProfileFioHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.SearchByFIODTO
	var helpers dto.SearchHelpersDTO

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid search params"})
		return
	}
	if err := c.ShouldBindQuery(&helpers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination params"})
		return
	}

	res, err := h.service.SearchByFioService(ctx, req, helpers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
