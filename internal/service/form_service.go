package service

import (
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
