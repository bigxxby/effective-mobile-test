package router

import (
	_ "github.com/bigxxby/effective-mobile-test/docs"
	"github.com/bigxxby/effective-mobile-test/internal/controller"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(router *gin.Engine, controller *controller.Controller) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/api/users", controller.GetUsers)
	router.GET("/api/users/:id", controller.GetUser)
	router.POST("/api/users", controller.CreateUser)
	router.PUT("/api/users/:id", controller.UpdateUser)

	router.GET("/api/users/:id/workloads", controller.GetUserWorkloadsByUserID)

	router.POST("/api/users/:id/tasks/:taskId/start", controller.StartTask)
	router.POST("/api/users/:id/tasks/:taskId/stop", controller.EndTask)
	router.GET("/api/tasks", controller.GetTasks)

	router.DELETE("/api/users/:id", controller.DeleteUser)
}
