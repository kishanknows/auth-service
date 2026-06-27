package routes

import (
	"auth-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	authHandler := handlers.NewAuthHandler()

	public := r.Group("/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/refresh-token", authHandler.RefreshToken)
	}
}