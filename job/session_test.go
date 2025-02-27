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
		sessionAttendanceApplier := NewMocksessionAttendanceApplier(controller)
		mockLogger := NewMocklogger(controller)
		executor := NewExecutor(sessionGetter, sessionAttendanceApplier, mockLogger, clock.NewMock())

		sessionGetter.EXPECT().GetOpenSessionsWithForm().Return([]session.Session{}, assert.AnError)
		mockLogger.EXPECT().Errorw("Failed to get open sessions with form", "error", assert.AnError.Error())
		executor.CloseExpiredSessions()
	})

	t.Run("Fails if it fails to close ession", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionGetter := NewMocksessionGetter(controller)
		sessionAttendanceApplier := NewMocksessionAttendanceApplier(controller)
		mockLogger := NewMocklogger(controller)
		clock := clock.NewMock()
		executor := NewExecutor(sessionGetter, sessionAttendanceApplier, mockLogger, clock)

		clock.Set(time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC))
		sessionGetter.EXPECT().GetOpenSessionsWithForm().Return([]session.Session{
			{Id: "sessionId1", StartsAt: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)},
			{Id: "sessionId2", StartsAt: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)},
			{Id: "sessionId3", StartsAt: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)},
		}, nil)
		sessionAttendanceApplier.EXPECT().ApplyAttendanceByFormSubmissions("sessionId1", "session-attendance-syncer").Return(errors.New("error1"))
		sessionAttendanceApplier.EXPECT().ApplyAttendanceByFormSubmissions("sessionId2", "session-attendance-syncer").Return(nil)
		sessionAttendanceApplier.EXPECT().ApplyAttendanceByFormSubmissions("sessionId3", "session-attendance-syncer").Return(errors.New("error2"))
		mockLogger.EXPECT().Infow("Closed sessions", "session_ids", "sessionId2")
		mockLogger.EXPECT().Errorw("Failed to close sessions", "session_ids", "sessionId1, sessionId3", "errors", "error1, error2")
		executor.CloseExpiredSessions()
	})

	t.Run("Successfully close open and also expired sessions", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionGetter := NewMocksessionGetter(controller)
		sessionAttendanceApplier := NewMocksessionAttendanceApplier(controller)
		mockLogger := NewMocklogger(controller)
		clock := clock.NewMock()
		executor := NewExecutor(sessionGetter, sessionAttendanceApplier, mockLogger, clock)

		clock.Set(time.Date(2024, 1, 2, 12, 30, 0, 0, time.UTC))
		sessionGetter.EXPECT().GetOpenSessionsWithForm().Return([]session.Session{
			{Id: "sessionId1", StartsAt: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)},
			{Id: "sessionId2", StartsAt: time.Date(2024, 1, 2, 12, 30, 0, 0, time.UTC)},
			{Id: "sessionId3", StartsAt: time.Date(2024, 1, 3, 12, 30, 0, 0, time.UTC)},
		}, nil)
		sessionAttendanceApplier.EXPECT().ApplyAttendanceByFormSubmissions("sessionId1", "session-attendance-syncer").Return(nil)
		sessionAttendanceApplier.EXPECT().ApplyAttendanceByFormSubmissions("sessionId2", "session-attendance-syncer").Return(nil)
		// sessionId3 is not expired yet.
		mockLogger.EXPECT().Infow("Closed sessions", "session_ids", "sessionId1, sessionId2")
		executor.CloseExpiredSessions()
	})
}
