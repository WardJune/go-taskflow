package handler

import (
	"errors"
	"strconv"

	"github.com/WardJune/taskflow/internal/domain"
	"github.com/WardJune/taskflow/pkg/response"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService domain.ProjectService
	taskService    domain.TaskService
}

func NewProjectHandler(projectService domain.ProjectService, taskService domain.TaskService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		taskService:    taskService,
	}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req domain.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	ownerId := c.GetInt64("user_id")
	project, err := h.projectService.Create(c.Request.Context(), ownerId, &req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Created(c, project)
}

func (h *ProjectHandler) GetByID(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		response.BadRequest(c, errors.New("invalid project id"))
		return
	}

	requesterID := c.GetInt64("user_id")
	project, err := h.projectService.GetByID(c.Request.Context(), projectID, requesterID)
	if err != nil {
		if err.Error() == "access denied" {
			response.Forbidden(c, err)
			return
		} else if err.Error() == "project not found" {
			response.NotFound(c, err)
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, project)
}

func (h *ProjectHandler) GetMyProjects(c *gin.Context) {
	userId := c.GetInt64("user_id")
	projects, err := h.projectService.GetUserProjects(c.Request.Context(), userId)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	response.OK(c, gin.H{"projects": projects})
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	projectId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid project id"))
		return
	}

	var req domain.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	ownerId := c.GetInt64("user_id")
	if err := h.projectService.AddMember(c.Request.Context(), projectId, ownerId, &req); err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, gin.H{"message": "member added successfully"})
}

func (h *ProjectHandler) CreateTask(c *gin.Context) {
	projectId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid project id"))
		return
	}

	var req domain.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	creatorID := c.GetInt64("user_id")
	task, err := h.taskService.Create(c.Request.Context(), projectId, creatorID, &req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Created(c, gin.H{"task": task})
}

func (h *ProjectHandler) UpdateTask(c *gin.Context) {
	taskId, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid task id"))
		return
	}

	var req domain.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	requesterID := c.GetInt64("user_id")
	task, err := h.taskService.Update(c.Request.Context(), taskId, requesterID, &req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, gin.H{"task": task})
}

func (h *ProjectHandler) DeleteTask(c *gin.Context) {
	taskId, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid task id"))
		return
	}

	requesterID := c.GetInt64("user_id")
	err = h.taskService.Delete(c.Request.Context(), taskId, requesterID)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, gin.H{"message": "task deleted successfully"})
}
