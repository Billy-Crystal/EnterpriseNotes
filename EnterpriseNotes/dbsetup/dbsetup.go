package dbsetup

import (
	"context"
	"fmt"

	"EnterpriseNotes/db" // Import the package for database connections
)

var DATABASE_URL = "postgres://postgres:postgres@localhost:5432/postgres"

func SetupDatabase() ( *db.PostgresDatabase, error) {
	db, err := db.NewPostgresDatabase(DATABASE_URL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the database: %v", err)
	}

	if err := createNotesTable(db); err != nil {
		return nil, fmt.Errorf("unable to create notes table: %v", err)
	}

	return db, nil
}

func createNotesTable(db *db.PostgresDatabase) error {
	notesTable := `
		CREATE TABLE IF NOT EXISTS notes (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			noteType TEXT NOT NULL,
			description TEXT NOT NULL,
			noteCreated TEXT NOT NULL,
			taskCompletionTime TEXT,
			taskCompletionDate TEXT,
			noteStatus TEXT,
			noteDelegation TEXT,
			sharedUsers TEXT,
			fts_text tsvector
		)
	`
	_, err := db.Conn.Exec(context.Background(), notesTable)
	return err
}
