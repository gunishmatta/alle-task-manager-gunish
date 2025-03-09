package handler

import (
	"alle-task-manager-gunish/internal/api/response"
	"alle-task-manager-gunish/internal/common/errors"
	"alle-task-manager-gunish/internal/common/pagination"
	"alle-task-manager-gunish/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (handler *TaskHandler) RegisterRoutes(router *gin.Engine) {
	tasks := router.Group("/tasks")
	{
		tasks.GET("", handler.ListTasks)
		tasks.POST("", handler.CreateTask)
		tasks.GET("/:id", handler.GetTask)
		tasks.PUT("/:id", handler.UpdateTask)
		tasks.DELETE("/:id", handler.DeleteTask)
	}
}

func (handler *TaskHandler) CreateTask(c *gin.Context) {
	var input service.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid request payload: "+err.Error())
		return
	}

	task, err := handler.taskService.CreateTask(c.Request.Context(), input)
	if err != nil {
		handler.handleError(c, err)
		return
	}

	response.Created(c, task)
}

func (handler *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := handler.taskService.GetTask(c.Request.Context(), id)
	if err != nil {
		handler.handleError(c, err)
		return
	}

	response.Success(c, task)
}

func (handler *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var input service.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid request payload: "+err.Error())
		return
	}

	task, err := handler.taskService.UpdateTask(c.Request.Context(), id, input)
	if err != nil {
		handler.handleError(c, err)
		return
	}

	response.Success(c, task)
}

func (handler *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := handler.taskService.DeleteTask(c.Request.Context(), id)
	if err != nil {
		handler.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (handler *TaskHandler) ListTasks(c *gin.Context) {
	status := c.Query("status")

	pageNum, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if pageNum < 1 {
		pageNum = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	page := &pagination.Page{
		Number: pageNum,
		Size:   pageSize,
	}

	tasks, pageInfo, err := handler.taskService.ListTasks(c.Request.Context(), status, page)
	if err != nil {
		handler.handleError(c, err)
		return
	}

	response.SuccessWithPagination(c, tasks, pageInfo)
}

func (handler *TaskHandler) handleError(c *gin.Context, err error) {
	switch err {
	case errors.ErrNotFound:
		response.NotFound(c, "Task not found")
	case errors.ErrDuplicateEntity:
		response.BadRequest(c, "Task with this ID already exists")
	case errors.ErrInvalidStatus:
		response.BadRequest(c, "Invalid task status")
	default:
		response.InternalServerError(c)
	}
}
