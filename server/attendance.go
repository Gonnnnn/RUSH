package server

import (
	"errors"
	"fmt"
	"rush/attendance"
	attendanceAggregation "rush/attendance/aggregation"
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

	// TODO(#134): Fetch active users only.
	users, err := s.userRepo.GetAll()
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to get users: %w", err))
	}

	activeUsers := array.Filter(users, func(user user.User) bool { return user.IsActive })

	sort.Slice(activeUsers, func(i, j int) bool {
		if activeUsers[i].Generation != activeUsers[j].Generation {
			return activeUsers[i].Generation > activeUsers[j].Generation
		}
		return activeUsers[i].Name < activeUsers[j].Name
	})

	if dbSession.GoogleFormUri != "" {
		return "", newBadRequestError(fmt.Errorf("form already exists: URI is %s", dbSession.GoogleFormUri))
	}

	formTitle := fmt.Sprintf("[출석] %s", dbSession.Name)
	startsAt := dbSession.StartsAt.In(s.formTimeLocation)
	expiresAt := startsAt.Add(-time.Second)
	formDescription := fmt.Sprintf(`%s을(를) 위한 출석용 구글폼입니다.
폼 마감 시간은 %s입니다. %s 이후 요청은 무시됩니다.`, dbSession.Name, expiresAt.Format("2006-01-02 15:04:05"), startsAt.Format("2006-01-02 15:04:05"))

	attendanceForm, err := s.attendanceFormHandler.GenerateForm(formTitle, formDescription, activeUsers)
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to generate form: %w", err))
	}

	_, err = s.openSessionRepo.UpdateOpenSession(sessionId, session.OpenSessionUpdateForm{
		GoogleFormId:  &attendanceForm.Id,
		GoogleFormUri: &attendanceForm.Uri,
	})
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
	Sessions []sessionForAttendance `json:"sessions"`
	// All the users who joined the sessions in the half year so far.
	Users []userForAttendance `json:"users"`
	// All the attendances in the half year so far.
	Attendances []Attendance `json:"attendances"`
}

type userForAttendance struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Generation float64 `json:"generation"`
}

type sessionForAttendance struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	StartedAt time.Time `json:"started_at"`
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
		Users: array.Map(activeUsers, func(user user.User) userForAttendance {
			return userForAttendance{
				Id:         user.Id,
				Name:       user.Name,
				Generation: user.Generation,
			}
		}),
		Attendances: convertedAttendances,
	}
	return halfYearAttendace, nil
}

func (s *Server) AggregateAttendance() error {
	attendances, err := s.attendanceRepo.GetAll()
	if err != nil {
		return newInternalServerError(fmt.Errorf("failed to get all attendances: %w", err))
	}

	userScoreMap := map[string]int{}
	for _, attendance := range attendances {
		userScoreMap[attendance.UserId] += attendance.SessionScore
	}

	userInfosForAggregation := array.Map(attendances, func(attendance attendance.Attendance) attendanceAggregation.UserInfo {
		userScore, ok := userScoreMap[attendance.UserId]
		if !ok {
			userScore = 0
		}
		return attendanceAggregation.UserInfo{
			UserId:     attendance.UserId,
			UserName:   attendance.UserExternalName,
			Generation: attendance.UserGeneration,
			Score:      userScore,
		}
	})

	if _, err := s.attendanceAggregationRepo.AddAggregation(
		array.Map(attendances, func(attendance attendance.Attendance) string { return attendance.SessionId }),
		userInfosForAggregation,
	); err != nil {
		return newInternalServerError(fmt.Errorf("failed to add aggregation: %w", err))
	}

	// TODO(#165) Apply user rankings so that their current ranking can be fetched from the user document.
	return nil
}
