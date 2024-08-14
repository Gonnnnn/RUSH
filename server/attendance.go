package server

import (
	"errors"
	"fmt"
	"rush/golang/array"
	"rush/session"
	"rush/user"
	"slices"
	"sort"
	"time"
)

func (s *Server) CreateAttendanceForm(sessionId string) (string, error) {
	dbSession, err := s.sessionRepo.Get(sessionId)
	if err != nil {
		return "", newNotFoundError(fmt.Errorf("failed to get session: %w", err))
	}
	if dbSession.IsClosed {
		return "", newBadRequestError(errors.New("session is already closed"))
	}

	users, err := s.userRepo.GetAll()
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to get users: %w", err))
	}

	sort.Slice(users, func(i, j int) bool {
		if users[i].Generation != users[j].Generation {
			return users[i].Generation > users[j].Generation
		}
		return users[i].Name < users[j].Name
	})

	if dbSession.GoogleFormUri != "" {
		return "", newBadRequestError(fmt.Errorf("form already exists: URI is %s", dbSession.GoogleFormUri))
	}

	formTitle := fmt.Sprintf("[출석] %s", dbSession.Name)
	startsAt := dbSession.StartsAt.In(s.formTimeLocation)
	expiresAt := startsAt.Add(-time.Second)
	formDescription := fmt.Sprintf(`%s을(를) 위한 출석용 구글폼입니다.
폼 마감 시각은 %s입니다. %s 이후 요청은 무시됩니다.`, dbSession.Name, expiresAt.Format("2006-01-02 15:04:05"), startsAt.Format("2006-01-02 15:04:05"))

	attendanceForm, err := s.attendanceFormHandler.GenerateForm(formTitle, formDescription, users)
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to generate form: %w", err))
	}

	_, err = s.sessionRepo.Update(sessionId, &session.UpdateForm{
		GoogleFormId: &attendanceForm.Id, GoogleFormUri: &attendanceForm.Uri, ReturnUpdatedSession: false})
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to update session: %w", err))
	}

	return attendanceForm.Uri, nil
}

func (s *Server) GetAttendanceByUserId(userId string) ([]Attendance, error) {
	attendances, err := s.attendanceRepo.FindByUserId(userId)
	if err != nil {
		return nil, newInternalServerError(fmt.Errorf("failed to find attendance by user ID: %w", err))
	}
	converted := []Attendance{}
	for _, attendance := range attendances {
		converted = append(converted, *fromAttendance(&attendance))
	}
	return converted, nil
}

type HalfYearAttendace struct {
	// All the sessions that are held in the half year so far.
	Sessions []sessionForAttendance
	// All the users who joined the sessions in the half year so far.
	Users []struct {
		Id         string
		Name       string
		Generation float64
	}
	// All the attendances in the half year so far.
	Attendances []Attendance
}

type sessionForAttendance struct {
	Id        string
	Name      string
	StartedAt time.Time
}

func (s *Server) GetHalfYearAttendance() (HalfYearAttendace, error) {
	// TODO(#113): Save the data of each generation, startsAt, finishesAt, name, etc.
	// And then replace it to a method to get attendance within certain period.
	users, err := s.userRepo.GetAll()
	if err != nil {
		return HalfYearAttendace{}, newInternalServerError(fmt.Errorf("failed to get users: %w", err))
	}
	activeUsers := array.Filter(users, func(user user.User) bool { return user.IsActive })
	slices.SortStableFunc(activeUsers, func(user1, user2 user.User) int {
		if user1.Generation > user2.Generation {
			return 1
		}
		if user1.Generation < user2.Generation {
			return -1
		}
		if user1.Name > user2.Name {
			return 1
		}
		if user1.Name < user2.Name {
			return -1
		}

		return 0
	})

	attendances, err := s.attendanceRepo.GetAll()
	if err != nil {
		return HalfYearAttendace{}, newInternalServerError(fmt.Errorf("failed to find half year attendance: %w", err))
	}
	convertedAttendances := []Attendance{}
	for _, attendance := range attendances {
		convertedAttendances = append(convertedAttendances, *fromAttendance(&attendance))
	}

	uniqueSessionMap := map[string]sessionForAttendance{}
	for _, attendance := range convertedAttendances {
		uniqueSessionMap[attendance.SessionName] = sessionForAttendance{
			Id:        attendance.SessionId,
			Name:      attendance.SessionName,
			StartedAt: attendance.SessionStartedAt,
		}
	}
	uniqueSessions := []sessionForAttendance{}
	for name := range uniqueSessionMap {
		uniqueSessions = append(uniqueSessions, uniqueSessionMap[name])
	}
	slices.SortStableFunc(uniqueSessions, func(session1, session2 sessionForAttendance) int {
		if session1.StartedAt.After(session2.StartedAt) {
			return 1
		}
		return -1
	})

	halfYearAttendace := HalfYearAttendace{
		Sessions: uniqueSessions,
		Users: array.Map(activeUsers, func(user user.User) struct {
			Id         string
			Name       string
			Generation float64
		} {
			return struct {
				Id         string
				Name       string
				Generation float64
			}{Id: user.Id, Name: user.Name, Generation: user.Generation}
		}),
		Attendances: convertedAttendances,
	}
	return halfYearAttendace, nil
}
