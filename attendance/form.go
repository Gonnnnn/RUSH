package attendance

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"google.golang.org/api/forms/v1"

	"rush/user"
)

type formHandler struct {
	googleFormService   *forms.Service
	userOptionDelimiter string
}

func NewFormHandler(googleFormService *forms.Service) *formHandler {
	return &formHandler{googleFormService: googleFormService, userOptionDelimiter: "-"}
}

func (f *formHandler) GenerateForm(title string, description string, users []user.User) (string, error) {
	newForm := &forms.Form{Info: &forms.Info{Title: title}}

	form, err := f.googleFormService.Forms.Create(newForm).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create form: %w", err)
	}

	question := &forms.Question{
		Required: true,
		ChoiceQuestion: &forms.ChoiceQuestion{
			Type:    "DROP_DOWN",
			Options: make([]*forms.Option, len(users)),
		},
	}

	for index, user := range users {
		question.ChoiceQuestion.Options[index] = &forms.Option{Value: f.attendanceOption(user)}
	}

	updateRequest := &forms.BatchUpdateFormRequest{
		Requests: []*forms.Request{
			{
				UpdateFormInfo: &forms.UpdateFormInfoRequest{
					Info: &forms.Info{
						Description: description,
					},
					UpdateMask: "description",
				},
			},
			{
				CreateItem: &forms.CreateItemRequest{
					Item: &forms.Item{
						Title:       "기수:이름",
						Description: "기수와 이름을 선택해주세요.\n선택지는 1. 기수 2. 이름 순으로 정렬돼있습니다.\nformat: `기수-이름`",
						QuestionItem: &forms.QuestionItem{
							Question: question,
						},
					},
					Location: &forms.Location{
						Index: 0,
						// 0 is the default value that it is ignored (omitempty). `Index` should be specified as `ForceSendFields`.
						ForceSendFields: []string{"Index"},
					},
				},
			},
		},
	}

	_, err = f.googleFormService.Forms.BatchUpdate(form.FormId, updateRequest).Do()
	if err != nil {
		return "", fmt.Errorf("failed to update form: %w", err)
	}

	return form.ResponderUri, nil
}

func (f *formHandler) ReadUsers(formId string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (f *formHandler) attendanceOption(user user.User) string {
	var generationStr string
	if math.Trunc(user.Generation) == user.Generation {
		generationStr = strconv.Itoa(int(user.Generation))
	} else {
		generationStr = fmt.Sprintf("%.1f", user.Generation)
	}

	return fmt.Sprintf("%s%s%s", generationStr, f.userOptionDelimiter, user.Name)
}
