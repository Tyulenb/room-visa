package service

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"io"
	"log"
	"os"
	"room-visa/internal/crypto"
	"room-visa/internal/model"
	"room-visa/internal/storage"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
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
func (fs *FormService) SaveForm(form *model.Form) (string, error) {
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

	return requestUUID.String(), fs.fr.InsertForm(request, data)
}

//Returns all forms that need to be checked
func (fs *FormService) GetAwaitingFormsData() ([]model.RequestData, error) {
    requests, err := fs.fr.SelectAwaitingRequests()
    if err != nil {
        return nil, err
    }

    requestData := make([]model.RequestData, 0)

    for _, r := range requests{
        data, err := fs.fr.SelectRequestDataByRequestId(r.Id)
        if err != nil {
            return nil, err
        }
        requestData = append(requestData, *data) 
    }
    return requestData, nil
}

func (fs *FormService) FindFormDataById(req uuid.UUID) (*model.RequestData, error) {
    return fs.fr.SelectRequestDataByRequestId(req)
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

func (fs *FormService) ChangeFormStatus(id uuid.UUID, status string) error {
    token, err := crypto.GenerateVisaToken(id) 
    if err != nil {
        return err
    }
    return fs.fr.UpdateRequestStatus(id, status, token) 
}

func (fs *FormService) FindFormById(id uuid.UUID) (*model.Request, error) {
    return fs.fr.SelectRequestById(id)
}

func (fs *FormService) ValidateVisaToken(tokenString string) (string, error) { 
    return crypto.GetVisaTokenClaimReq(tokenString)
}

func (fs *FormService) GenerateVisaQR(id uuid.UUID) (string, error) {
    visa, err := fs.fr.SelectVisaByRequestId(id)
    if err != nil {
        return "", err
    }

    domain := os.Getenv("DOMAIN")
    port := os.Getenv("PORT")
    encodeString := domain+port+"/validate?token="+visa.Token
    log.Println(encodeString)
    
    var buf bytes.Buffer 
    qrcode, err := qr.Encode(encodeString, qr.L, qr.Auto)
    if err != nil {
        return "", err
    }
    qrcode, err = barcode.Scale(qrcode, 200, 200)
    if err != nil {
        return "", err
    }
    if err = png.Encode(&buf, qrcode); err != nil {
        return "", err
    }
    enc := base64.StdEncoding.EncodeToString(buf.Bytes())
    return enc, nil
}

