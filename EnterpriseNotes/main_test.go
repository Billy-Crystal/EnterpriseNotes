package main

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
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
	// Create a test context
	

	// Perform your test
	err := ListNotes(context.Background())
	if err != nil {
		t.Errorf("ListNotes returned an error: %v", err)
	}
}

