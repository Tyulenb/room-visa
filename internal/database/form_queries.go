package database

import (
	"database/sql"
	"room-visa/internal/model"
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
