package job

import (
	"errors"
	"rush/session"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestCloseExpiredSessions(t *testing.T) {
	t.Run("Fails if it fails to get open sessions", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionGetter := NewMocksessionGetter(controller)
		sessionCloser := NewMocksessionCloser(controller)
		mockLogger := NewMockLogger(controller)
		executor := NewExecutor(sessionGetter, sessionCloser, mockLogger, clock.NewMock())

		sessionGetter.EXPECT().GetOpenSessionsWithForm().Return([]session.Session{}, assert.AnError)
		mockLogger.EXPECT().Errorw("Failed to get open sessions with form", "error", assert.AnError.Error())
		executor.CloseExpiredSessions()
	})

	t.Run("Fails if it fails to close ession", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionGetter := NewMocksessionGetter(controller)
		sessionCloser := NewMocksessionCloser(controller)
		mockLogger := NewMockLogger(controller)
		clock := clock.NewMock()
		executor := NewExecutor(sessionGetter, sessionCloser, mockLogger, clock)

		clock.Set(time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC))
		sessionGetter.EXPECT().GetOpenSessionsWithForm().Return([]session.Session{
			{Id: "sessionId1", StartsAt: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)},
			{Id: "sessionId2", StartsAt: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)},
			{Id: "sessionId3", StartsAt: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)},
		}, nil)
		sessionCloser.EXPECT().CloseSession("sessionId1").Return(errors.New("error1"))
		sessionCloser.EXPECT().CloseSession("sessionId2").Return(nil)
		sessionCloser.EXPECT().CloseSession("sessionId3").Return(errors.New("error2"))
		mockLogger.EXPECT().Infow("Closed sessions", "session_ids", "sessionId2")
		mockLogger.EXPECT().Errorw("Failed to close sessions", "session_ids", "sessionId1, sessionId3", "errors", "error1, error2")
		executor.CloseExpiredSessions()
	})

	t.Run("Successfully close open and also expired sessions", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionGetter := NewMocksessionGetter(controller)
		sessionCloser := NewMocksessionCloser(controller)
		mockLogger := NewMockLogger(controller)
		clock := clock.NewMock()
		executor := NewExecutor(sessionGetter, sessionCloser, mockLogger, clock)

		clock.Set(time.Date(2024, 1, 2, 12, 30, 0, 0, time.UTC))
		sessionGetter.EXPECT().GetOpenSessionsWithForm().Return([]session.Session{
			{Id: "sessionId1", StartsAt: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)},
			{Id: "sessionId2", StartsAt: time.Date(2024, 1, 2, 12, 30, 0, 0, time.UTC)},
			{Id: "sessionId3", StartsAt: time.Date(2024, 1, 3, 12, 30, 0, 0, time.UTC)},
		}, nil)
		sessionCloser.EXPECT().CloseSession("sessionId1").Return(nil)
		sessionCloser.EXPECT().CloseSession("sessionId2").Return(nil)
		// sessionId3 is not expired yet.
		mockLogger.EXPECT().Infow("Closed sessions", "session_ids", "sessionId1, sessionId2")
		executor.CloseExpiredSessions()
	})
}
