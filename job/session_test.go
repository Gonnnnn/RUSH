package job

import (
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
		executor := NewExecutor(sessionGetter, sessionCloser, clock.NewMock())

		sessionGetter.EXPECT().GetOpenSessions().Return([]session.Session{}, assert.AnError)
		result, err := executor.CloseExpiredSessions()

		assert.Equal(t, CloseExpiredSessionsResult{}, result)
		assert.EqualError(t, err, "failed to get open sessions: "+assert.AnError.Error())
	})

	t.Run("Fails if it fails to close ession", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionGetter := NewMocksessionGetter(controller)
		sessionCloser := NewMocksessionCloser(controller)
		clock := clock.NewMock()
		executor := NewExecutor(sessionGetter, sessionCloser, clock)

		clock.Set(time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC))
		sessionGetter.EXPECT().GetOpenSessions().Return([]session.Session{
			{Id: "sessionId1", StartsAt: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)},
			{Id: "sessionId2", StartsAt: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)},
			{Id: "sessionId3", StartsAt: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)},
		}, nil)
		sessionCloser.EXPECT().CloseSession("sessionId1").Return(assert.AnError)
		sessionCloser.EXPECT().CloseSession("sessionId2").Return(nil)
		sessionCloser.EXPECT().CloseSession("sessionId3").Return(assert.AnError)
		result, err := executor.CloseExpiredSessions()

		assert.Equal(t, CloseExpiredSessionsResult{
			SucceededSessionIds: []string{"sessionId2"},
			FailedSessionIds:    []string{"sessionId1", "sessionId3"},
			Errors:              []error{assert.AnError, assert.AnError},
		}, result)
		assert.EqualError(t, err, "failed to close some sessions")
	})

	t.Run("Successfully close open and also expired sessions", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionGetter := NewMocksessionGetter(controller)
		sessionCloser := NewMocksessionCloser(controller)
		clock := clock.NewMock()
		executor := NewExecutor(sessionGetter, sessionCloser, clock)

		clock.Set(time.Date(2024, 1, 2, 12, 30, 0, 0, time.UTC))
		sessionGetter.EXPECT().GetOpenSessions().Return([]session.Session{
			{Id: "sessionId1", StartsAt: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)},
			{Id: "sessionId2", StartsAt: time.Date(2024, 1, 2, 12, 30, 0, 0, time.UTC)},
			{Id: "sessionId3", StartsAt: time.Date(2024, 1, 3, 12, 30, 0, 0, time.UTC)},
		}, nil)
		sessionCloser.EXPECT().CloseSession("sessionId1").Return(nil)
		sessionCloser.EXPECT().CloseSession("sessionId2").Return(nil)
		// sessionId3 is not expired yet.
		result, err := executor.CloseExpiredSessions()

		assert.Equal(t, CloseExpiredSessionsResult{
			SucceededSessionIds: []string{"sessionId1", "sessionId2"},
		}, result)
		assert.NoError(t, err)
	})
}
