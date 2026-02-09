package errs

import "net/http"

type MessageError interface {
	error
	Status() int
	Message() string
}

type messageError struct {
	ErrStatus  int    `json:"status"`
	ErrMessage string `json:"message"`
}

func (e *messageError) Error() string {
	return e.ErrMessage
}

func (e *messageError) Status() int {
	return e.ErrStatus
}

func (e *messageError) Message() string {
	return e.ErrMessage
}

func BadRequest(msg string) MessageError {
	return &messageError{ErrStatus: http.StatusBadRequest, ErrMessage: msg}
}

func NotFound(msg string) MessageError {
	return &messageError{ErrStatus: http.StatusNotFound, ErrMessage: msg}
}

func InternalServerError(msg string) MessageError {
	return &messageError{ErrStatus: http.StatusInternalServerError, ErrMessage: msg}
}

func Unauthorized(msg string) MessageError {
	return &messageError{ErrStatus: http.StatusUnauthorized, ErrMessage: msg}
}

func Forbidden(msg string) MessageError {
	return &messageError{ErrStatus: http.StatusForbidden, ErrMessage: msg}
}
