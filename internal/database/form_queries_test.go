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
    
    mock.ExpectQuery("SELECT \\* FROM request$").WithoutArgs().WillReturnRows(rows)

    fr := NewFormRepository(db)
    reqs, err := fr.SelectRequests()
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("Unmet expectations during queries: %v", err)
    }

    if len(reqs) != len(rowsData) {
        t.Fatalf("Expected %d requests, got %d", len(rowsData), len(reqs))
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

func TestSelectRequestDataByRequest(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatal("Failed to open sqlmock: ", err)
    }
    defer db.Close()

    requestData := []model.RequestData {
        {
           Id: 1, Request_id: uuid.New(), Name: "Ivan", Surname: "Ivanov",
           Sex: "male", Ethnicity: "Russian", Citizenship: "Russian",
           Purpose: "Business Trip", Photo: uuid.NewString(),
        },
        {
           Id: 2, Request_id: uuid.New(), Name: "John", Surname: "Doe",
           Sex: "male", Ethnicity: "American", Citizenship: "USA",
           Purpose: "Tourism", Photo: uuid.NewString(),
        },
    }

    columns := []string{
        "id", "request_id", "name", "surname", "sex",
        "ethnicity", "citizenship", "purpose", "photopath",
    }
    rows := []*sqlmock.Rows {
        mock.NewRows(columns),
        mock.NewRows(columns),
    }

    for i, v := range requestData {
        rows[i].AddRow(
            v.Id, v.Request_id, v.Name, v.Surname, v.Sex,
            v.Ethnicity, v.Citizenship, v.Purpose, v.Photo,
        )
        mock.ExpectQuery("SELECT \\* FROM request_data WHERE request_id = \\$1$").
        WithArgs(v.Request_id).WillReturnRows(rows[i])
    }

    fr := NewFormRepository(db)

    for _, v := range requestData {
        data, err := fr.SelectRequestDataByRequest(v.Request_id)
        if err != nil {
            t.Fatalf("Unexpected error: %v", err)
        }
        if v.Id != data.Id{
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", v.Id, data.Id)
        }
        if v.Request_id != data.Request_id{
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", v.Request_id, data.Request_id)
        }
        if v.Name != data.Name {
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", v.Name, data.Name)
        }
        if v.Surname != data.Surname {
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", v.Surname, data.Surname)
        }
        if v.Sex != data.Sex {
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", v.Sex, data.Sex)
        }
        if v.Ethnicity != data.Ethnicity {
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", v.Ethnicity, data.Ethnicity)
        }
        if v.Citizenship != data.Citizenship {
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", v.Citizenship, data.Citizenship)
        }
        if v.Purpose != data.Purpose {
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", v.Purpose, data.Purpose)
        }
        if v.Photo != data.Photo {
            t.Errorf("Unexpected error during values comparising\nExpected: %v, but got: %v", v.Photo, data.Photo)
        }
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("Unmet expectations during queries: %v", err)
    }
}
