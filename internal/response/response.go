package response

import (
	"auth-service/internal/errors"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool `json:"success"`
	Message string `json:"message,omitempty"`
	Data any `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func Success(ctx *gin.Context, code int, message string, data any) {
	ctx.JSON(code, APIResponse{
		Success: true,
		Message: message,
		Data: data,
	})
}

func Failure(ctx *gin.Context, err *errors.AppError) {
	ctx.JSON(err.Code, APIResponse{
		Success: false,
		Error: err.Message,
	})
}