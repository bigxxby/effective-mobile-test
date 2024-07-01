package controller

import (
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
	users, err := c.Service.GetUsers()
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"users": users,
	})

}
