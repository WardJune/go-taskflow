package handler

import (
	"context"
	"time"

	"github.com/WardJune/taskflow/internal/domain"
	"github.com/WardJune/taskflow/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req domain.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	resp, err := h.userService.Register(ctx, &req)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	response.Created(c, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	resp, err := h.userService.Login(ctx, &req)
	if err != nil {
		response.Unauthorized(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetInt64("user_id")
	email := c.GetString("user_email")

	response.OK(c, gin.H{
		"user_id": userID,
		"email":   email,
	})
}
