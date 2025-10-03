package service

import (
	"encoding/base64"
	"io"
	"room-visa/internal/model"
	"room-visa/internal/storage"
	"time"

	"github.com/google/uuid"
)

type FormService struct {
	st storage.Storage
	fr model.FormRepository
}

func NewFormService(fr model.FormRepository, st storage.Storage) *FormService {
	return &FormService{
		fr: fr,
		st: st,
	}
}

//Saves forms to database and profile photo to storage
func (fs *FormService) SaveForm(form *model.Form) error {
	photoName := uuid.NewString()
	fs.st.Save(photoName, form.Photo)
	defer form.Photo.Close()

	requestUUID := uuid.New()

	timeNow := time.Now()
	request := &model.Request{
		Id:         requestUUID,
		Status:     "Awaiting",
		Created_at: timeNow,
		Updated_at: timeNow,
	}

	data := &model.RequestData{
		Request_id:  requestUUID,
		Name:        form.Name,
		Surname:     form.Surname,
		Sex:         form.Sex,
		Ethnicity:   form.Ethnicity,
		Citizenship: form.Citizenship,
		Purpose:     form.Purpose,
		Photo:       photoName,
	}

	return fs.fr.InsertForm(request, data)
}

//Returns all forms that need to be checked
func (fs *FormService) GetForms() ([]model.RequestData, error) {
    requests, err := fs.fr.SelectAwaitingRequests()
    if err != nil {
        return nil, err
    }

    requestData := make([]model.RequestData, 0)

    for _, r := range requests{
        data, err := fs.fr.SelectRequestDataByRequest(r.Id)
        if err != nil {
            return nil, err
        }
        requestData = append(requestData, *data) 
    }
    return requestData, nil
}

//Loads photo from storage
//Encodes it in base64 string format
func (fs *FormService) LoadFormPhoto(photoName string) (string, error) {
    photo, err := fs.st.GetByName(photoName)
    if err != nil {
        return "", err
    }
    defer photo.Close()

    dataBytes, err := io.ReadAll(photo)
    if err != nil {
        return "", err
    }
    enc := base64.StdEncoding.EncodeToString(dataBytes)
    return enc, nil
}
