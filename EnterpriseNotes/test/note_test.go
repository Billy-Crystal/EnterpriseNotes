package test

import (
	"testing"

	"github.com/AlexGithub777/EnterpriseNotes/internal/app" // Import application package.
)

func TestCreateNote(t *testing.T) {
	// Create a note using the CreateNote function
	note := app.CreateNote("Test Title", "Test Content")
	expectedTitle := "Test Title"
	expectedContent := "Test Content"

	// Check if the note's title matches the expected title
	if note.Title != expectedTitle {
		t.Errorf("Expected title: %s, but got: %s", expectedTitle, note.Title)
	}

	// Check if the note's content matches the expected content
	if note.Content != expectedContent {
		t.Errorf("Expected content: %s, but got: %s", expectedContent, note.Content)
	}
}


