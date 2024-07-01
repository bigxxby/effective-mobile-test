package internal

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bigxxby/effective-mobile-test/internal/controller"
	"github.com/bigxxby/effective-mobile-test/internal/repository"
	routes "github.com/bigxxby/effective-mobile-test/internal/router"
	"github.com/bigxxby/effective-mobile-test/internal/service"
	config "github.com/bigxxby/effective-mobile-test/pkg/config"
	"github.com/bigxxby/effective-mobile-test/pkg/migrations"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Run the application

// 1. Initialize the configuration
// 2. Initialize the logger
// 3. Initialize the database
// 4. Initialize the router
// 5. Initialize the server
// 6. Start the server
// 7. Gracefully shutdown the server

// GET /api/users?filter[name]=John&page=1&perPage=10
// GET /api/users/{id}/worklog?startDate=2024-06-01&endDate=2024-06-30
// POST /api/users/{id}/tasks/{taskId}/start
// POST /api/users/{id}/tasks/{taskId}/stop
// DELETE /api/users/{id}
func Run() {
	config.LoadEnv()
	host := config.GetEnv("DB_HOST")
	port := config.GetEnv("DB_PORT")
	user := config.GetEnv("DB_USER")
	password := config.GetEnv("DB_PASSWORD")
	dbname := config.GetEnv("DB_NAME")
	fmt.Println(host, port, user, password, dbname)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	if err := migrations.ApplyMigrations(db); err != nil {
		log.Fatal(err)
	}

	repo := repository.New(db)
	if err != nil {
		log.Println(err)
		return
	}
	service := service.New(repo)

	controller := controller.New(service)

	router := gin.Default()
	router.Use(gin.Logger())
	routes.RegisterRoutes(router, &controller)

	log.Println("Server started on http://localhost:8080")
	err = router.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}
