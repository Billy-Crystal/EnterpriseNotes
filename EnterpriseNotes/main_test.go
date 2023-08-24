package main

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

var testConn *pgx.Conn

func setupTest() {
	// Connect to the test database
	connConfig, _ := pgx.ParseConfig("postgres://postgres:postgres@localhost:5432/postgres")
	conn, _ := pgx.ConnectConfig(context.Background(), connConfig)
	testConn = conn
}

func teardownTest() {
	// Close the testConn connection
	testConn.Close(context.Background())
}

func TestListNotes(t *testing.T) {
	setupTest()
	defer teardownTest()

	
	err := ListNotes()
	assert.NoError(t, err) // Ensure no error occurred
}

// ... Other test functions for addNote, updateNote, and removeNote

func TestAddNote (t *testing.T) {
	setupTest()
	defer teardownTest()
	// Insert test data into the test database
	err := AddNote("Note Test")
	if err != nil {
		t.Fatalf("Failed to add note: %v", err)
	}
}