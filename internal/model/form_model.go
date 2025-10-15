package model

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

// Form from web page
type Form struct {
	Name        string
	Surname     string
	Sex         string
	Ethnicity   string
	Citizenship string
	Purpose     string
	Photo       multipart.File
}

// Definition of table Request
// Used for storage crucial data about request
type Request struct {
	Id         uuid.UUID
	Status     string
	Created_at time.Time
	Updated_at time.Time
}

// Definition of table request_data
// Contains all data from form
type RequestData struct {
	Id          int
	Request_id  uuid.UUID
	Name        string
	Surname     string
	Sex         string
	Ethnicity   string
	Citizenship string
	Purpose     string
	Photo       string //filepath
}

// Definition of table visa
// Conatains all non-expired visas
// All data contains in token
type Visa struct {
	Id         uuid.UUID
	Request_id uuid.UUID
	Token      string
}

type FormService interface {
	SaveForm(*Form) error
    GetAwaitingForms() ([]RequestData, error)
    FindFormById(req uuid.UUID) (*Request, error)
    LoadFormPhoto(photoName string) (string, error)
    ChangeFormStatus(uuid.UUID, string) error
}

type FormRepository interface {
	InsertForm(*Request, *RequestData) error
    SelectRequests() ([]Request, error)
    SelectRequestById(req uuid.UUID) (*Request, error)
    SelectAwaitingRequests() ([]Request, error)
    SelectRequestDataByRequest(req uuid.UUID) (*RequestData, error)
    UpdateRequestStatus(req uuid.UUID, status, token string) error
}
