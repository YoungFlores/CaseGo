package handlers

import (
	"github.com/YoungFlores/Case_Go/Profile/pkg/middleware/rs256"
	"github.com/gin-gonic/gin"
)

func (h *SearchHandler) RegisterRoutes(rg *gin.RouterGroup, middleware *rs256.JWTAuthMiddleware) {
	router := rg.Group("/search")
	router.Use(middleware.Handler())
	{
		router.GET("", h.GetProfilesHandler)
		router.GET("/fio", h.GetProfileFioHandler)
	}

}
