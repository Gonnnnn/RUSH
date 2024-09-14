package job

import (
	"rush/session"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestCloseSessions(t *testing.T) {
	t.Run("Fails if it fails to get open sessions", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionGetter := NewMocksessionGetter(controller)
		sessionCloser := NewMocksessionCloser(controller)
		executor := NewExecutor(sessionGetter, sessionCloser)

		sessionGetter.EXPECT().GetOpenSessions().Return([]session.Session{}, assert.AnError)
		result, err := executor.CloseSessions()

		assert.Equal(t, CloseSessionsResult{}, result)
		assert.EqualError(t, err, "failed to get open sessions: "+assert.AnError.Error())
	})

	t.Run("Fails if it fails to close ession", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionGetter := NewMocksessionGetter(controller)
		sessionCloser := NewMocksessionCloser(controller)
		executor := NewExecutor(sessionGetter, sessionCloser)

		sessionGetter.EXPECT().GetOpenSessions().Return([]session.Session{
			{Id: "sessionId1"},
			{Id: "sessionId2"},
			{Id: "sessionId3"},
		}, nil)
		sessionCloser.EXPECT().CloseSession("sessionId1").Return(assert.AnError)
		sessionCloser.EXPECT().CloseSession("sessionId2").Return(nil)
		sessionCloser.EXPECT().CloseSession("sessionId3").Return(assert.AnError)
		result, err := executor.CloseSessions()

		assert.Equal(t, CloseSessionsResult{
			SessionIdsSucceeded: []string{"sessionId2"},
			FailedSessionIds:    []string{"sessionId1", "sessionId3"},
			Errors:              []error{assert.AnError, assert.AnError},
		}, result)
		assert.EqualError(t, err, "failed to close some sessions")
	})

	t.Run("Succeeds if all sessions are closed", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionGetter := NewMocksessionGetter(controller)
		sessionCloser := NewMocksessionCloser(controller)
		executor := NewExecutor(sessionGetter, sessionCloser)

		sessionGetter.EXPECT().GetOpenSessions().Return([]session.Session{
			{Id: "sessionId1"},
			{Id: "sessionId2"},
		}, nil)
		sessionCloser.EXPECT().CloseSession("sessionId1").Return(nil)
		sessionCloser.EXPECT().CloseSession("sessionId2").Return(nil)
		result, err := executor.CloseSessions()

		assert.Equal(t, CloseSessionsResult{
			SessionIdsSucceeded: []string{"sessionId1", "sessionId2"},
		}, result)
		assert.NoError(t, err)
	})
}
