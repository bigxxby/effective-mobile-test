package router

import (
	"github.com/bigxxby/effective-mobile-test/internal/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, controller *controller.Controller) {
	router.GET("/api/users", controller.GetUsers)
	// router.GET("/api/users", controller.GetUsers)
	// router.GET("/api/users/:id/worklog", controller.GetWorklog)
	// router.POST("/api/users/:id/tasks/:taskId/start", controller.StartTask)
	// router.POST("/api/users/:id/tasks/:taskId/stop", controller.StopTask)
	// router.DELETE("/api/users/:id", controller.DeleteUser)
}
