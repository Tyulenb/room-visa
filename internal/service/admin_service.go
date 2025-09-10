package service

import (
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

    return bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
}
