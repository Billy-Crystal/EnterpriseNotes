package test

import (
	"testing"

	"github.com/AlexGithub777/EnterpriseNotes/internal/app"
)

func TestCreateNote(t *testing.T) {
	// Create sample notes using the CreateNote function
	note1 := app.CreateNote("Note 1 Title", "Note 1 Content")
	note2 := app.CreateNote("Note 2 Title", "Note 2 Content")

	// Test that noteIds are automatically assigned
	if note1.ID == note2.ID {
		t.Errorf("Expected different noteIds for note1 and note2, but got the same: %d", note1.ID)
	}


	// Add assertions to verify the properties of the created notes
	if note1.Title != "Note 1 Title" {
		t.Errorf("Expected title: %s, but got: %s", "Note 1 Title", note1.Title)
	}
	if note1.Content != "Note 1 Content" {
		t.Errorf("Expected content: %s, but got: %s", "Note 1 Content", note1.Content)
	}
	// ... Add assertions for other properties if needed ...

	// Similarly, add assertions for the properties of note2
	if note2.Title != "Note 2 Title" {
		t.Errorf("Expected title: %s, but got: %s", "Note 2 Title", note2.Title)
	}
	if note2.Content != "Note 2 Content" {
		t.Errorf("Expected content: %s, but got: %s", "Note 2 Content", note2.Content)
	}

	
}

func TestGetNoteByID(t *testing.T) {
	// Create sample notes using the CreateNote function
	note1 := app.CreateNote("Note 1 Title", "Note 1 Content")

	// Test getting a note by its ID
	note, found := app.GetNoteByID(note1.ID)
	if !found {
		t.Error("Expected note to be found, but it wasn't")
	}
	if note.Title != "Note 1 Title" {
		t.Errorf("Expected title: %s, but got: %s", "Note 1 Title", note.Title)
	}
	// ... Add assertions for other properties if needed ...
}

func TestUpdateNote(t *testing.T) {
	// Create a sample note using the CreateNote function
	note := app.CreateNote("Note Title", "Note Content")

	// Test updating a note's title and content
	updated := app.UpdateNote(note.ID, "Updated Title", "Updated Content")
	if !updated {
		t.Error("Expected note to be updated, but it wasn't")
	}

	// ... Add assertions to verify the updated properties ...
}

func TestDeleteNote(t *testing.T) {
	// Create a sample note using the CreateNote function
	note := app.CreateNote("Note Title", "Note Content")

	// Test deleting a note
	deleted := app.DeleteNote(note.ID)
	if !deleted {
		t.Error("Expected note to be deleted, but it wasn't")
	}

	// Add assertions to verify that the note is indeed deleted
	_, found := app.GetNoteByID(note.ID)
	if found {
		t.Error("Expected note not to be found after deletion, but it was")
	}
}
