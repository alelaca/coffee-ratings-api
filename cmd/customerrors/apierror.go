package customerrors

import (
	"fmt"
)

type ApiError struct {
	HttpStatusCode int
	Message        string
}

func CreateApiError(httpStatusCode int, message string) ApiError {
	return ApiError{
		HttpStatusCode: httpStatusCode,
		Message:        message,
	}
}

func (e ApiError) Error() string {
	return fmt.Sprintf(`Error: %s`, e.Message)
}
