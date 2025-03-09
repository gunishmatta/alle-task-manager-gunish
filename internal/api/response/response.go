package response

import (
	"alle-task-manager-gunish/internal/common/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Err        `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Err struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Pagination *pagination.PageInfo `json:"pagination,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func SuccessWithPagination(c *gin.Context, data interface{}, pageInfo *pagination.PageInfo) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Pagination: pageInfo,
		},
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, errorCode string, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &Err{
			Code:    errorCode,
			Message: message,
		},
	})
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "NOT_FOUND", message)
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

func InternalServerError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
}
