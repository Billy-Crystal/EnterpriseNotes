package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/AlexGithub777/EnterpriseNotes/dbsetup"
)


func main() {
	var err error

	dbInstance, err := dbsetup.SetupDatabase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to the database: %v\n", err)
		os.Exit(1)
	}
	defer dbInstance.Conn.Close(context.Background())

	if len(os.Args) == 1 {
		printHelp()
		os.Exit(0)
	}

    // Switch statement to handle different commands
    switch os.Args[1] {
    case "list":

        err := dbInstance.ListNotes(context.Background())
        if err != nil {
            fmt.Fprintf(os.Stderr, "Unable to list notes: %v\n", err)
            os.Exit(1)
        }

	case "add":
		if len(os.Args) > 5 {
			fmt.Fprintln(os.Stderr, "Missing type or title or description")
			os.Exit(1)
		}

		var noteType, taskCompletionDate, taskCompletionTime string

		if os.Args[2] == "note" {
			noteType = "note"
			taskCompletionDate = "none"
			taskCompletionTime = "none"
		} else if os.Args[2] == "task" {
			noteType = "task"
			fmt.Print("Enter task completion date (dd-mm-yyyy): ")
			fmt.Scan(&taskCompletionDate)

			fmt.Print("Enter task completion time (HH:MM): ")
			fmt.Scan(&taskCompletionTime)
		} else {
			fmt.Fprintln(os.Stderr, "Invalid note type. Please use 'note' or 'task'.")
			os.Exit(1)
		}

		title := os.Args[3]
		description := os.Args[4]
		noteCreated := time.Now()
		formattedNoteCreated := noteCreated.Format(time.ANSIC)
		noteStatus := "none"
		noteDelegation := "none"
		sharedUsers := "none"
		

		err = dbInstance.AddNote(title, noteType, description, formattedNoteCreated, taskCompletionDate, taskCompletionTime, noteStatus, noteDelegation, sharedUsers)
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
		err = dbInstance.UpdateNote(int(n), os.Args[3])
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
		err = dbInstance.RemoveNote(int(n))
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
        err := dbInstance.SearchNotes(searchPattern)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Unable to search notes: %v\n", err)
            os.Exit(1)
        }

	case "find":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Missing note ID and/or text snippet pattern")
			os.Exit(1)
		}
		noteID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid note ID: %v\n", err)
			os.Exit(1)
		}
		snippetPattern := os.Args[3]
	
		count, description, err := dbInstance.FindTextSnippetInNote(noteID, snippetPattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while finding text snippet: %v\n", err)
			os.Exit(1)
		}
	
		analysisCount := dbInstance.AnalyzeTextSnippet(description)
	
		fmt.Printf("Text snippet '%s' found %d times in the note:\n%s\n", snippetPattern, count, description)
		fmt.Printf("Analysis: Text snippet patterns found %d times in the note\n", analysisCount)
	
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
  notes add <type> <title> <description>
  notes search <search-pattern>
  notes find <id> <text-pattern> (only searches over )
  notes update <id> <description>
  notes remove <id>
  

Example:

  notes list
  notes add note shopping oranges
  notes add task jobs "feed cat"
  notes search shopping
  notes search oranges
  notes find 1 agenda
  notes update 1 updated
  notes remove 1
`)
}
