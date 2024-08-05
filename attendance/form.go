package attendance

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/forms/v1"

	"rush/user"
)

type Form struct {
	// The Google form ID.
	Id string
	// The Google form URI. It's what users access to fill out the form.
	Uri string
}

type FormSubmission struct {
	// The external ID of the user that is exposed to the form.
	// Use it to match the submission with the user. E.g., "abc123"
	UserExternalId string
	// The time when the form was submitted.
	SubmissionTime time.Time
}

type formHandler struct {
	// The Google Forms service to make forms.
	googleFormService *forms.Service
	// The Google Drive service to manage permissions to the form.
	googleDriveService *drive.Service
	// The delimiter to separate the user's generation, name, and external ID in the form option.
	userOptionDelimiter string
}

// 김건, 양현우
var adminEmails = []string{"geonkim23@gmail.com", "hyeonyi30754@gmail.com"}

func NewFormHandler(googleFormService *forms.Service, googleDriveService *drive.Service) *formHandler {
	return &formHandler{googleFormService: googleFormService, googleDriveService: googleDriveService, userOptionDelimiter: " - "}
}

func (f *formHandler) GenerateForm(title string, description string, users []user.User) (Form, error) {
	newForm := &forms.Form{Info: &forms.Info{Title: title, DocumentTitle: title}}

	form, err := f.googleFormService.Forms.Create(newForm).Do()
	if err != nil {
		return Form{}, fmt.Errorf("failed to create form: %w", err)
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
						Description: "기수와 이름을 선택해주세요.\n선택지는 1. 기수 2. 이름 순으로 정렬돼있습니다.\nformat: `기수 - 이름 - ID`",
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
		return Form{}, fmt.Errorf("failed to update form: %w", err)
	}

	for _, adminEmail := range adminEmails {
		permission := &drive.Permission{
			Type:         "user",
			Role:         "writer",
			EmailAddress: adminEmail,
		}

		_, err = f.googleDriveService.Permissions.Create(form.FormId, permission).Do()
		if err != nil {
			return Form{}, fmt.Errorf("failed to create permission: %w", err)
		}
	}

	return Form{Id: form.FormId, Uri: form.ResponderUri}, nil
}

func (f *formHandler) GetSubmissions(formId string) ([]FormSubmission, error) {
	responses, err := f.googleFormService.Forms.Responses.List(formId).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch form responses: %w", err)
	}

	var userExternalIdSubmissionTimeMap = make(map[string]time.Time)
	for _, response := range responses.Responses {
		submissionTime, err := time.Parse(time.RFC3339, response.LastSubmittedTime)
		if err != nil {
			return nil, fmt.Errorf("failed to parse submission time: %w", err)
		}

		externalId, err := f.getUserExternalIdFromResponse(response)
		if err != nil {
			return nil, fmt.Errorf("failed to get user external ID from response: %w", err)
		}

		timeFoundBefore, ok := userExternalIdSubmissionTimeMap[externalId]
		if ok && submissionTime.After(timeFoundBefore) {
			// The user might have submitted the form multiple times. We only keep the first submission.
			continue
		}

		userExternalIdSubmissionTimeMap[externalId] = submissionTime
	}

	var submissions []FormSubmission
	for externalId, submissionTime := range userExternalIdSubmissionTimeMap {
		submissions = append(submissions, FormSubmission{UserExternalId: externalId, SubmissionTime: submissionTime})
	}

	return submissions, nil
}

// TODO(#42): Save the user ID somewhere else. Not safe to include it here and couple it with the parsing logic.
func (f *formHandler) attendanceOption(user user.User) string {
	var generationStr string
	if math.Trunc(user.Generation) == user.Generation {
		generationStr = strconv.Itoa(int(user.Generation))
	} else {
		generationStr = fmt.Sprintf("%.1f", user.Generation)
	}

	return fmt.Sprintf("%s%s%s%s%s", generationStr, f.userOptionDelimiter, user.Name, f.userOptionDelimiter, user.ExternalId)
}

func (f *formHandler) parseUserExternalId(option string) (string, error) {
	parts := strings.Split(option, f.userOptionDelimiter)
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid option format: %s", option)
	}

	return parts[2], nil
}

func (f *formHandler) getUserExternalIdFromResponse(response *forms.FormResponse) (string, error) {
	for _, answer := range response.Answers {
		if answer.TextAnswers == nil || len(answer.TextAnswers.Answers) == 0 {
			return "", fmt.Errorf("invalid answer was read from the form which seems like the form API problem: %v", answer)
		}

		selectedOption := answer.TextAnswers.Answers[0].Value
		externalId, err := f.parseUserExternalId(selectedOption)
		if err != nil {
			return "", fmt.Errorf("failed to parse user external ID: %w", err)
		}

		return externalId, nil
	}

	return "", fmt.Errorf("no answer was found in the response")
}
