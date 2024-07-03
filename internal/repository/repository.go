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

// Получение всех пользователей
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

// Получение рабочей нагрузки пользователей по их ID
func (r *Repository) GetUserWorkloadsByUserID(userID int, startDate, endDate time.Time) ([]models.UserWorkload, error) {
	query := `
		SELECT l.task_id, t.task_name, 
		       EXTRACT(EPOCH FROM (l.end_time - l.start_time)) AS total_seconds
		FROM task_logs l
		INNER JOIN tasks t ON l.task_id = t.id
		WHERE l.user_id = $1 AND l.start_time >= $2 AND l.end_time <= $3
		ORDER BY l.task_id, l.start_time
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

// Запуск задачи
func (r *Repository) StartTask(userID, taskID int) error {
	query := `
		INSERT INTO task_logs (user_id, task_id, start_time)
		VALUES ($1, $2, NOW())
	`
	_, err := r.DB.Exec(query, userID, taskID)
	if err != nil {
		return err
	}

	return nil
}
func (r *Repository) IsTaskInProgress(userID, taskID int) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM task_logs
		WHERE user_id = $1 AND task_id = $2 AND end_time IS NULL
	`
	var count int
	err := r.DB.QueryRow(query, userID, taskID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) IsTaskExists(taskID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM tasks
			WHERE id = $1
		)
	`
	var exists bool
	err := r.DB.QueryRow(query, taskID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Завершение задачи
func (r *Repository) EndTask(userID, taskID int) error {
	query := `
		UPDATE task_logs
		SET end_time = NOW()
		WHERE user_id = $1 AND task_id = $2 AND end_time IS NULL
	`
	res, err := r.DB.Exec(query, userID, taskID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrTaskNotStarted
	}
	return nil
}

// Получение всех задач
func (r *Repository) GetTasks() ([]models.Task, error) {
	query := `
		SELECT id, task_name
		FROM tasks
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Name)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) GetUser(userID int) (models.User, error) {
	query := `
		SELECT id, passport_number, surname, name
		FROM users
		WHERE id = $1
	`
	var user models.User
	err := r.DB.QueryRow(query, userID).Scan(&user.ID, &user.PassportNumber, &user.Surname, &user.Name)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *Repository) UserExistsByPassportNumber(passportNumber string) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM users
            WHERE passport_number = $1
        )
    `
	var exists bool
	err := r.DB.QueryRow(query, passportNumber).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) CreateUser(user models.UserData) (int, error) {
	query := `
		INSERT INTO users (passport_number)
		VALUES ($1)
		RETURNING id
	`
	var id int
	err := r.DB.QueryRow(query, user.PassportNumber).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
func (r *Repository) DeleteUser(userID int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := r.DB.Exec(query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateUser(userID int, user models.User) error {
	query := `
		UPDATE users
		SET  surname = $2, name = $3
		WHERE id = $1
	`
	_, err := r.DB.Exec(query, userID, user.Surname, user.Name)
	if err != nil {
		return err
	}

	return nil
}
