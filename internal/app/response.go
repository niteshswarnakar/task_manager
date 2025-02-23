package app

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaginatedJSONResult struct {
	JSONResult
	Pagination Pagination `json:"pagination"`
}

type JSONResult struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
}

type HTTPError struct {
	Message string `json:"message"`
}

type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func PaginatedResponse(c *gin.Context, data any, total int64, page int, limit int) {
	result := PaginatedJSONResult{
		JSONResult: JSONResult{
			Data: data,
		},
		Pagination: Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}
	c.JSON(http.StatusOK, result)
}

func HttpErrorResult(ctx *gin.Context, statusCode int, message string) {
	result := HTTPError{
		Message: message,
	}
	ctx.JSON(statusCode, result)
}

func NotFoundResponse(ctx *gin.Context, message string) {
	HttpErrorResult(ctx, http.StatusNotFound, message)
}

func InvalidRequest(ctx *gin.Context) {
	HttpErrorResult(ctx, http.StatusBadRequest, "invalid request")
}

func BadRequest(ctx *gin.Context, message string) {
	HttpErrorResult(ctx, http.StatusBadRequest, message)
}

func ErrorResponse(ctx *gin.Context, err error) {
	var appError AppError
	if errors.As(err, &appError) {
		HttpErrorResult(ctx, appError.Code, appError.Message)
	}
	InternalServerError(ctx)
}

func InternalServerError(c *gin.Context) {
	ServerError(c, "Internal Server Error")
}

func SuccessResponse(ctx *gin.Context, data any) {
	result := JSONResult{
		Data: data,
	}
	ctx.JSON(http.StatusOK, result)
}

func ServerError(c *gin.Context, message string) {
	result := HTTPError{
		Message: message,
	}
	c.JSON(http.StatusInternalServerError, result)
}

type AppError struct {
	Code    int
	Message string
}
