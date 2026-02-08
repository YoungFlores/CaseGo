package handlers

import (
	"net/http"
	"strconv"

	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/models"
	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/service"
	userModel "github.com/YoungFlores/Case_Go/Profile/internal/profile/models"

	"github.com/gin-gonic/gin"
)

type ProfessionCategoryHandler struct {
	service *service.ProfessionCategoryService
}

func NewProfessionCategoryHandler(service *service.ProfessionCategoryService) *ProfessionCategoryHandler {
	return &ProfessionCategoryHandler{
		service: service,
	}
}

func (h *ProfessionCategoryHandler) GetRole(c *gin.Context) (userModel.UserRole, bool) {
	userRole, exist := c.Get("role")
	if !exist {
		return 0, false
	}

	userRoleInt, ok := userRole.(userModel.UserRole)
	if !ok {
		return 0, false
	}

	return userRoleInt, true

}

func (h *ProfessionCategoryHandler) CreateCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userRole, exist := h.GetRole(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if userRole != userModel.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var dto models.CategoryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.service.CreateCategoryService(ctx, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)

}

func (h *ProfessionCategoryHandler) GetCategoriesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	category, err := h.service.GetCategoriesService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *ProfessionCategoryHandler) GetCategoryByParentIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	parentID, err := strconv.Atoi(c.Param("parentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parentID"})
		return
	}

	categories, err := h.service.GetCategoriesByParentService(ctx, int16(parentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)

}

func (h *ProfessionCategoryHandler) GetCategoryByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	category, err := h.service.GetCategoryByIDService(ctx, int16(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)

}
