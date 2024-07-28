package server

import (
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
