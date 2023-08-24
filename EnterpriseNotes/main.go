package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

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
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		noteCreated TEXT NOT NULL,
		noteStatus TEXT,
		noteDelegation TEXT,
		sharedUsers TEXT
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
		err = ListNotes()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to list notes: %v\n", err)
			os.Exit(1)
		}

	case "add":
		var title, description string

		scanner := bufio.NewScanner(os.Stdin)
	
		// Prompt the user to enter the title
		fmt.Print("Enter title: ")
		scanner.Scan()
        title = scanner.Text()
	
		// Prompt the user to enter the description
		fmt.Print("Enter description: ")
		scanner.Scan()
        description = scanner.Text()
	
		// You can set the default noteCreated and noteStatus values as needed
		noteCreated := time.Now()
		

		// Format noteCreated into a string
		formattedNoteCreated := noteCreated.Format(time.ANSIC)
		//fmt.Println(formattedNoteCreated)

		noteStatus := "none"
		noteDelegation := "none"
		sharedUsers := "none"


	
		err = AddNote(title, description, formattedNoteCreated, noteStatus, noteDelegation, sharedUsers)
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

func ListNotes() error {
	rows, _ := conn.Query(context.Background(), "SELECT * FROM notes")

	for rows.Next() {
		var id int
		var title string
		var description string
		var noteCreated string
		var noteStatus string
		var noteDelegation string
		var sharedUsers string
		err := rows.Scan(&id, &title, &description, &noteCreated, &noteStatus, &noteDelegation, &sharedUsers)
		if err != nil {
			return err
		}
		fmt.Printf("Note ID:%d.\n Title: %s\n Description: %s\n %s\n Note Status: %s\n Note Delegation: %s\n Shared users: %s \n", id, title, description, noteCreated, noteStatus, noteDelegation, sharedUsers)
	}

	return rows.Err()
}

func AddNote(title string, description string, noteCreated string, noteStatus string, noteDelegation string, sharedUsers string) error {
	_, err := conn.Exec(
		context.Background(),
		"INSERT INTO notes(title, description, noteCreated, noteStatus, noteDelegation, sharedUsers) VALUES($1, $2, $3, $4, $5, $6)",
		title, description, noteCreated, noteStatus, noteDelegation, sharedUsers,
	)
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
