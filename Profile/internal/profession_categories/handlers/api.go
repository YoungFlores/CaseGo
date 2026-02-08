package handlers

import (
	"github.com/YoungFlores/Case_Go/Profile/pkg/middleware/rs256"
	"github.com/gin-gonic/gin"
)

func (h *ProfessionCategoryHandler) RegisterRoutes(rg *gin.RouterGroup, jwtMiddleware *rs256.JWTAuthMiddleware) {
	routers := rg.Group("/profession_categories")
	{
		routers.GET("", h.GetCategoriesHandler)
		routers.GET("/:id", h.GetCategoryByIDHandler)
		routers.GET("/parent/:id", h.GetCategoryByParentIDHandler)
	}
	protected := routers.Group("")
	protected.Use(jwtMiddleware.Handler())
	{
		protected.POST("", h.CreateCategoryHandler)
	}
}
