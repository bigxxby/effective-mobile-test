package router

import (
	"github.com/bigxxby/effective-mobile-test/internal/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, controller *controller.Controller) {
	router.GET("/api/users", controller.GetUsers)
	router.GET("/api/users/:id/workloads", controller.GetUserWorkloadsByUserID)
	// router.POST("/api/users/:id/tasks/:taskId/start", controller.StartTask)
	// router.POST("/api/users/:id/tasks/:taskId/stop", controller.StopTask)
	// router.DELETE("/api/users/:id", controller.DeleteUser)
}

// Получение данных пользователей:
// Фильтрация по всем полям.
// Пагинация.
// Получение трудозатрат по пользователю за период задача-сумма часов и минут с сортировкой от большей затраты к меньшей
// Начать отсчет времени по задаче для пользователя
// Закончить отсчет времени по задаче для пользователя
// Удаление пользователя
// Изменение данных пользователя
// Добавление нового пользователя в формате:
