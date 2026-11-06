package response

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	TotalRows  int64 `json:"total_rows,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

type Response struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func JSON(c *gin.Context, statusCode int, success bool, message string, data any, meta *Meta, errors any) {
	c.JSON(statusCode, Response{
		Success: success,
		Code:    statusCode,
		Message: message,
		Data:    data,
		Meta:    meta,
		Errors:  errors,
	})
}

func Success(c *gin.Context, statusCode int, message string, data any) {
	JSON(c, statusCode, true, message, data, nil, nil)
}

func OK(c *gin.Context, message string, data any) {
	Success(c, http.StatusOK, message, data)
}

func Created(c *gin.Context, message string, data any) {
	Success(c, http.StatusCreated, message, data)
}

func Paginated(c *gin.Context, message string, data any, page, limit int, totalRows int64) {
	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(totalRows) / float64(limit)))
	}

	meta := &Meta{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}

	JSON(c, http.StatusOK, true, message, data, meta, nil)
}

func Error(c *gin.Context, statusCode int, message string, errDetail any) {
	JSON(c, statusCode, false, message, nil, nil, errDetail)
}

func BadRequest(c *gin.Context, message string, errDetail any) {
	Error(c, http.StatusBadRequest, message, errDetail)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, nil)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}

func InternalError(c *gin.Context, message string, errDetail any) {
	Error(c, http.StatusInternalServerError, message, errDetail)
}
