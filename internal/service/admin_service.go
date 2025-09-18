package service

import (
	"database/sql"
	"fmt"
	"room-visa/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	ar model.AdminRepository
}

func NewAdminService(ar model.AdminRepository) *AdminService {
	return &AdminService{
		ar: ar,
	}
}

func (as *AdminService) CheckPassword(login, password string) error {
	user, err := as.ar.SelectAdminByLogin(login)
	if err != nil {
		return err
	}
	fmt.Println(password)
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (as *AdminService) CreateAdmin(login, password string) (*model.Admin, error) {
	_, err := as.ar.SelectAdminByLogin(login)
	if err == nil {
		return nil, fmt.Errorf("User with such login already exists")
	} else {
		if err == sql.ErrNoRows {
			hashpass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			return as.ar.InsertAdmin(login, string(hashpass))
		} else {
			return nil, err
		}
	}
}
