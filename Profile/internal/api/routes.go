package api

import (
	profileHandler "github.com/YoungFlores/Case_Go/Profile/internal/profile/handlers"
	"github.com/YoungFlores/Case_Go/Profile/pkg/middleware/rs256"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/YoungFlores/Case_Go/Profile/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func SetupRouter(handler *profileHandler.ProfileHandler, jwtMiddleware *rs256.JWTAuthMiddleware) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // setup later
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	protected := r.Group("/profile/api/v1")

	protected.Use(jwtMiddleware.Handler())
	{
		protected.POST("/profile", handler.CreateProfileHandler)
		protected.GET("/profile", handler.GetUserProfileHandler)
		protected.GET("/profile/:id", handler.GetUserByProfileIDHandler)
		protected.PUT("/profile", handler.UpdateProfileHandler)
		protected.PATCH("/profile", handler.PatchProfileHandler)
		protected.DELETE("/profile", handler.DeleteProfileHandler)
		protected.DELETE("/profile/:id", handler.HardDeleteHandler)

		protected.POST("/profile/social", handler.AddSocialLinkHandler)
		protected.PUT("/profile/social/:id", handler.UpdateLinkHandler)
		protected.DELETE("profile/social/:id", handler.DeleteSocialLinkHandler)

		protected.POST("/profile/purpose", handler.AddPurposesHandler)
		protected.PUT("/profile/purpose/:id", handler.UpdatePurposeHandler)
		protected.DELETE("/profile/purpose/:id", handler.DeletePuposeHandler)
	}

	return r
}
