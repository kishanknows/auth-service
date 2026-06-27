package middleware

import (
	"auth-service/internal/errors"
	"auth-service/internal/models"
	"auth-service/internal/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.Failure(ctx, errors.ErrUnauthorized)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Failure(ctx, errors.ErrUnauthorized)
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		claims := &models.TokenClaims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return jwtSecret, nil
			},
		)

		if err != nil || !token.Valid {
			response.Failure(ctx, errors.ErrUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("email", claims.Email)

		ctx.Next()
	}
}
