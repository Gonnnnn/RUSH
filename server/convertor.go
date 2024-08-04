package server

import (
	"rush/attendance"
	"rush/session"
	"rush/user"
	"strconv"
)

func fromUser(user *user.User) *User {
	return &User{
		Id:         user.Id,
		Name:       user.Name,
		University: user.University,
		Phone:      user.Phone,
		Generation: user.Generation,
		IsActive:   user.IsActive,
	}
}

func fromSession(session *session.Session) *Session {
	return &Session{
		Id:            session.Id,
		Name:          session.Name,
		Description:   session.Description,
		HostedBy:      strconv.Itoa(session.HostedBy),
		CreatedBy:     strconv.Itoa(session.CreatedBy),
		GoogleFormUri: session.GoogleFormUri,
		JoinningUsers: session.JoinningUsers,
		CreatedAt:     session.CreatedAt,
		StartsAt:      session.StartsAt,
		Score:         session.Score,
		IsClosed:      session.IsClosed,
	}
}

func fromAttendance(attendance *attendance.Attendance) *Attendance {
	return &Attendance{
		Id:          attendance.Id,
		SessionId:   attendance.SessionId,
		SessionName: attendance.SessionName,
		UserId:      attendance.UserId,
		UserName:    attendance.UserName,
		JoinedAt:    attendance.JoinedAt,
		CreatedAt:   attendance.CreatedAt,
	}
}
