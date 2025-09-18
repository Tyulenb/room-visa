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

func (a *AdminRepository) singleRowQuery(query string, val ...any) (*model.Admin, error) {
	row := a.db.QueryRow(query, val...)
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

func (a *AdminRepository) SelectAdminByLogin(login string) (*model.Admin, error) {
	query := `SELECT * FROM admin WHERE login=$1`
	return a.singleRowQuery(query, login)
}

func (a *AdminRepository) InsertAdmin(login, hashpass string) (*model.Admin, error) {
	query := `INSERT INTO admin (login, hashpass) VALUES($1, $2) RETURNING *`
	return a.singleRowQuery(query, login, hashpass)
}
