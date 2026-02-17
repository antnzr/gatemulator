package apperror

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewHTTPError(status int, message string) error {
	return &HTTPError{
		StatusCode: status,
		Message:    message,
	}
}

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func WriteError(w http.ResponseWriter, err error) {
	// Check if the error is an HTTPError with a specific status code
	if httpErr, ok := err.(*HTTPError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpErr.StatusCode)
		json.NewEncoder(w).Encode(ErrorResponse{Message: httpErr.Message, Code: httpErr.StatusCode})
		return
	}

	// Default case for unhandled errors
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error(), Code: http.StatusInternalServerError})
}

type ErrorHandler func(w http.ResponseWriter, r *http.Request) error

// ErrorHandlerWrapper wraps an ErrorHandler function, handling errors in one place
func ErrorHandlerWrapper(fn ErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			slog.Error("", slog.Any("err", err))
			WriteError(w, err)
		}
	}
}
