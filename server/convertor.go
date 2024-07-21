package server

import (
	"rush/session"
	"rush/user"
	"strconv"
	"strings"
)

func fromUser(user *user.User) *User {
	return &User{
		Id:         strconv.Itoa(user.Id),
		Name:       user.Name,
		University: user.University,
		Phone:      user.Phone,
		Generation: user.Generation,
		IsActive:   user.IsActive,
	}
}

func fromSession(session *session.Session) *Session {
	return &Session{
		Id:            strconv.Itoa(session.Id),
		Name:          session.Name,
		Description:   session.Description,
		HostedBy:      strconv.Itoa(session.HostedBy),
		CreatedBy:     strconv.Itoa(session.CreatedBy),
		JoinningUsers: strings.Split(string(session.JoinningUsers), ","),
		CreatedAt:     session.CreatedAt,
		StartsAt:      session.StartsAt,
		Score:         session.Score,
		IsClosed:      session.IsClosed,
	}
}
