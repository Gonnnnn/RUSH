package session

import (
	"errors"

	"google.golang.org/api/forms/v1"

	"rush/user"
)

type formHandler struct {
	googleFormService *forms.Service
}

func NewFormHandler(googleFormService *forms.Service) *formHandler {
	return &formHandler{googleFormService: googleFormService}
}

func (f *formHandler) GenerateForm(title string, description string, users []user.User) (string, error) {
	newForm := &forms.Form{
		Info: &forms.Info{
			Title: title,
		},
	}

	form, err := f.googleFormService.Forms.Create(newForm).Do()
	if err != nil {
		return "", err
	}

	return form.ResponderUri, nil
}

func (f *formHandler) ReadUsers(formId string) ([]string, error) {
	return nil, errors.New("not implemented")
}
