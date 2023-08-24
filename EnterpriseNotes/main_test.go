package main

import (
	"testing"

	"github.com/pashagolub/pgxmock/v2"
)


func TestListNotes(t *testing.T) {
	// Mock the database connection
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer mock.Close()

	// Create a new expected query and rows
	expectedQuery := "SELECT * FROM notes"
	rows := mock.NewRows([]string{"id", "description"}).
		AddRow(1, "Note 1").
		AddRow(2, "Note 2")

	// Expect the query and return the mocked rows
	mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

	// Call the function you want to test
	err = ListNotes()

	// Check for any errors during
}