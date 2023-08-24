// db.go
package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Database interface {
	ListNotes(ctx context.Context) error
	AddNote(title, description, noteCreated, noteStatus, noteDelegation, sharedUsers string) error
	RemoveNote(noteID int) error
	SearchNotes(pattern string) error
	UpdateNote(noteID int, description string) error
	// Define other methods as needed
}

type PostgresDatabase struct {
    conn *pgx.Conn // The connection to your PostgreSQL database
}

func (db *PostgresDatabase) ListNotes(ctx context.Context) error {
    query := `
        SELECT id, title, description, noteCreated, noteStatus, noteDelegation, sharedUsers
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
        var description string
        var noteCreated string
        var noteStatus string
        var noteDelegation string
        var sharedUsers string
        err := rows.Scan(&id, &title, &description, &noteCreated, &noteStatus, &noteDelegation, &sharedUsers)
        if err != nil {
            return err
        }
        fmt.Printf("Note ID:%d.\n Title: %s\n Description: %s\n Note Created: %s\n Note Status: %s\n Note Delegation: %s\n Shared users: %s \n",
            id, title, description, noteCreated, noteStatus, noteDelegation, sharedUsers)
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
    query := fmt.Sprintf("SELECT id, title, description, noteCreated, noteStatus, noteDelegation, sharedUsers FROM notes WHERE fts_text @@ to_tsquery('english', $1)")
    rows, err := db.conn.Query(context.Background(), query, pattern)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var title, description, noteCreated, noteStatus, noteDelegation, sharedUsers string
        if err := rows.Scan(&id, &title, &description, &noteCreated, &noteStatus, &noteDelegation, &sharedUsers); err != nil {
            return err
        }
        fmt.Printf("Note ID:%d.\n Title: %s\n Description: %s\n %s\n Note Status: %s\n Note Delegation: %s\n Shared users: %s \n", id, title, description, noteCreated, noteStatus, noteDelegation, sharedUsers)
    }

    return nil
}

func (db *PostgresDatabase) AddNote(title, description, noteCreated, noteStatus, noteDelegation, sharedUsers string) error {
	_, err := db.conn.Exec(
        context.Background(),
        "INSERT INTO notes(title, description, noteCreated, noteStatus, noteDelegation, sharedUsers, fts_text) VALUES($1, $2, $3, $4, $5, $6, to_tsvector('english', $1 || ' ' || $2 || ' ' || $4 || ' ' || $5 || ' ' || $6))",
        title, description, noteCreated, noteStatus, noteDelegation, sharedUsers,
    )
    return err
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