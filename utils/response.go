package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Api Response
type ApiResponse struct {
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func Success(ctx *gin.Context, message string, statusCode int, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	ctx.JSON(statusCode, ApiResponse{
		StatusCode: statusCode,
		Success:    true,
		Message:    message,
		Data:       data,
	})
}

func Error(ctx *gin.Context, message string, statusCode int) {
	ctx.JSON(statusCode, ApiResponse{
		StatusCode: statusCode,
		Success:    false,
		Message:    message,
		Data:       gin.H{},
	})
}

func BadRequest(ctx *gin.Context, message string) {
	Error(ctx, message, http.StatusBadRequest)
}

func InternalServerError(ctx *gin.Context, message string) {
	Error(ctx, message, http.StatusInternalServerError)
}

func Unauthorized(ctx *gin.Context, message string) {
	Error(ctx, message, http.StatusUnauthorized)
}

func NotFound(ctx *gin.Context, message string) {
	Error(ctx, message, http.StatusNotFound)
}

func Forbidden(ctx *gin.Context, message string) {
	Error(ctx, message, http.StatusForbidden)
}

func Conflict(ctx *gin.Context, message string) {
	Error(ctx, message, http.StatusConflict)
}

func UnprocessableEntity(ctx *gin.Context, message string) {
	Error(ctx, message, http.StatusUnprocessableEntity)
}

func ServiceUnavailable(ctx *gin.Context, message string) {
	Error(ctx, message, http.StatusServiceUnavailable)
}
