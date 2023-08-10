package app

// Note represents a single note in the application.
type Note struct {
    ID      int    // Unique identifier for the note.
    Title   string // Title of the note.
    Content string // Content of the note.
    // Add other fields as needed.
}

var notes []Note

// CreateNote adds a new note to the list of notes.
func CreateNote(title, content string) Note {
    newID := len(notes) + 1
    note := Note{ID: newID, Title: title, Content: content}
    notes = append(notes, note)
    return note
}

// GetNoteByID retrieves a note by its ID.
func GetNoteByID(id int) (Note, bool) {
    for _, note := range notes {
        if note.ID == id {
            return note, true
        }
    }
    return Note{}, false
}

// UpdateNote updates an existing note with new data.
func UpdateNote(id int, title, content string) bool {
    for i := range notes {
        if notes[i].ID == id {
            notes[i].Title = title
            notes[i].Content = content
            return true
        }
    }
    return false
}

// DeleteNote deletes a note by its ID.
func DeleteNote(id int) bool {
    for i, note := range notes {
        if note.ID == id {
            // Remove the note from the slice.
            notes = append(notes[:i], notes[i+1:]...)
            return true
        }
    }
    return false


    
}