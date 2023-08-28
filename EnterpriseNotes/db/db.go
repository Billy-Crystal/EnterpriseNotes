// db.go
package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Database interface {
	ListNotes(ctx context.Context) error
	AddNote(title, noteType, description, noteCreated, taskCompletionDate, taskCompletionTime, noteStatus, noteDelegation, sharedUsers string) error
	RemoveNote(noteID int) error
	SearchNotes(pattern string) error
	UpdateNote(noteID int, description string) error
    AnalyzeTextSnippet(description string) error

	// Define other methods as needed
}

type PostgresDatabase struct {
    conn *pgx.Conn // The connection to your PostgreSQL database
}

func (db *PostgresDatabase) ListNotes(ctx context.Context) error {
    query := `
        SELECT id, title, noteType, description, noteCreated, taskCompletionDate, taskCompletionTime, noteStatus, noteDelegation, sharedUsers
        FROM notes
    `
    rows, err := db.conn.Query(ctx, query)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var title string
        var noteType string
        var description string
        var noteCreated string
        var taskCompletionTime string
		var taskCompletionDate string
        var noteStatus string
        var noteDelegation string
        var sharedUsers string
        err := rows.Scan(&id, &title, &noteType, &description, &noteCreated, &taskCompletionDate, &taskCompletionTime, &noteStatus, &noteDelegation, &sharedUsers)
        if err != nil {
            return err
        }
        fmt.Printf(" Note ID: %d.\n Title: %s\n Note Type: %s\n Description: %s\n Note Created: %s\n Task Completion Date: %s\n Task Completion Time: %s\n Note Status: %s\n Note Delegation: %s\n Shared users: %s \n",
            id, title, noteType, description, noteCreated, taskCompletionDate, taskCompletionTime, noteStatus, noteDelegation, sharedUsers)
    }

    return rows.Err()
}

// Implement the methods of the Database interface
func (db *PostgresDatabase) UpdateNote(noteID int, description string) error {
    _, err := db.conn.Exec(context.Background(), "UPDATE notes SET description=$1, fts_text=to_tsvector('english', $1) WHERE id=$2", description, noteID)
	return err
}

func (db *PostgresDatabase) RemoveNote(noteID int) error {
    _, err := db.conn.Exec(context.Background(), "DELETE FROM notes WHERE id=$1", noteID)
	return err
}

func (db *PostgresDatabase) SearchNotes(pattern string) error {
    query := fmt.Sprintf("SELECT id, title, noteType, description, noteCreated, taskCompletionDate, taskCompletionTime, noteStatus, noteDelegation, sharedUsers FROM notes WHERE fts_text @@ to_tsquery('english', $1)")
    rows, err := db.conn.Query(context.Background(), query, pattern)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var title, noteType, description, noteCreated, taskCompletionDate, taskCompletionTime, noteStatus, noteDelegation, sharedUsers string
        if err := rows.Scan(&id, &title, &noteType, &description, &noteCreated, &taskCompletionDate, &taskCompletionTime, &noteStatus, &noteDelegation, &sharedUsers); err != nil {
            return err
        }
        fmt.Printf(" Note ID: %d.\n Title: %s\n Note Type: %s\n Description: %s\n Note Created: %s\n Task Completion Date: %s\n Task Completion Time: %s\n Note Status: %s\n Note Delegation: %s\n Shared users: %s \n",
            id, title, noteType, description, noteCreated, taskCompletionDate, taskCompletionTime, noteStatus, noteDelegation, sharedUsers)
    }

    return nil
}



func (db *PostgresDatabase) AddNote(title, noteType, description, noteCreated, taskCompletionDate, taskCompletionTime, noteStatus, noteDelegation, sharedUsers string) error {
    _, err := db.conn.Exec(
        context.Background(),
        "INSERT INTO notes(title, noteType, description, noteCreated, taskCompletionDate, taskCompletionTime, noteStatus, noteDelegation, sharedUsers, fts_text) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, to_tsvector('english', $1 || ' ' || $2 || ' ' || $3 || ' ' || $4 || ' ' || $5 || ' ' || $6 || ' ' || $7 || ' ' || $8 || ' ' || $9))",
        title, noteType, description, noteCreated, taskCompletionDate, taskCompletionTime, noteStatus, noteDelegation, sharedUsers,
    )
    return err
}






func (db *PostgresDatabase) FindTextSnippetInNote(noteID int, snippetPattern string) (int, string, error) {
    var description string
    row := db.conn.QueryRow(context.Background(), "SELECT description FROM notes WHERE id=$1", noteID)
    if err := row.Scan(&description); err != nil {
        return 0, "", err
    }

    count := countOccurrences(description, snippetPattern)

    return count, description, nil
}



func countOccurrences(text, snippetPattern string) int {
    // Convert both the text and snippetPattern to lowercase for case-insensitive matching
    lowerText := strings.ToLower(text)
    lowerSnippetPattern := strings.ToLower(snippetPattern)

    // Count the occurrences of the snippetPattern in the text
    count := strings.Count(lowerText, lowerSnippetPattern)

    return count
}

func (db *PostgresDatabase) AnalyzeTextSnippet(description string) int {
    analysisCount := 0

    // Count occurrences of sentences with a given prefix and/or suffix
    prefix := "Dear"
    suffix := "Sincerely"
    sentencePattern := fmt.Sprintf("%s.*%s", prefix, suffix)
    sentenceCount := countOccurrences(description, sentencePattern)

    // Count occurrences of phone numbers with a specific area code
    areaCode := "555"
    phonePattern := fmt.Sprintf(`\b%s-\d{4}\b`, areaCode)
    phoneCount := countOccurrences(description, phonePattern)

    // Count occurrences of email addresses with a partial domain
    domain := "example.com"
    emailPattern := fmt.Sprintf(`\b\w+@%s\b`, domain)
    emailCount := countOccurrences(description, emailPattern)

    // Count occurrences of specific keywords
    keywords := []string{"meeting", "minutes", "agenda", "action", "attendees", "apologies"}
    keywordCount := countOccurrencesByKeywords(description, keywords)

    // Count occurrences of words in all capitals of three characters or more
    capitalWordPattern := `\b[A-Z]{3,}\b`
    capitalWordCount := countOccurrences(description, capitalWordPattern)

    // Sum up all the analysis counts
    analysisCount = sentenceCount + phoneCount + emailCount + keywordCount + capitalWordCount

    return analysisCount
}

func countOccurrencesByKeywords(text string, keywords []string) int {
    count := 0
    lowerText := strings.ToLower(text)

    for _, keyword := range keywords {
        lowerKeyword := strings.ToLower(keyword)
        count += strings.Count(lowerText, lowerKeyword)
    }

    return count
}




func NewPostgresDatabase(connectionString string) (*PostgresDatabase, error) {
    conn, err := pgx.Connect(context.Background(), connectionString)
    if err != nil {
        return nil, err
    }

    return &PostgresDatabase{
        conn: conn,
    }, nil
}