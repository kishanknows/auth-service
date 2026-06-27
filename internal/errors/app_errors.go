package errors

type AppError struct {
	Code int
	Message string
	Err error
}

func New(code int, message string, err error) *AppError {
	return &AppError{
		Code: code,
		Message: message,
		Err: err,
	}
}