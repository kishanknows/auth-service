package handlers

import (
	"auth-service/internal/dto"
	"auth-service/internal/errors"
	"auth-service/internal/response"
	"auth-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService 
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		service: services.NewAuthService(),
	}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var user dto.UserRegister

	if err := ctx.ShouldBindJSON(&user); err != nil {
		response.Failure(ctx, errors.ErrInvalidRequestBody)
		return
	}

	if err := user.Validate(); err != nil {
		response.Failure(ctx, errors.ErrInvalidRequestBody)
		return
	}

	res, err := h.service.RegisterUser(&user)

	if err != nil {
		response.Failure(ctx, err)
		return
	}

	response.Success(ctx, http.StatusCreated, "user created", res)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var login dto.UserLogin

	if err := ctx.ShouldBindJSON(&login); err != nil {
		response.Failure(ctx, errors.ErrInvalidRequestBody)
		return
	}

	res, err := h.service.LoginUser(&login)

	if err != nil {
		response.Failure(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "login successful", *res)
}

func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var refresh dto.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&refresh); err != nil {
		response.Failure(ctx, errors.ErrInvalidRequestBody)
		return
	}

	res, err := h.service.RefreshToken(&refresh)

	if err != nil {
		response.Failure(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "token refreshed", res)
}
