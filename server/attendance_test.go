package server

import (
	"errors"
	"fmt"
	"rush/attendance"
	"rush/session"
	"rush/user"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestGetHalfYearAttendance(t *testing.T) {
	t.Run("Fails if it fails to get users", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockUserRepo := NewMockuserRepo(controller)
		server := New(nil, nil, mockUserRepo, nil, nil, nil, nil, nil, nil, nil, nil)

		mockUserRepo.EXPECT().GetAllActive().Return(nil, errors.New("failed to get active users"))
		_, err := server.GetHalfYearAttendance()

		assert.Equal(t, newInternalServerError(fmt.Errorf("failed to get users: %w",
			errors.New("failed to get active users"))), err)
	})

	t.Run("Fails if it fails to get attendances", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockUserRepo := NewMockuserRepo(controller)
		mockAttendanceRepo := NewMockattendanceRepo(controller)
		server := New(nil, nil, mockUserRepo, nil, nil, nil, nil, nil, mockAttendanceRepo, nil, nil)

		mockUserRepo.EXPECT().GetAllActive().Return([]user.User{
			{Id: "1", Name: "김건", Generation: 9},
		}, nil)
		mockAttendanceRepo.EXPECT().GetAll().Return(nil, errors.New("failed to get attendances"))
		_, err := server.GetHalfYearAttendance()

		assert.Equal(t, newInternalServerError(fmt.Errorf("failed to get all attendances: %w",
			errors.New("failed to get attendances"))), err)
	})

	t.Run("Return the active users, sessions and their attendances after sorting", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockUserRepo := NewMockuserRepo(controller)
		mockAttendanceRepo := NewMockattendanceRepo(controller)
		server := New(nil, nil, mockUserRepo, nil, nil, nil, nil, nil, mockAttendanceRepo, nil, nil)

		// Different generations, different names for the same generation.
		mockUserRepo.EXPECT().GetAllActive().Return([]user.User{
			{Id: "1", Name: "김건", Generation: 9, ExternalName: "김건ExName", IsActive: true},
			{Id: "2", Name: "양현우", Generation: 8, ExternalName: "양현우", IsActive: true},
			{Id: "3", Name: "강민경", Generation: 8, ExternalName: "강민경", IsActive: true},
			{Id: "4", Name: "어떤10기", Generation: 10, ExternalName: "어떤10기ExName", IsActive: true},
			{Id: "5", Name: "어떤다른사람", Generation: 10, ExternalName: "어떤다른사람", IsActive: false},
		}, nil)
		// Sessions with different startedAt.
		mockAttendanceRepo.EXPECT().GetAll().Return([]attendance.Attendance{
			{
				Id:               "attendance_id_1",
				SessionId:        "session_id_1",
				SessionName:      "연트",
				SessionScore:     2,
				SessionStartedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				UserId:           "1",
				UserExternalName: "김건ExName",
				UserGeneration:   9,
			},
			{
				Id:               "attendance_id_2",
				SessionId:        "session_id_1",
				SessionName:      "연트",
				SessionScore:     2,
				SessionStartedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				UserId:           "2",
				UserExternalName: "unknown",
				UserGeneration:   9,
			},
			{
				Id:               "attendance_id_3",
				SessionId:        "session_id_1",
				SessionName:      "연트",
				SessionScore:     2,
				SessionStartedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				UserId:           "3",
				UserExternalName: "강민경",
				UserGeneration:   8,
			},
			{
				Id:               "attendance_id_4",
				SessionId:        "session_id_2",
				SessionName:      "다른 세션",
				SessionScore:     2,
				SessionStartedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				UserId:           "1",
				UserExternalName: "김건ExName",
				UserGeneration:   9,
			},
		}, nil)

		halfYearAttendance, err := server.GetHalfYearAttendance()
		assert.NoError(t, err)
		assert.Equal(t, HalfYearAttendace{
			// Should be sorted by startedAt.
			Sessions: []sessionForAttendance{
				{Id: "session_id_2", Name: "다른 세션", StartedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				{Id: "session_id_1", Name: "연트", StartedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			// Should be sorted by generation and then by name. Active users only.
			Users: []userForAttendance{
				{Id: "3", Name: "강민경", Generation: 8},
				{Id: "2", Name: "양현우", Generation: 8},
				// Use the actual name, not external name.
				{Id: "1", Name: "김건", Generation: 9},
				{Id: "4", Name: "어떤10기", Generation: 10},
			},
			Attendances: []Attendance{
				{Id: "attendance_id_1", SessionId: "session_id_1", SessionName: "연트",
					SessionScore:     2,
					SessionStartedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					UserId:           "1", UserExternalName: "김건ExName", UserGeneration: 9,
				},
				{Id: "attendance_id_2", SessionId: "session_id_1", SessionName: "연트",
					SessionScore:     2,
					SessionStartedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					UserId:           "2", UserExternalName: "unknown", UserGeneration: 9,
				},
				{Id: "attendance_id_3", SessionId: "session_id_1", SessionName: "연트",
					SessionScore:     2,
					SessionStartedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					UserId:           "3", UserExternalName: "강민경", UserGeneration: 8,
				},
				{Id: "attendance_id_4", SessionId: "session_id_2", SessionName: "다른 세션",
					SessionScore:     2,
					SessionStartedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UserId:           "1", UserExternalName: "김건ExName", UserGeneration: 9,
				},
			},
		}, halfYearAttendance)
	})
}

func TestMarkUsersAsPresent(t *testing.T) {
	t.Run("Fails if it fails to get session", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockSessionRepo := NewMocksessionRepo(controller)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().Get("session_id").Return(session.Session{}, errors.New("failed to get session"))
		err := server.MarkUsersAsPresent("session_id", []string{"user_id"}, false /* =forceApply */, "caller")

		assert.Equal(t, newNotFoundError(fmt.Errorf("failed to get session: %w",
			errors.New("failed to get session"))), err)
	})

	t.Run("Fails if the session is already closed", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockSessionRepo := NewMocksessionRepo(controller)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		// TODO(#54): Fix it to mock the returned session's CanApplyAttendanceManually as it's hard to test.
		mockSessionRepo.EXPECT().Get("session_id").Return(session.Session{
			Id:               "session_id",
			AttendanceStatus: session.AttendanceStatusApplied,
		}, nil)
		err := server.MarkUsersAsPresent("session_id", []string{"user_id"}, false /* =forceApply */, "caller")

		assert.Equal(t, newBadRequestError(errors.New("session is already closed")), err)
	})

	t.Run("Fails if it fails to get attendances", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockSessionRepo := NewMocksessionRepo(controller)
		mockAttendanceRepo := NewMockattendanceRepo(controller)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, mockAttendanceRepo, nil, nil)

		mockSessionRepo.EXPECT().Get("session_id").Return(session.Session{
			Id:               "session_id",
			AttendanceStatus: session.AttendanceStatusNotAppliedYet,
		}, nil)
		mockAttendanceRepo.EXPECT().FindBySessionId("session_id").Return(nil, errors.New("failed to get attendances"))
		err := server.MarkUsersAsPresent("session_id", []string{"user_id"}, false /* =forceApply */, "caller")

		assert.Equal(t, newInternalServerError(fmt.Errorf("failed to get attendances: %w",
			errors.New("failed to get attendances"))), err)
	})

	t.Run("Fails if it fails to get users", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockSessionRepo := NewMocksessionRepo(controller)
		mockAttendanceRepo := NewMockattendanceRepo(controller)
		mockUserRepo := NewMockuserRepo(controller)
		server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, nil, nil, mockAttendanceRepo, nil, nil)

		mockSessionRepo.EXPECT().Get("session_id").Return(session.Session{
			Id:               "session_id",
			AttendanceStatus: session.AttendanceStatusNotAppliedYet,
		}, nil)
		mockAttendanceRepo.EXPECT().FindBySessionId("session_id").Return([]attendance.Attendance{
			{Id: "attendance_id_1", SessionId: "session_id", UserId: "user_id_1"},
		}, nil)
		mockUserRepo.EXPECT().GetAllActive().Return(nil, errors.New("failed to get users"))
		err := server.MarkUsersAsPresent("session_id", []string{"user_id"}, false /* =forceApply */, "caller")

		assert.Equal(t, newInternalServerError(fmt.Errorf("failed to get users: %w",
			errors.New("failed to get users"))), err)
	})

	t.Run("Fails if it tries to mark inactive users as present", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockSessionRepo := NewMocksessionRepo(controller)
		mockAttendanceRepo := NewMockattendanceRepo(controller)
		mockUserRepo := NewMockuserRepo(controller)
		server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, nil, nil, mockAttendanceRepo, nil, nil)

		mockSessionRepo.EXPECT().Get("session_id").Return(session.Session{
			Id:               "session_id",
			AttendanceStatus: session.AttendanceStatusNotAppliedYet,
		}, nil)
		mockAttendanceRepo.EXPECT().FindBySessionId("session_id").Return([]attendance.Attendance{
			{Id: "attendance_id_1", SessionId: "session_id", UserId: "user_id_3"},
		}, nil)
		mockUserRepo.EXPECT().GetAllActive().Return([]user.User{
			{Id: "user_id_1", IsActive: false},
			{Id: "user_id_2", IsActive: true},
			{Id: "user_id_3", IsActive: true},
		}, nil)
		err := server.MarkUsersAsPresent("session_id", []string{"user_id_1", "user_id_2"}, false /* =forceApply */, "caller")

		assert.Equal(t, newBadRequestError(
			fmt.Errorf(
				"it received 2 user IDs (user_id_1,user_id_2) where 2 users (user_id_1,user_id_2) are not attended yet but only 1 users (user_id_2) are active among them",
			)), err)
	})

	t.Run("Fails if it fails to insert attendances", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockSessionRepo := NewMocksessionRepo(controller)
		mockAttendanceRepo := NewMockattendanceRepo(controller)
		mockUserRepo := NewMockuserRepo(controller)
		server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, nil, nil, mockAttendanceRepo, nil, nil)

		mockSessionRepo.EXPECT().Get("session_id").Return(session.Session{
			Id:               "session_id",
			AttendanceStatus: session.AttendanceStatusNotAppliedYet,
		}, nil)
		mockAttendanceRepo.EXPECT().FindBySessionId("session_id").Return([]attendance.Attendance{
			{Id: "attendance_id_1", SessionId: "session_id", UserId: "user_id_1"},
		}, nil)
		mockUserRepo.EXPECT().GetAllActive().Return([]user.User{
			{Id: "user_id_1", IsActive: true},
		}, nil)
		mockAttendanceRepo.EXPECT().BulkInsert(gomock.Any()).Return(errors.New("failed to insert attendances"))
		err := server.MarkUsersAsPresent("session_id", []string{"user_id_1"}, false /* =forceApply */, "caller")

		assert.Equal(t, newInternalServerError(fmt.Errorf("failed to bulk insert attendances: %w",
			errors.New("failed to insert attendances"))), err)
	})

	t.Run("Fails if it fails to mark the session as attendance applied", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockSessionRepo := NewMocksessionRepo(controller)
		mockAttendanceRepo := NewMockattendanceRepo(controller)
		mockUserRepo := NewMockuserRepo(controller)
		mockOpenSessionRepo := NewMockopenSessionRepo(controller)
		server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, mockOpenSessionRepo, nil, mockAttendanceRepo, nil, nil)

		mockSessionRepo.EXPECT().Get("session_id").Return(session.Session{
			Id:               "session_id",
			AttendanceStatus: session.AttendanceStatusNotAppliedYet,
		}, nil)
		mockAttendanceRepo.EXPECT().FindBySessionId("session_id").Return([]attendance.Attendance{
			{Id: "attendance_id_1", SessionId: "session_id", UserId: "user_id_1"},
		}, nil)
		mockUserRepo.EXPECT().GetAllActive().Return([]user.User{
			{Id: "user_id_1", IsActive: true},
		}, nil)
		mockAttendanceRepo.EXPECT().BulkInsert(gomock.Any()).Return(nil)
		mockOpenSessionRepo.EXPECT().MarkAsAttendanceApplied("session_id").Return(errors.New("failed to mark the session as attendance applied"))
		err := server.MarkUsersAsPresent("session_id", []string{"user_id_1"}, false /* =forceApply */, "caller")

		assert.Equal(t, newInternalServerError(fmt.Errorf("failed to close the session: %w",
			errors.New("failed to mark the session as attendance applied"))), err)
	})

	t.Run("Apply attendance for the active users who are not attended yet", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockSessionRepo := NewMocksessionRepo(controller)
		mockAttendanceRepo := NewMockattendanceRepo(controller)
		mockUserRepo := NewMockuserRepo(controller)
		mockOpenSessionRepo := NewMockopenSessionRepo(controller)
		mockClock := clock.NewMock()
		server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, mockOpenSessionRepo, nil, mockAttendanceRepo, nil, mockClock)

		mockClock.Set(time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC))
		mockSessionRepo.EXPECT().Get("session_id").Return(session.Session{
			Id:               "session_id",
			AttendanceStatus: session.AttendanceStatusNotAppliedYet,
			Name:             "session_name",
			Score:            2,
			StartsAt:         time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		}, nil)
		mockAttendanceRepo.EXPECT().FindBySessionId("session_id").Return([]attendance.Attendance{
			{Id: "attendance_id_1", SessionId: "session_id", UserId: "user_id_1"},
		}, nil)
		mockUserRepo.EXPECT().GetAllActive().Return([]user.User{
			{Id: "user_id_1", IsActive: true, ExternalName: "user_external_name_1",
				Generation: 9, Name: "user_name_1"},
			{Id: "user_id_2", IsActive: true, ExternalName: "user_external_name_2",
				Generation: 9, Name: "user_name_2"},
			{Id: "user_id_3", IsActive: false, ExternalName: "user_external_name_3",
				Generation: 9, Name: "user_name_3"},
		}, nil)
		mockAttendanceRepo.EXPECT().BulkInsert([]attendance.AddAttendanceReq{
			{
				SessionId:        "session_id",
				SessionName:      "session_name",
				SessionScore:     2,
				SessionStartedAt: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				UserId:           "user_id_2",
				UserExternalName: "user_external_name_2",
				UserGeneration:   9,
				UserJoinedAt:     time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				CreatedBy:        "caller",
			},
		}).Return(nil)
		mockOpenSessionRepo.EXPECT().MarkAsAttendanceApplied("session_id").Return(nil)

		err := server.MarkUsersAsPresent("session_id", []string{"user_id_1", "user_id_2"}, false /* =forceApply */, "caller")
		assert.NoError(t, err)
	})

	t.Run("Apply attendance with force apply even if the session is already closed", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockSessionRepo := NewMocksessionRepo(controller)
		mockAttendanceRepo := NewMockattendanceRepo(controller)
		mockUserRepo := NewMockuserRepo(controller)
		mockOpenSessionRepo := NewMockopenSessionRepo(controller)
		mockClock := clock.NewMock()
		server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, mockOpenSessionRepo, nil, mockAttendanceRepo, nil, mockClock)

		mockClock.Set(time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC))
		mockSessionRepo.EXPECT().Get("session_id").Return(session.Session{
			Id:               "session_id",
			AttendanceStatus: session.AttendanceStatusNotAppliedYet,
			Name:             "session_name",
			Score:            2,
			StartsAt:         time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		}, nil)
		mockAttendanceRepo.EXPECT().FindBySessionId("session_id").Return([]attendance.Attendance{}, nil)
		mockUserRepo.EXPECT().GetAllActive().Return([]user.User{
			{Id: "user_id_1", IsActive: true, ExternalName: "user_external_name_1",
				Generation: 9, Name: "user_name_1"},
		}, nil)
		mockAttendanceRepo.EXPECT().BulkInsert([]attendance.AddAttendanceReq{
			{
				SessionId:        "session_id",
				SessionName:      "session_name",
				SessionScore:     2,
				SessionStartedAt: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				UserId:           "user_id_1",
				UserExternalName: "user_external_name_1",
				UserGeneration:   9,
				UserJoinedAt:     time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				CreatedBy:        "caller",
				ForceApply:       true,
			},
		}).Return(nil)
		mockOpenSessionRepo.EXPECT().MarkAsAttendanceApplied("session_id").Return(nil)

		err := server.MarkUsersAsPresent("session_id", []string{"user_id_1"}, true /* =forceApply */, "caller")
		assert.NoError(t, err)
	})
}
