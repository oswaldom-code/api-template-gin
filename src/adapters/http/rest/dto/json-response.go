package dto

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	ErrUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrForbidden      ErrorCode = "FORBIDDEN"
	ErrNotFound       ErrorCode = "NOT_FOUND"
	ErrBadRequest     ErrorCode = "BAD_REQUEST"
	ErrValidation     ErrorCode = "VALIDATION_ERROR"
	ErrConflict       ErrorCode = "CONFLICT"
	ErrInternalServer ErrorCode = "INTERNAL_ERROR"
	ErrServiceUnavail ErrorCode = "SERVICE_UNAVAILABLE"
)

type Meta struct {
	Timestamp string `json:"timestamp"`
}

type SuccessResponse struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}

type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func newMeta() Meta {
	return Meta{Timestamp: time.Now().UTC().Format(time.RFC3339)}
}

func Success(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, SuccessResponse{
		Data: data,
		Meta: newMeta(),
	})
}

func Error(c *gin.Context, statusCode int, code ErrorCode, message string) {
	c.JSON(statusCode, ErrorResponse{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

func AbortWithError(c *gin.Context, statusCode int, code ErrorCode, message string) {
	c.AbortWithStatusJSON(statusCode, ErrorResponse{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

func OK(c *gin.Context, data any) {
	Success(c, http.StatusOK, data)
}

func Created(c *gin.Context, data any) {
	Success(c, http.StatusCreated, data)
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, ErrBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	AbortWithError(c, http.StatusUnauthorized, ErrUnauthorized, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, ErrNotFound, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, ErrInternalServer, message)
}
