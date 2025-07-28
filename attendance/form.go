package attendance

import (
	"fmt"
	"math"
	"rush/golang/array"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/forms/v1"
)

// The option that will be shown in the form.
// Users select one of them to indicate that they joined the session.
type UserOption struct {
	// The generation of the user. E.g. 9.5
	Generation float64
	// The external name of the user. E.g. "김건3"
	// There could be multiple users with the same name.
	// External name is used to distinguish them by adding a number as suffix.
	ExternalName string
}

type Form struct {
	// The Google form ID.
	Id string
	// The Google form URI. It's the public URL that users access to fill out the form.
	Uri string
}

type FormSubmission struct {
	// The external name of the user that is exposed to the form.
	// Use it to match the submission with the user. E.g., "김건3"
	UserExternalName string
	// The time when the form was submitted.
	SubmissionTime time.Time
}

type formHandler struct {
	// The Google Forms service to make forms.
	googleFormService *forms.Service
	// The Google Drive service to manage permissions to the form.
	googleDriveService *drive.Service
	// The form option parser to get the user external name from the form.
	formOptionParser *formOptionParser
	// The delimiter to separate the generation and the external name in the form option.
	delimiter string
	// The emails of the admins who will be able to manage the form.
	adminEmails []string
}

func NewFormHandler(googleFormService *forms.Service, googleDriveService *drive.Service) *formHandler {
	delimiter := " - "
	// 김건, 양현우
	adminEmails := []string{"geonkim23@gmail.com", "alsrudrkd13@gmail.com"}
	return &formHandler{
		googleFormService:  googleFormService,
		googleDriveService: googleDriveService,
		formOptionParser:   newFormOptionParser(delimiter),
		delimiter:          delimiter,
		adminEmails:        adminEmails,
	}
}

// Generates a new Google form with the given arguments. It returns a Form object after uploading it to the Google Drive.
func (f *formHandler) GenerateForm(title string, description string, userOptions []UserOption) (Form, error) {
	newForm := &forms.Form{Info: &forms.Info{Title: title, DocumentTitle: title}}
	form, err := f.googleFormService.Forms.Create(newForm).Do()
	if err != nil {
		return Form{}, fmt.Errorf("failed to create form: %w", err)
	}

	question := &forms.Question{
		Required: true,
		ChoiceQuestion: &forms.ChoiceQuestion{
			Type: "DROP_DOWN",
			Options: array.Map(userOptions, func(userOption UserOption) *forms.Option {
				return &forms.Option{Value: newFormOption(userOption.Generation, userOption.ExternalName, f.delimiter).string()}
			}),
		},
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
						Description: "기수와 이름을 선택해주세요.\n선택지는 1. 기수 2. 이름 순으로 정렬돼있습니다.\nformat: `기수 - 이름`",
						QuestionItem: &forms.QuestionItem{
							Question: question,
						},
					},
					Location: &forms.Location{
						Index: 0,
						// 0 is the default value in golang that it is ignored (omitempty).
						// Thus, `Index` should be specified as `ForceSendFields`.
						ForceSendFields: []string{"Index"},
					},
				},
			},
		},
	}

	_, err = f.googleFormService.Forms.BatchUpdate(form.FormId, updateRequest).Do()
	if err != nil {
		return Form{}, fmt.Errorf("failed to add the question to the form: %w", err)
	}

	for _, adminEmail := range f.adminEmails {
		permission := &drive.Permission{
			Type:         "user",
			Role:         "writer",
			EmailAddress: adminEmail,
		}

		_, err = f.googleDriveService.Permissions.Create(form.FormId, permission).Do()
		if err != nil {
			return Form{}, fmt.Errorf("failed to create writer permission for admin: %w", err)
		}
	}

	return Form{Id: form.FormId, Uri: form.ResponderUri}, nil
}

// Fetches all the submissions of the form.
func (f *formHandler) GetSubmissions(formId string) ([]FormSubmission, error) {
	responses, err := f.googleFormService.Forms.Responses.List(formId).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch form responses: %w", err)
	}

	var userExternalNameSubmissionTimeMap = make(map[string]time.Time)
	for _, response := range responses.Responses {
		submissionTime, err := time.Parse(time.RFC3339, response.LastSubmittedTime)
		if err != nil {
			return nil, fmt.Errorf("failed to parse submission time: %w", err)
		}

		externalName, err := f.getUserExternalName(response)
		if err != nil {
			return nil, fmt.Errorf("failed to get user external name from response: %w", err)
		}

		timeFoundBefore, ok := userExternalNameSubmissionTimeMap[externalName]
		if ok && submissionTime.After(timeFoundBefore) {
			// The user might have submitted the form multiple times.
			// Only keep the first submission.
			continue
		}

		userExternalNameSubmissionTimeMap[externalName] = submissionTime
	}

	var submissions []FormSubmission
	for externalName, submissionTime := range userExternalNameSubmissionTimeMap {
		submissions = append(submissions, FormSubmission{UserExternalName: externalName, SubmissionTime: submissionTime})
	}

	return submissions, nil
}

func (f *formHandler) getUserExternalName(response *forms.FormResponse) (string, error) {
	for _, answer := range response.Answers {
		if answer.TextAnswers == nil || len(answer.TextAnswers.Answers) == 0 {
			return "", fmt.Errorf("invalid answer was read from the form which seems like the form API problem: %v", answer)
		}

		selectedOption := answer.TextAnswers.Answers[0].Value
		option, err := f.formOptionParser.parse(selectedOption)
		if err != nil {
			return "", fmt.Errorf("failed to parse user external name: %w", err)
		}

		return option.ExternalName, nil
	}

	return "", fmt.Errorf("no answer was found in the response")
}

type formOption struct {
	Generation   float64
	ExternalName string
	delimiter    string
}

func newFormOption(generation float64, externalName string, delimiter string) formOption {
	return formOption{Generation: generation, ExternalName: externalName, delimiter: delimiter}
}

func (f formOption) string() string {
	if math.Trunc(f.Generation) == f.Generation {
		return fmt.Sprintf("%s%s%s", strconv.Itoa(int(f.Generation)), f.delimiter, f.ExternalName)
	}
	return fmt.Sprintf("%s%s%s", fmt.Sprintf("%.1f", f.Generation), f.delimiter, f.ExternalName)
}

type formOptionParser struct {
	delimiter string
}

func newFormOptionParser(delimiter string) *formOptionParser {
	return &formOptionParser{delimiter: delimiter}
}

func (f *formOptionParser) parse(option string) (formOption, error) {
	parts := strings.Split(option, f.delimiter)
	if len(parts) != 2 {
		return formOption{}, fmt.Errorf("invalid option format: %s", option)
	}

	generation, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return formOption{}, fmt.Errorf("failed to parse generation: %w", err)
	}

	return formOption{Generation: generation, ExternalName: parts[1], delimiter: f.delimiter}, nil
}
