package dto

type LoginResponse struct {
	RefreshToken *string `json:"refresh_token"`
	Token *string `json:"token"`
	ExpiresIn float64 `json:"expires_in"`
	TokenType string `json:"token_type"`
}