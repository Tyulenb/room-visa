package database

import (
	"room-visa/internal/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestSelectRequests(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatal("Failed to open sqlmock: ", err)
    }
    defer db.Close()

    rowsData := []model.Request{
        {
            Id: uuid.New(),
            Status: "Awaiting",
            Created_at: time.Now(),
            Updated_at: time.Now().Add(time.Hour),
        },
        {
            Id: uuid.New(),
            Status: "Rejected",
            Created_at: time.Now().Add(-2*time.Hour),
            Updated_at: time.Now().Add(-1*time.Hour),
        },
    }

    rows := sqlmock.NewRows([]string{"id", "status", "created_at", "updated_at"})
    for _, rd := range rowsData {
        rows.AddRow(rd.Id, rd.Status, rd.Created_at, rd.Updated_at)
    }
    
    mock.ExpectQuery("SELECT \\* FROM request$").WillReturnRows(rows)

    fr := NewFormRepository(db)
    reqs, err := fr.SelectRequests()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations during queries: %v", err)
    }

    if len(reqs) != len(rowsData) {
        t.Fatalf("expected %d requests, got %d", len(rowsData), len(reqs))
    }

    for i := range reqs {
        if reqs[i].Id != rowsData[i].Id{
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", rowsData[i].Id, reqs[i].Id)
        }
        if reqs[i].Status != rowsData[i].Status{
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", rowsData[i].Status, reqs[i].Status)
        }
        if reqs[i].Created_at != rowsData[i].Created_at{
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", rowsData[i].Created_at, reqs[i].Created_at)
        }
        if reqs[i].Updated_at != rowsData[i].Updated_at{
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", rowsData[i].Updated_at, reqs[i].Updated_at)
        }
    }
    
}
