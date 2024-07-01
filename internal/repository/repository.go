package repository

import "database/sql"

type Repository struct {
	DB *sql.DB
}

func New(db *sql.DB) Repository {
	return Repository{
		DB: db,
	}
}

func (r *Repository) GetUsers() ([]string, error) {
	rows, err := r.DB.Query("SELECT name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
