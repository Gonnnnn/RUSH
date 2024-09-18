package job

import (
	"fmt"
	"rush/golang/array"
	"rush/session"

	"github.com/benbjohnson/clock"
)

type executor struct {
	sessionGetter sessionGetter
	sessionCloser sessionCloser
	clock         clock.Clock
}

type CloseExpiredSessionsResult struct {
	// The IDs of the sessions that succeeded to close.
	SucceededSessionIds []string
	// The IDs of the sessions that failed to close.
	FailedSessionIds []string
	// The errors that occurred while closing the sessions.
	// The order of the errors corresponds to the order of the sessions.
	Errors []error
}

func NewExecutor(sessionGetter sessionGetter, sessionCloser sessionCloser, clock clock.Clock) *executor {
	return &executor{
		sessionGetter: sessionGetter,
		sessionCloser: sessionCloser,
		clock:         clock,
	}
}

// Closes the open sessions that are past the start time.
func (e *executor) CloseExpiredSessions() (CloseExpiredSessionsResult, error) {
	openSessions, err := e.sessionGetter.GetOpenSessions()
	if err != nil {
		return CloseExpiredSessionsResult{}, fmt.Errorf("failed to get open sessions: %w", err)
	}

	openSessions = array.Filter(openSessions, func(session session.Session) bool {
		return session.StartsAt.Before(e.clock.Now()) || session.StartsAt.Equal(e.clock.Now())
	})

	failedSessionIds := []string{}
	succeededSessionIds := []string{}
	errors := []error{}
	for _, session := range openSessions {
		if err := e.sessionCloser.CloseSession(session.Id); err != nil {
			failedSessionIds = append(failedSessionIds, session.Id)
			errors = append(errors, err)
			continue
		}
		succeededSessionIds = append(succeededSessionIds, session.Id)
	}

	if len(failedSessionIds) == 0 {
		return CloseExpiredSessionsResult{
			SucceededSessionIds: succeededSessionIds,
		}, nil
	}

	return CloseExpiredSessionsResult{
		SucceededSessionIds: succeededSessionIds,
		FailedSessionIds:    failedSessionIds,
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
