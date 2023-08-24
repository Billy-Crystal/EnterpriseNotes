package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
)

// DATABASE_URL should be set with your PostgreSQL database connection details
var DATABASE_URL = "postgres://postgres:postgres@localhost:5432/postgres"


func main() {
	var err error

	db, err := NewPostgresDatabase(DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to the database: %v\n", err)
		os.Exit(1)
	}
	defer db.conn.Close(context.Background())

	// Create table called notes
	notesTable := `
	CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		noteCreated TEXT NOT NULL,
		noteStatus TEXT,
		noteDelegation TEXT,
		sharedUsers TEXT,
		fts_text tsvector
	)
`
	_, err = db.conn.Exec(context.Background(), notesTable)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute SQL command: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		printHelp()
		os.Exit(0)
	}

	// Create a new instance of PostgresDatabase
    database := &PostgresDatabase{conn: db.conn}

    // Switch statement to handle different commands
    switch os.Args[1] {
    case "list":
        err := database.ListNotes(context.Background())
        if err != nil {
            fmt.Fprintf(os.Stderr, "Unable to list notes: %v\n", err)
            os.Exit(1)
        }

	case "add":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Missing title and/or description")
			os.Exit(1)
		}
		title := os.Args[2]
		description := os.Args[3]
	
		// You can set the default noteCreated and noteStatus values as needed
		noteCreated := time.Now()
		

		// Format noteCreated into a string
		formattedNoteCreated := noteCreated.Format(time.ANSIC)
		//fmt.Println(formattedNoteCreated)

		noteStatus := "none"
		noteDelegation := "none"
		sharedUsers := "none"


	
		err = database.AddNote(title, description, formattedNoteCreated, noteStatus, noteDelegation, sharedUsers)
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
		err = database.UpdateNote(int(n), os.Args[3])
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
		err = database.RemoveNote(int(n))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to remove note: %v\n", err)
			os.Exit(1)
		}

	case "search":
        if len(os.Args) < 3 {
            fmt.Fprintln(os.Stderr, "Missing search pattern")
            os.Exit(1)
        }
        searchPattern := os.Args[2]

        // Perform full-text search
        err := database.SearchNotes(searchPattern)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Unable to search notes: %v\n", err)
            os.Exit(1)
        }

	default:
		fmt.Fprintln(os.Stderr, "Invalid command")
		printHelp()
		os.Exit(1)
	}
}


func printHelp() {
	fmt.Print(`
	Enterprise Notes

Usage:

  notes list
  notes add <title> <description>
  notes search <search-pattern>
  notes update <id> <description>
  notes remove <id>
  

Example:

  notes list
  notes add shopping oranges
  notes add "Note title" "This is a note"
  notes search shopping
  notes search oranges
  notes update 1 updated
  notes remove 1
`)
}
