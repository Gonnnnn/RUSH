package server

import (
	"rush/attendance"
	"rush/session"
	"rush/user"
)

func fromUser(user *user.User) *User {
	return &User{
		Id:           user.Id,
		Name:         user.Name,
		Generation:   user.Generation,
		IsActive:     user.IsActive,
		Email:        user.Email,
		ExternalName: user.ExternalName,
	}
}

func fromSessionToSessionForAdmin(sessionData session.Session) SessionForAdmin {
	return SessionForAdmin{
		Id:               sessionData.Id,
		Name:             sessionData.Name,
		Description:      sessionData.Description,
		CreatedBy:        sessionData.CreatedBy,
		GoogleFormUri:    sessionData.GoogleFormUri,
		GoogleFormId:     sessionData.GoogleFormId,
		CreatedAt:        sessionData.CreatedAt,
		StartsAt:         sessionData.StartsAt,
		Score:            sessionData.Score,
		AttendanceStatus: sessionData.AttendanceStatus,
		AttendanceAppliedBy: func() SessionAttendanceAppliedBy {
			if sessionData.AttendanceAppliedBy() == session.AttendanceAppliedByUnspecified {
				return SessionAttendanceAppliedByUnspecified
			}
			if sessionData.AttendanceAppliedBy() == session.AttendanceAppliedByManual {
				return SessionAttendanceAppliedByManual
			}
			// TODO(#54): Hard to test as the internal logic of session should be checked. Mock it or find another way.
			if sessionData.AttendanceAppliedBy() == session.AttendanceAppliedByForm {
				return SessionAttendanceAppliedByForm
			}
			return SessionAttendanceAppliedByUnknown
		}(),
	}
}

func fromSessionToSessionForUser(sessionData session.Session) Session {
	return Session{
		Id:          sessionData.Id,
		Name:        sessionData.Name,
		Description: sessionData.Description,
		CreatedBy:   sessionData.CreatedBy,
		CreatedAt:   sessionData.CreatedAt,
		StartsAt:    sessionData.StartsAt,
		Score:       sessionData.Score,
	}
}

func fromAttendance(attendance *attendance.Attendance) *Attendance {
	return &Attendance{
		Id:               attendance.Id,
		SessionId:        attendance.SessionId,
		SessionName:      attendance.SessionName,
		SessionScore:     attendance.SessionScore,
		SessionStartedAt: attendance.SessionStartedAt,
		UserId:           attendance.UserId,
		UserExternalName: attendance.UserExternalName,
		UserGeneration:   attendance.UserGeneration,
		UserJoinedAt:     attendance.UserJoinedAt,
		CreatedAt:        attendance.CreatedAt,
	}
}
