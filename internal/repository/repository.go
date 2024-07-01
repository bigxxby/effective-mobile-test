package repository

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/bigxxby/effective-mobile-test/internal/models"
)

type Repository struct {
	DB *sql.DB
}

func New(db *sql.DB) Repository {
	return Repository{
		DB: db,
	}
}

func (r *Repository) GetUsers(filter models.Filter, pagination models.Pagination, sortBy, sortOrder string) ([]models.User, error) {
	query := "SELECT id, passport_number, surname, name FROM users WHERE 1=1"

	var args []interface{}
	argCount := 1

	if filter.PassportNumber != "" {
		query += " AND passport_number = $" + strconv.Itoa(argCount)
		args = append(args, filter.PassportNumber)
		argCount++
	}
	if filter.Surname != "" {
		query += " AND surname = $" + strconv.Itoa(argCount)
		args = append(args, filter.Surname)
		argCount++
	}
	if filter.Name != "" {
		query += " AND name = $" + strconv.Itoa(argCount)
		args = append(args, filter.Name)
		argCount++
	}

	query += " ORDER BY " + sortBy + " " + sortOrder
	query += " LIMIT $" + strconv.Itoa(argCount) + " OFFSET $" + strconv.Itoa(argCount+1)
	args = append(args, pagination.PageSize, (pagination.Page-1)*pagination.PageSize)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.PassportNumber, &user.Surname, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) GetUserWorkloadsByUserID(userID int, startDate, endDate time.Time) ([]models.UserWorkload, error) {
	query := `
		SELECT t.id, t.task_name, 
		       SUM(EXTRACT(EPOCH FROM (t.end_time - t.start_time))) AS total_seconds
		FROM tasks t
		WHERE t.user_id = $1 AND t.start_time >= $2 AND t.end_time <= $3
		GROUP BY t.id, t.task_name
		ORDER BY total_seconds DESC
	`

	rows, err := r.DB.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userWorkloads []models.UserWorkload
	for rows.Next() {
		var userWorkload models.UserWorkload
		var totalTimeSeconds float64
		err := rows.Scan(&userWorkload.TaskID, &userWorkload.TaskName, &totalTimeSeconds)
		if err != nil {
			return nil, err
		}
		userWorkload.TotalHours = int(totalTimeSeconds / 3600)
		userWorkload.TotalMinutes = int((totalTimeSeconds - float64(userWorkload.TotalHours)*3600) / 60)
		userWorkloads = append(userWorkloads, userWorkload)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userWorkloads, nil
}
