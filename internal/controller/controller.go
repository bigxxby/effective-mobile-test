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

func (c *Controller) GetUsers(ctx *gin.Context) {
	var filter models.Filter
	var pagination models.Pagination

	filter.PassportNumber = ctx.Query("passport_number")
	filter.Surname = ctx.Query("surname")
	filter.Name = ctx.Query("name")

	//////////////////////////// basic validation
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
	//////////////////////////// basic validation

	pagination.Page = page
	pagination.PageSize = pageSize

	users, err := c.Service.GetUsers(filter, pagination, sortBy, sortOrder)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(404, gin.H{
				"error": "Users not found",
			})
			return
		}
		ctx.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"users": users,
	})
}

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
