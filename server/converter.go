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

func fromSession(session *session.Session) *Session {
	return &Session{
		Id:            session.Id,
		Name:          session.Name,
		Description:   session.Description,
		CreatedBy:     session.CreatedBy,
		GoogleFormUri: session.GoogleFormUri,
		GoogleFormId:  session.GoogleFormId,
		CreatedAt:     session.CreatedAt,
		StartsAt:      session.StartsAt,
		Score:         session.Score,
		IsClosed:      session.IsClosed,
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
