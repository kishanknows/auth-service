package services

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/dto"
	"auth-service/internal/errors"
	"auth-service/internal/models"
	"auth-service/internal/utils"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) RegisterUser(user *dto.UserRegister) (*dto.LoginResponse, *errors.AppError) {
	// Check for duplicate registration
	var exists bool
	err := database.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE email = $1)
		`, user.Email).Scan(&exists)
	
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	if exists {
		fmt.Println("user already exists")
		return nil, errors.ErrDuplicateRegistration
	}

	// Generate password hash
	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	// Begin database transaction
	tx, err := database.DB.Begin()
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}
	defer tx.Rollback()

	// Insert user into the database
	var userID int
	err = tx.QueryRow(
		`INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id`, user.Email, hashedPassword, user.Role).Scan(&userID)

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	// Create refresh token and Store hash in DB
	refreshToken, refreshTokenHash, err := utils.GenerateRefreshToken()

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	_, err = tx.Exec(`
		INSERT INTO refresh_tokens (user_id, token_hash, created_at, expires_at)
		VALUES ($1, $2, $3, $4)
	`, userID, refreshTokenHash, time.Now(), time.Now().Add(time.Hour * 24 * 30))

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	// Generate Bearer Token
	claims := models.TokenClaims{
		UserID: userID,
		Email: user.Email,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Conf.JWT.TokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	bearerToken, err := token.SignedString([]byte(config.Conf.JWT.Secret))

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	// Commit the DB transaction
	if err := tx.Commit(); err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	return &dto.LoginResponse{
		RefreshToken: &refreshToken,
		Token: &bearerToken,
		ExpiresIn: config.Conf.JWT.TokenExpiry.Seconds(),
		TokenType: "Bearer",
	}, nil
}

func (s *AuthService) LoginUser(login *dto.UserLogin) (*dto.LoginResponse, *errors.AppError) {
	// fetch the user from database
	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, email, password_hash, role
		FROM users
		WHERE email = $1
	`, login.Email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)

	if err == sql.ErrNoRows {
		fmt.Println(err.Error())
		return nil, errors.ErrUserNotFound
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	// Verify user credentials
	if !utils.CheckPasswordHash(login.Password, user.PasswordHash) {
		fmt.Println("Invalid password")
		return nil, errors.ErrInvalidCredentials
	}

	// Generate Bearer Token
	claims := models.TokenClaims{
		UserID: user.ID,
		Email: user.Email,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Conf.JWT.TokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	bearerToken, err := token.SignedString([]byte(config.Conf.JWT.Secret))
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Generate Refresh Token and store in DB
	refreshToken, refreshTokenHash, err := utils.GenerateRefreshToken()

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	_, err = database.DB.Exec(`
		UPDATE refresh_tokens 
		SET token_hash = $1, created_at = $2, expires_at = $3
		WHERE user_id = $4
	`, refreshTokenHash, time.Now(), time.Now().Add(time.Hour * 24 * 30), user.ID)

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	return &dto.LoginResponse{
		RefreshToken: &refreshToken,
		Token: &bearerToken,
		ExpiresIn: config.Conf.JWT.TokenExpiry.Seconds(),
		TokenType: "Bearer",
	}, nil
}

func (s *AuthService) RefreshToken(refresh *dto.RefreshTokenRequest) (*dto.LoginResponse, *errors.AppError) {
	// Verify the refresh token
	tokenHash := utils.GenerateTokenHash(refresh.RefreshToken)

	var userID int
	var tokenExpiry time.Time

	err := database.DB.QueryRow(`
		SELECT user_id, expires_at 
		FROM refresh_tokens 
		WHERE token_hash = $1
	`, tokenHash).Scan(&userID, &tokenExpiry)
	
	if err == sql.ErrNoRows {
		fmt.Println("invalid refresh token")
		return nil, errors.ErrSessionExpired
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	if time.Now().After(tokenExpiry) {
		fmt.Println("session expired")
		return nil, errors.ErrSessionExpired
	}

	// Generate new refresh token and store in DB
	refreshToken, refreshTokenHash, err := utils.GenerateRefreshToken()

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	_, err = database.DB.Exec(`
		UPDATE refresh_tokens 
		SET token_hash = $1, created_at = $2, expires_at = $3
		WHERE user_id = $4
	`, refreshTokenHash, time.Now(), time.Now().Add(time.Hour * 24 * 30), userID)

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	var email string
	var role models.UserRole
	err = database.DB.QueryRow(`
		SELECT email, role
		FROM users
		WHERE id = $1`, userID).Scan(&email, &role)

	// Create Bearer Token
	claims := models.TokenClaims{
		UserID: userID,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Conf.JWT.TokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Conf.JWT.Secret))

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	return &dto.LoginResponse{
		RefreshToken: &refreshToken,
		Token: &tokenString,
		ExpiresIn: config.Conf.JWT.TokenExpiry.Seconds(),
		TokenType: "Bearer",
	}, nil
}