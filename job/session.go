package job

import (
	"rush/golang/array"
	"rush/session"
	"strings"

	"github.com/benbjohnson/clock"
)

type executor struct {
	sessionGetter sessionGetter
	sessionCloser sessionCloser
	logger        Logger
	clock         clock.Clock
}

func NewExecutor(sessionGetter sessionGetter, sessionCloser sessionCloser, logger Logger, clock clock.Clock) *executor {
	return &executor{
		sessionGetter: sessionGetter,
		sessionCloser: sessionCloser,
		logger:        logger,
		clock:         clock,
	}
}

// Closes the open sessions that are past the start time.
func (e *executor) CloseExpiredSessions() {
	openSessions, err := e.sessionGetter.GetOpenSessionsWithForm()
	if err != nil {
		e.logger.Errorw("Failed to get open sessions with form", "error", err.Error())
		return
	}

	now := e.clock.Now()
	sessionsToClose := array.Filter(openSessions, func(session session.Session) bool {
		return now.After(session.StartsAt) || now.Equal(session.StartsAt)
	})

	failedSessionIds := []string{}
	succeededSessionIds := []string{}
	closeErr := []error{}
	for _, session := range sessionsToClose {
		if err := e.sessionCloser.CloseSession(session.Id); err != nil {
			failedSessionIds = append(failedSessionIds, session.Id)
			closeErr = append(closeErr, err)
			continue
		}
		succeededSessionIds = append(succeededSessionIds, session.Id)
	}

	e.logger.Infow("Closed sessions", "session_ids", strings.Join(succeededSessionIds, ", "))
	if len(failedSessionIds) > 0 {
		e.logger.Errorw("Failed to close sessions", "session_ids", strings.Join(failedSessionIds, ", "),
			"errors", strings.Join(array.Map(closeErr, func(err error) string { return err.Error() }), ", "))
		return
	}
}

//go:generate mockgen -source=session.go -destination=session_mock.go -package=job
type sessionCloser interface {
	CloseSession(sessionId string) error
}

type sessionGetter interface {
	GetOpenSessionsWithForm() ([]session.Session, error)
}

type Logger interface {
	// Logs the given info with the info level.
	// Info level indicates any information that should be logged.
	Infow(msg string, keysAndValues ...any)
	// Logs the given info with the error level.
	// Error level indicates any issue that should be resolved as soon as possible.
	Errorw(msg string, keysAndValues ...any)
}
