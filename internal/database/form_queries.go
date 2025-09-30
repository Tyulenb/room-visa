package database

import (
	"database/sql"
	"room-visa/internal/model"

	"github.com/google/uuid"
)

type FormRepository struct {
    db *sql.DB
}

func NewFormRepository(db *sql.DB) *FormRepository {
    return &FormRepository{
        db: db,
    }
}

func (fr *FormRepository) InsertForm(req *model.Request, data *model.RequestData) error {
    queryReg := `INSERT INTO request VALUES($1,$2,$3,$4)` 
    queryData := `
        INSERT INTO request_data (
            request_id, name, surname, sex, 
            ethnicity, citizenship, purpose, photopath
            ) 
        VALUES($1,$2,$3,$4,$5,$6,$7,$8)
    `

    tx, err := fr.db.Begin()
    if err != nil {
        return err
    }

    _, err = tx.Exec(queryReg, req.Id, req.Status, req.Created_at, req.Updated_at)
    if err != nil {
        tx.Rollback()
        return err
    }

    _, err = tx.Exec(
        queryData, 
        data.Request_id, data.Name, data.Surname, data.Sex,
        data.Ethnicity, data.Citizenship, data.Purpose, data.Photo,
    )
    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit() 
}

func (fr *FormRepository) selectFromRequests(query string, args ...any) ([]model.Request, error){
    rows, err := fr.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    requests := make([]model.Request, 0)
    for rows.Next() {
        request := new(model.Request)
        err := rows.Scan(&request.Id, &request.Status, &request.Created_at, &request.Updated_at)
        if err != nil {
            return nil, err
        }
        requests = append(requests, *request)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return requests, nil
}

//SELECTS all data from table requests
func (fr *FormRepository) SelectRequests() ([]model.Request, error) {
    query := `SELECT * FROM request`
    return fr.selectFromRequests(query)
}

func (fr *FormRepository) SelectAwaitingRequests() ([]model.Request, error) {
    query := `SELECT * FROM requests WHERE status = $1`
    return fr.selectFromRequests(query, "Awaiting")
}

func (fr *FormRepository) SelectRequestDataByRequest(req uuid.UUID) (*model.RequestData, error) {
    query := `SELECT * FROM request_data WHERE request_id = $1`

    row := fr.db.QueryRow(query, req)

    if err := row.Err(); err != nil {
        return nil, err
    }
    
    data := new(model.RequestData)
    err := row.Scan(
        &data.Id, &data.Request_id, &data.Name, &data.Surname,
        &data.Sex, &data.Ethnicity, &data.Citizenship,
        &data.Purpose, &data.Photo,
    )
    if err != nil {
        return nil, err
    }

    return data, nil
}

