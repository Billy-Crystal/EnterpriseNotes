package main

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

func setup() (*pgx.Conn, error) {
	// Create a connection to the test database
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/test_database")
	if err != nil {
		return nil, err
	}

	// Create the 'notes' table (if it doesn't exist)
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS notes (id SERIAL PRIMARY KEY, description TEXT NOT NULL)")
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func TestListNotes(t *testing.T) {
	// Set up the test database connection
	conn, err := setup()
	if err != nil {
		t.Fatalf("Failed to set up the test database: %v", err)
	}
	defer conn.Close(context.Background())

	// Insert some test data into the 'notes' table
	_, err = conn.Exec(context.Background(), "INSERT INTO notes(description) VALUES('Note 1')")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Call the 'listNotes' function
	err = EnterpriseNotes.listNotes(conn)
	if err != nil {
		t.Fatalf("Failed to list notes: %v", err)
	}
}

func TestAddNote(t *testing.T) {
	// Set up the test database connection
	conn, err := setup()
	if err != nil {
		t.Fatalf("Failed to set up the test database: %v", err)
	}
	defer conn.Close(context.Background())

	// Call the 'addNote' function
	err = EnterpriseNotes.addNote(conn, "New note")
	if err != nil {
		t.Fatalf("Failed to add note: %v", err)
	}
}

func TestUpdateNote(t *testing.T) {
	// Set up the test database connection
	conn, err := setup()
	if err != nil {
		t.Fatalf("Failed to set up the test database: %v", err)
	}
	defer conn.Close(context.Background())

	// Insert a test note into the 'notes' table
	_, err = conn.Exec(context.Background(), "INSERT INTO notes(description) VALUES('Old description')")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Call the 'updateNote' function
	err = EnterpriseNotes.updateNote(conn, 1, "New description")
	if err != nil {
		t.Fatalf("Failed to update note: %v", err)
	}
}

func TestRemoveNote(t *testing.T) {
	// Set up the test database connection
	conn, err := setup()
	if err != nil {
		t.Fatalf("Failed to set up the test database: %v", err)
	}
	defer conn.Close(context.Background())

	// Insert a test note into the 'notes' table
	_, err = conn.Exec(context.Background(), "INSERT INTO notes(description) VALUES('Note to be deleted')")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Call the 'removeNote' function
	err = EnterpriseNotes.removeNote(conn, 1)
	if err != nil {
		t.Fatalf("Failed to remove note: %v", err)
	}
}
