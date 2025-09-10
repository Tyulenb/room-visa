package model

type Admin struct {
    Login string `json:"login"`
    Password string `json:"password"`
}

type AdminService interface {
    CheckPassword(login, password string) error

}

type AdminRepository interface {
    SelectAdminByLogin(login string) (Admin, error)
}
