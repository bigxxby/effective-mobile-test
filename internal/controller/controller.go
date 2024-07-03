package controller

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/bigxxby/effective-mobile-test/internal/models"
	"github.com/bigxxby/effective-mobile-test/internal/service"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service *service.Service
}

func New(service service.Service) Controller {
	return Controller{
		Service: &service,
	}
}

// GetUsers godoc
// @Summary Get users with optional filtering, pagination, and sorting.
// @Description Retrieves a list of users based on optional filters, paginated results, and sorting criteria.
// @Produce json
// @Param passport_number query string false "Passport number to filter users"
// @Param surname query string false "Surname to filter users"
// @Param name query string false "Name to filter users"
// @Param page query int false "Page number for pagination (default 1)"
// @Param page_size query int false "Number of users per page (default 10)"
// @Param sort_by query string false "Field to sort by (default 'id')"
// @Param sort_order query string false "Sort order, either 'asc' or 'desc' (default 'asc')"
// @Success 200 {object} models.ResponseUsersList "Successful response with list of users"
// @Failure 404 {object} models.ErrorResponse "Users not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/users [get]
func (c *Controller) GetUsers(ctx *gin.Context) {
	var filter models.Filter
	var pagination models.Pagination

	filter.PassportNumber = ctx.Query("passport_number")
	filter.Surname = ctx.Query("surname")
	filter.Name = ctx.Query("name")

	// Валидация параметров
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil {
		pageSize = 10
	}

	sortBy := ctx.DefaultQuery("sort_by", "id")
	sortOrder := ctx.DefaultQuery("sort_order", "asc")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	pagination.Page = page
	pagination.PageSize = pageSize

	users, err := c.Service.GetUsers(filter, pagination, sortBy, sortOrder)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(404, gin.H{"error": "Users not found"})
			return
		}
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"users": users})
}

// GetUser godoc
// @Summary Get a user by ID.
// @Description Retrieves a user by their ID.
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User "Successful response with user details"
// @Failure 400 {object} models.ErrorResponse "Invalid user ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/users/{id} [get]
func (c *Controller) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	uid, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}
	if uid <= 0 {
		ctx.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}

	user, err := c.Service.GetUser(uid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(404, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"user": user})
}

// CreateUser godoc
// @Summary Create a new user.
// @Description Creates a new user with the provided data.
// @Accept json
// @Produce json
// @Param request body models.UserData true "User data to create"
// @Success 201 {object} models.OKresponse "User created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request body or user already exists"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/users [post]
func (c *Controller) CreateUser(ctx *gin.Context) {
	var userData models.UserData

	err := ctx.ShouldBindJSON(&userData)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	userId, err := c.Service.CreateUser(userData)
	if err != nil {
		if err == models.ErrUserAlreadyExists {
			ctx.JSON(400, gin.H{"error": "User already exists"})
			return
		}
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(201, gin.H{"message ": "User created", "user_id": userId})
}

// UpdateUser godoc
// @Summary Update a user by ID.
// @Description Updates a user with the provided data.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.UserUpdate true "Updated user data"
// @Success 200 {object} models.OKresponse "User updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid user ID or request body"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/users/{id} [put]
func (c *Controller) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	uid, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}
	if uid <= 0 {
		ctx.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}

	var userData models.User

	err = ctx.ShouldBindJSON(&userData)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	err = c.Service.UpdateUser(uid, userData)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(404, gin.H{"error": "User not found"})
			return
		}
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "User updated"})

}

// DeleteUser godoc
// @Summary Delete a user by ID.
// @Description Deletes a user by their ID.
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.OKresponse "User deleted successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid user ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/users/{id} [delete]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	uid, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}
	if uid <= 0 {
		ctx.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}

	err = c.Service.DeleteUser(uid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(404, gin.H{"error": "User not found"})
			return
		}
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "User deleted"})
}

