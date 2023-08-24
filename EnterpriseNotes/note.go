package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

// DATABASE_URL should be set with your PostgreSQL database connection details
var DATABASE_URL = "postgres://postgres:postgres@localhost:5432/postgres"
var conn *pgx.Conn

func main() {
	var err error

	conn, err = pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to the database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Create table called notes
	notesTable := `
	CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		description TEXT NOT NULL
	)
`
	_, err = conn.Exec(context.Background(), notesTable)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute SQL command: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "list":
		err = listNotes()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to list notes: %v\n", err)
			os.Exit(1)
		}

	case "add":
		err = addNote(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to add note: %v\n", err)
			os.Exit(1)
		}

	case "update":
		n, err := strconv.ParseInt(os.Args[2], 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to convert note_id into int32: %v\n", err)
			os.Exit(1)
		}
		err = updateNote(int(n), os.Args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to update note: %v\n", err)
			os.Exit(1)
		}

	case "remove":
		n, err := strconv.ParseInt(os.Args[2], 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to convert note_id into int32: %v\n", err)
			os.Exit(1)
		}
		err = removeNote(int(n))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to remove note: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid command")
		printHelp()
		os.Exit(1)
	}
}

func listNotes() error {
	rows, _ := conn.Query(context.Background(), "SELECT * FROM notes")

	for rows.Next() {
		var id int
		var description string
		err := rows.Scan(&id, &description)
		if err != nil {
			return err
		}
		fmt.Printf("%d. %s\n", id, description)
	}

	return rows.Err()
}

func addNote(description string) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO notes(description) VALUES($1)", description)
	return err
}

func updateNote(noteID int, description string) error {
	_, err := conn.Exec(context.Background(), "UPDATE notes SET description=$1 WHERE id=$2", description, noteID)
	return err
}

func removeNote(noteID int) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM notes WHERE id=$1", noteID)
	return err
}

func printHelp() {
	fmt.Print(`Notes pgx demo

Usage:

  notes list
  notes add description
  notes update note_id description
  notes remove note_id

Example:

  notes add 'Important note'
  notes list
`)
}
