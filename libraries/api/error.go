package api

import "net/http"

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Error is used to pass an error during the request through the
// application with web specific context.
type Error struct {
	Err           error
	Status        string
	MessageStatus string
	HTTPStatus    int
}

// ErrNewRequest wraps a provided error with an HTTP status code and custome status code. This
// function should be used when handlers encounter expected errors.
func ErrNewRequest(err error, status string, messageStatus string, httpStatus int) error {
	return &Error{err, status, messageStatus, httpStatus}
}

// ErrBadRequest wraps a provided error with an HTTP status code and custome status code for bad request. This
// function should be used when handlers encounter expected errors.
func ErrBadRequest(err error, message string) error {
	if len(message) <= 0 || message == "" {
		message = StatusMessageBadRequest
	}
	return &Error{err, StatusCodeBadRequest, message, http.StatusBadRequest}
}

// ErrNotFound wraps a provided error with an HTTP status code and custome status code for not found. This
// function should be used when handlers encounter expected errors.
func ErrNotFound(err error, message string) error {
	if len(message) <= 0 || message == "" {
		message = StatusMessageNotFound
	}
	return &Error{err, StatusCodeNotFound, message, http.StatusNotFound}
}

// ErrForbidden wraps a provided error with an HTTP status code and custome status code for forbidden. This
// function should be used when handlers encounter expected errors.
func ErrForbidden(err error, message string) error {
	if len(message) <= 0 || message == "" {
		message = StatusMessageForbidden
	}
	return &Error{err, StatusCodeForbidden, message, http.StatusForbidden}
}

// ErrInternal wraps a provided error with an HTTP status code and custome status code for internal server error. This
// function should be used when handlers encounter expected errors.
func ErrInternal(err error, message string) error {
	if len(message) <= 0 || message == "" {
		message = StatusMessageInternalServerError
	}
	return &Error{err, StatusCodeInternalServerError, message, http.StatusInternalServerError}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (err *Error) Error() string {
	return err.Err.Error()
}
