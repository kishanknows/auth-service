package main

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	err := database.Connect(config.Conf.GetDSN())

	if err != nil {
		panic(err)
	}

	defer database.DB.Close()

	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	routes.RegisterAuthRoutes(r)

	r.Run(":8000")
}