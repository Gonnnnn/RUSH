package job

import (
	"fmt"
	"rush/golang/array"
	"rush/session"
)

type executor struct {
	sessionGetter sessionGetter
	sessionCloser sessionCloser
}

type CloseSessionsResult struct {
	// The IDs of the sessions that succeeded to close.
	SessionIdsSucceeded []string
	// The IDs of the sessions that failed to close.
	FailedSessionIds []string
	// The errors that occurred while closing the sessions.
	// The order of the errors corresponds to the order of the sessions.
	Errors []error
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

	sessionsSucceeded := array.Filter(openSessions, func(session session.Session) bool {
		return !array.Contains(failedSessions, session.Id)
	})
	sessionIdsSucceeded := array.Map(sessionsSucceeded, func(session session.Session) string {
		return session.Id
	})
	if len(failedSessions) == 0 {
		return CloseSessionsResult{
			SessionIdsSucceeded: sessionIdsSucceeded,
		}, nil
	}

	return CloseSessionsResult{
		SessionIdsSucceeded: sessionIdsSucceeded,
		FailedSessionIds:    failedSessions,
		Errors:              errors,
	}, fmt.Errorf("failed to close some sessions")
}

//go:generate mockgen -source=session.go -destination=session_mock.go -package=job
type sessionCloser interface {
	CloseSession(sessionId string) error
}

type sessionGetter interface {
	GetOpenSessions() ([]session.Session, error)
}
