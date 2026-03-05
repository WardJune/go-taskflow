package handler

import (
	"errors"
	"strconv"

	"github.com/WardJune/taskflow/internal/domain"
	ws "github.com/WardJune/taskflow/internal/websocket"
	"github.com/WardJune/taskflow/pkg/response"
	"github.com/gin-gonic/gin"
)

type WSHandler struct {
	hub         *ws.Hub
	projectRepo domain.ProjectRepository
}

func NewWSHandler(hub *ws.Hub, projectRepo domain.ProjectRepository) *WSHandler {
	return &WSHandler{
		hub:         hub,
		projectRepo: projectRepo,
	}
}

func (h *WSHandler) HandleWS(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid project id"))
		return
	}

	userID := c.GetInt64("user_id")

	isMember, err := h.projectRepo.IsMember(c.Request.Context(), projectID, userID)

	if err != nil || !isMember {
		response.Forbidden(c, err)
		return
	}

	ws.ServeWS(h.hub, userID, projectID, c.Writer, c.Request)
}
