package errors

import "net/http"

var (
	ErrInternalServer = New(
		http.StatusInternalServerError,
		"internal server error",
		nil,
	)

	ErrInvalidRequestBody = New(
		http.StatusBadRequest,
		"Invalid input format",
		nil,
	)

	ErrDuplicateRegistration = New(
		http.StatusConflict,
		"user already exists",
		nil,
	)

	ErrInvalidCredentials = New(
		http.StatusUnauthorized,
		"invalid credentials",
		nil,
	)

	ErrUserNotFound = New(
		http.StatusNotFound,
		"user not found",
		nil,
	)

	ErrUnauthorized = New(
		http.StatusUnauthorized,
		"user not authenticated",
		nil,
	)

	ErrSessionExpired = New(
		http.StatusUnauthorized,
		"session expired",
		nil,
	)
)