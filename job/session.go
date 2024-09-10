package job

import (
	"fmt"
	"rush/session"
)

type executor struct {
	sessionGetter sessionGetter
	sessionCloser sessionCloser
}

type CloseSessionsResult struct {
	// The IDs of the sessions that failed to close.
	failedSessionIds []string
	// The errors that occurred while closing the sessions.
	// The order of the errors corresponds to the order of the sessions.
	errors []error
}

func NewExecutor(sessionGetter sessionGetter, sessionCloser sessionCloser) *executor {
	return &executor{
		sessionGetter: sessionGetter,
		sessionCloser: sessionCloser,
	}
}

func (e *executor) CloseSessions() (CloseSessionsResult, error) {
	openSessions, err := e.sessionGetter.GetOpenSessions()
	if err != nil {
		return CloseSessionsResult{}, fmt.Errorf("failed to get open sessions: %w", err)
	}

	failedSessions := []string{}
	errors := []error{}
	for _, session := range openSessions {
		if err := e.sessionCloser.CloseSession(session.Id); err != nil {
			failedSessions = append(failedSessions, session.Id)
			errors = append(errors, err)
		}
	}

	if len(failedSessions) == 0 {
		return CloseSessionsResult{}, nil
	}

	return CloseSessionsResult{
		failedSessionIds: failedSessions,
		errors:           errors,
	}, fmt.Errorf("failed to close some sessions")
}

//go:generate mockgen -source=session.go -destination=session_mock.go -package=job
type sessionCloser interface {
	CloseSession(sessionId string) error
}

type sessionGetter interface {
	GetOpenSessions() ([]*session.Session, error)
}
