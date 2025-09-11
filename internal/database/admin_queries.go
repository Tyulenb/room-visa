package database

import (
	"database/sql"
	"room-visa/internal/model"
)

type AdminRepository struct {
    db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
    return &AdminRepository{
        db: db,
    }
}

func (a *AdminRepository) SelectAdminByLogin(login string) (*model.Admin, error) {
    query := `SELECT * FROM admin WHERE login=$1`    
    row := a.db.QueryRow(query, login)
    if err := row.Err(); err != nil {
        return nil, err
    }

    admin := new(model.Admin)
    err := row.Scan(&admin.Login, &admin.Password)
    if err != nil {
        return nil, err
    }

    return admin, nil 
}
