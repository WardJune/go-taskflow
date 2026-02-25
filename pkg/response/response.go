package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error:   err.Error(),
	})
}

func NotFound(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error:   err.Error(),
	})
}

func InternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error:   "internal server error",
	})
}

func Unauthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Error:   err.Error(),
	})
}

func Forbidden(c *gin.Context, err error) {
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Error:   err.Error(),
	})
}
