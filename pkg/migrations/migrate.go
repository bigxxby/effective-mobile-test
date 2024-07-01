package migrations

import (
	"database/sql"
	"fmt"
	"os"
)

func ApplyMigrations(db *sql.DB) error {
	migrationFiles := []string{
		"pkg/migrations/sql/drop.sql",
		"pkg/migrations/sql/table.sql",
		"pkg/migrations/sql/mock.sql",
	}

	for _, file := range migrationFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", file, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %v", file, err)
		}

		fmt.Printf("Applied migration %s\n", file)
	}

	fmt.Println("All migrations applied successfully")
	return nil
}