// StartTask godoc
// @Summary Start a task for a user by ID and task ID.
// @Description Starts a task for a user by their IDs.
// @Produce json
// @Param id path int true "User ID"
// @Param taskId path int true "Task ID"
// @Success 200 {object} models.OKresponse "Task started successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid user ID or task ID"
// @Failure 404 {object} models.ErrorResponse "User or task not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/users/{id}/tasks/{taskId}/start [post]
func (c *Controller) StartTask(ctx *gin.Context) {
	userID := ctx.Param("id")
	taskID := ctx.Param("taskId")

	uid, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}
	if uid <= 0 {
		ctx.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}

	tid, err := strconv.Atoi(taskID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid taskID"})
		return
	}
	if tid <= 0 {
		ctx.JSON(400, gin.H{"error": "Invalid taskID"})
		return
	}

	err = c.Service.StartTask(uid, tid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(404, gin.H{"error": "User or task not found"})
			return
		}
		if err == models.ErrTaskAlreadyStarted {
			ctx.JSON(400, gin.H{"error": "Task already started"})
			return
		}
		if err == models.ErrTaskNotFound {
			ctx.JSON(404, gin.H{"error": "Task not found"})
			return
		}
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Task started"})
}

// EndTask godoc
// @Summary End a task for a user by ID and task ID.
// @Description Ends a task for a user by their IDs.
// @Produce json
// @Param id path int true "User ID"
// @Param taskId path int true "Task ID"
// @Success 200 {object} models.OKresponse "Task ended successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid user ID or task ID"
// @Failure 404 {object} models.ErrorResponse "User or task not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/users/{id}/tasks/{taskId}/stop [post]
func (c *Controller) EndTask(ctx *gin.Context) {
	userID := ctx.Param("id")
	taskID := ctx.Param("taskId")

	uid, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}
	if uid <= 0 {
		ctx.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}

	tid, err := strconv.Atoi(taskID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid taskID"})
		return
	}
	if tid <= 0 {
		ctx.JSON(400, gin.H{"error": "Invalid taskID"})
		return
	}

	err = c.Service.EndTask(uid, tid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(404, gin.H{"error": "User or task not found"})
			return
		}
		if err == models.ErrTaskNotStarted {
			ctx.JSON(400, gin.H{"error": "Task not started yet"})
			return
		}
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Task ended"})
}

// Получение нагрузки пользователя по ID пользователя
func (c *Controller) GetUserWorkloadsByUserID(ctx *gin.Context) {
	userID := ctx.Param("id")

	var startDate, endDate time.Time
	var err error

	startDate, err = time.Parse("2006-01-02", ctx.Query("start_date"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid start_date, format should be YYYY-MM-DD"})
		return
	}

	endDate, err = time.Parse("2006-01-02", ctx.Query("end_date"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid end_date, format should be YYYY-MM-DD"})
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}
	if uid <= 0 {
		ctx.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}

	userWorkloads, err := c.Service.GetUserWorkloadsByUserID(uid, startDate, endDate)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(404, gin.H{"error": "User workloads not found"})
			return
		}
		if err == models.ErrStartDateAfterEndDate {
			ctx.JSON(400, gin.H{"error": "start_date should be before end_date"})
			return
		}
		if err == models.ErrStartDateInFuture {
			ctx.JSON(400, gin.H{"error": "start_date should be in the past"})
			return
		}
		log.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"user_workloads": userWorkloads})
}

// GetTasks godoc
// @Summary Get all tasks.
// @Description Retrieves all tasks.
// @Produce json
// @Success 200 {object} models.ResponseTasksList "Successful response with tasks"
// @Failure 404 {object} models.ErrorResponse "Tasks not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/tasks [get]
func (c *Controller) GetTasks(ctx *gin.Context) {
	tasks, err := c.Service.GetTasks()
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(404, gin.H{"error": "Tasks not found"})
			return
		}
		log.Println(err.Error())
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"tasks": tasks})
}
