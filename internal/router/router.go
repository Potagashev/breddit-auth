package router

import (
	"github.com/Potagashev/breddit_auth/internal/auth"
	"github.com/gin-gonic/gin"
)

func NewRouter(AuthHandler *auth.AuthHandler) *gin.Engine {
	router := gin.Default()

	threadRoutes := router.Group("/api/v1/auth")
	{
		threadRoutes.POST("/signUp", AuthHandler.SignUp)
		threadRoutes.POST("/signIn", AuthHandler.SignIn)
		threadRoutes.GET("/verify", AuthHandler.VerifyAuthToken)
	}

	return router
}