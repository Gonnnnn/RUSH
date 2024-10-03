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

func fromSession(dbSession session.Session) Session {
	return Session{
		Id:            dbSession.Id,
		Name:          dbSession.Name,
		Description:   dbSession.Description,
		CreatedBy:     dbSession.CreatedBy,
		GoogleFormUri: dbSession.GoogleFormUri,
		GoogleFormId:  dbSession.GoogleFormId,
		CreatedAt:     dbSession.CreatedAt,
		StartsAt:      dbSession.StartsAt,
		Score:         dbSession.Score,
		IsClosed:      dbSession.AttendanceStatus == session.AttendanceStatusApplied,
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
