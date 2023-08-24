package main

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)


func TestListNotes(t *testing.T) {
	
	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()

	mockDB.ExpectQuery("SELECT id, title, description, noteCreated, noteStatus, noteDelegation, sharedUsers FROM notes").
		WillReturnRows(pgxmock.NewRows([]string{"id", "title", "description", "noteCreated", "noteStatus", "noteDelegation", "sharedUsers"}).
			AddRow(1, "Title 1", "Description 1", "2023-08-01", "Status 1", "Delegation 1", "User 1").
			AddRow(2, "Title 2", "Description 2", "2023-08-02", "Status 2", "Delegation 2", "User 2"))



	


	conn, err := pgxmock.NewPool()
	if err != nil {
        return 
    }



	db := &MockDatabase{mockDB: conn.AsConn()}

	// Call the ListNotes function and check the results
	err = db.ListNotes(context.Background())
	assert.NoError(t, err)

	// Check if expectations were met
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}


}