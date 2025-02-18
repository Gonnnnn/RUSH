package server

import (
	"errors"
	"fmt"
	"rush/attendance"
	"rush/session"
	"rush/user"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestAdminGetSession(t *testing.T) {
	t.Run("Returns not found error when session is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSessionRepo := NewMocksessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{}, session.ErrNotFound)
		dbSession, err := server.AdminGetSession("session-id")

		assert.Equal(t, SessionForAdmin{}, dbSession)
		assert.Equal(t, &NotFoundError{originalError: fmt.Errorf("failed to get session: %w", session.ErrNotFound)}, err)
	})

	t.Run("Returns internal server error when failed to get session", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSessionRepo := NewMocksessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{}, assert.AnError)
		dbSession, err := server.AdminGetSession("session-id")

		assert.Equal(t, SessionForAdmin{}, dbSession)
		assert.Equal(t, &InternalServerError{originalError: fmt.Errorf("failed to get session: %w", assert.AnError)}, err)
	})

	t.Run("Returns session when session is found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSessionRepo := NewMocksessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
			Id:               "session-id",
			Name:             "session-name",
			Description:      "session-description",
			CreatedBy:        "user-id",
			GoogleFormId:     "google-form-id",
			GoogleFormUri:    "google-form-uri",
			CreatedAt:        time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			Score:            1,
			AttendanceStatus: session.AttendanceStatusNotAppliedYet,
		}, nil)
		dbSession, err := server.AdminGetSession("session-id")

		assert.Equal(t, SessionForAdmin{
			Id:                  "session-id",
			Name:                "session-name",
			Description:         "session-description",
			CreatedBy:           "user-id",
			GoogleFormId:        "google-form-id",
			GoogleFormUri:       "google-form-uri",
			CreatedAt:           time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			StartsAt:            time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			Score:               1,
			AttendanceStatus:    session.AttendanceStatusNotAppliedYet,
			AttendanceAppliedBy: SessionAttendanceAppliedByUnspecified,
		}, dbSession)
		assert.NoError(t, err)
	})
}

func TestGetSession(t *testing.T) {
	t.Run("Returns not found error when session is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSessionRepo := NewMocksessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{}, session.ErrNotFound)
		dbSession, err := server.GetSession("session-id")

		assert.Equal(t, Session{}, dbSession)
		assert.Equal(t, &NotFoundError{originalError: fmt.Errorf("failed to get session: %w", session.ErrNotFound)}, err)
	})

	t.Run("Returns internal server error when failed to get session", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSessionRepo := NewMocksessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{}, assert.AnError)
		dbSession, err := server.GetSession("session-id")

		assert.Equal(t, Session{}, dbSession)
		assert.Equal(t, &InternalServerError{originalError: fmt.Errorf("failed to get session: %w", assert.AnError)}, err)
	})

	t.Run("Returns session when session is found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSessionRepo := NewMocksessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
			Id:               "session-id",
			Name:             "session-name",
			Description:      "session-description",
			CreatedBy:        "user-id",
			GoogleFormId:     "google-form-id",
			GoogleFormUri:    "google-form-uri",
			CreatedAt:        time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			Score:            1,
			AttendanceStatus: session.AttendanceStatusNotAppliedYet,
		}, nil)
		dbSession, err := server.GetSession("session-id")

		assert.Equal(t, Session{
			Id:          "session-id",
			Name:        "session-name",
			Description: "session-description",
			CreatedBy:   "user-id",
			CreatedAt:   time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			StartsAt:    time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			Score:       1,
		}, dbSession)
		assert.NoError(t, err)
	})
}

func TestAdminListSessions(t *testing.T) {
	t.Run("Returns internal server error when failed to list sessions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSessionRepo := NewMocksessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().List(1, 2).Return(nil, assert.AnError)
		listResult, err := server.AdminListSessions(1, 2)

		assert.Nil(t, listResult)
		assert.Equal(t, &InternalServerError{originalError: fmt.Errorf("failed to list sessions: %w", assert.AnError)}, err)
	})

	t.Run("Returns sessions when successfully fetches sessions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSessionRepo := NewMocksessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().List(1, 2).Return(
			&session.ListResult{
				Sessions: []session.Session{
					{
						Id:               "session-id",
						Name:             "session-name",
						Description:      "session-description",
						CreatedBy:        "user-id",
						GoogleFormId:     "google-form-id",
						GoogleFormUri:    "google-form-uri",
						CreatedAt:        time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
						StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
						Score:            1,
						AttendanceStatus: session.AttendanceStatusNotAppliedYet,
					},
					{
						Id:               "session-id2",
						Name:             "session-name2",
						Description:      "session-description2",
						CreatedBy:        "user-id2",
						GoogleFormId:     "google-form-id2",
						GoogleFormUri:    "google-form-uri2",
						CreatedAt:        time.Date(2024, 1, 1, 21, 0, 0, 0, time.UTC),
						StartsAt:         time.Date(2024, 1, 1, 21, 0, 0, 0, time.UTC),
						Score:            2,
						AttendanceStatus: session.AttendanceStatusNotAppliedYet,
					},
				},
				IsEnd:      true,
				TotalCount: 2,
			}, nil)
		listResult, err := server.AdminListSessions(1, 2)

		assert.Equal(t, &AdminListSessionsResult{
			Sessions: []SessionForAdmin{
				{
					Id:                  "session-id",
					Name:                "session-name",
					Description:         "session-description",
					CreatedBy:           "user-id",
					GoogleFormId:        "google-form-id",
					GoogleFormUri:       "google-form-uri",
					CreatedAt:           time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
					StartsAt:            time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
					Score:               1,
					AttendanceStatus:    session.AttendanceStatusNotAppliedYet,
					AttendanceAppliedBy: SessionAttendanceAppliedByUnspecified,
				},
				{
					Id:                  "session-id2",
					Name:                "session-name2",
					Description:         "session-description2",
					CreatedBy:           "user-id2",
					GoogleFormId:        "google-form-id2",
					GoogleFormUri:       "google-form-uri2",
					CreatedAt:           time.Date(2024, 1, 1, 21, 0, 0, 0, time.UTC),
					StartsAt:            time.Date(2024, 1, 1, 21, 0, 0, 0, time.UTC),
					Score:               2,
					AttendanceStatus:    session.AttendanceStatusNotAppliedYet,
					AttendanceAppliedBy: SessionAttendanceAppliedByUnspecified,
				},
			},
			IsEnd:      true,
			TotalCount: 2,
		}, listResult)
		assert.NoError(t, err)
	})
}

func TestListSessions(t *testing.T) {
	t.Run("Returns internal server error when failed to list sessions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSessionRepo := NewMocksessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().List(1, 2).Return(nil, assert.AnError)
		listResult, err := server.ListSessions(1, 2)

		assert.Nil(t, listResult)
		assert.Equal(t, &InternalServerError{originalError: fmt.Errorf("failed to list sessions: %w", assert.AnError)}, err)
	})

	t.Run("Returns sessions when successfully fetches sessions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSessionRepo := NewMocksessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().List(1, 2).Return(
			&session.ListResult{
				Sessions: []session.Session{
					{
						Id:               "session-id",
						Name:             "session-name",
						Description:      "session-description",
						CreatedBy:        "user-id",
						GoogleFormId:     "google-form-id",
						GoogleFormUri:    "google-form-uri",
						CreatedAt:        time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
						StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
						Score:            1,
						AttendanceStatus: session.AttendanceStatusNotAppliedYet,
					},
					{
						Id:               "session-id2",
						Name:             "session-name2",
						Description:      "session-description2",
						CreatedBy:        "user-id2",
						GoogleFormId:     "google-form-id2",
						GoogleFormUri:    "google-form-uri2",
						CreatedAt:        time.Date(2024, 1, 1, 21, 0, 0, 0, time.UTC),
						StartsAt:         time.Date(2024, 1, 1, 21, 0, 0, 0, time.UTC),
						Score:            2,
						AttendanceStatus: session.AttendanceStatusNotAppliedYet,
					},
				},
				IsEnd:      true,
				TotalCount: 2,
			}, nil)
		listResult, err := server.ListSessions(1, 2)

		assert.Equal(t, &ListSessionsResult{
			Sessions: []Session{
				{
					Id:          "session-id",
					Name:        "session-name",
					Description: "session-description",
					CreatedBy:   "user-id",
					CreatedAt:   time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
					StartsAt:    time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
					Score:       1,
				},
				{
					Id:          "session-id2",
					Name:        "session-name2",
					Description: "session-description2",
					CreatedBy:   "user-id2",
					CreatedAt:   time.Date(2024, 1, 1, 21, 0, 0, 0, time.UTC),
					StartsAt:    time.Date(2024, 1, 1, 21, 0, 0, 0, time.UTC),
					Score:       2,
				},
			},
			IsEnd:      true,
			TotalCount: 2,
		}, listResult)
		assert.NoError(t, err)
	})
}

func TestAddSession(t *testing.T) {
	t.Run("Returns internal server error when failed to add session", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSessionRepo := NewMocksessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().Add("session-name", "session-description", "user-id", time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC), 1).Return("", assert.AnError)
		id, err := server.AddSession("session-name", "session-description", "user-id", time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC), 1)

		assert.Equal(t, "", id)
		assert.Equal(t, &InternalServerError{originalError: fmt.Errorf("failed to add session: %w", assert.AnError)}, err)
	})

	t.Run("Returns session id when successfully adds session", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSessionRepo := NewMocksessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

		mockSessionRepo.EXPECT().Add("session-name", "session-description", "user-id", time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC), 1).Return("session-id", nil)
		id, err := server.AddSession("session-name", "session-description", "user-id", time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC), 1)

		assert.Equal(t, "session-id", id)
		assert.NoError(t, err)
	})
}

func TestDeleteSession(t *testing.T) {
	t.Run("Returns internal server error when failed to delete session", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, nil, mockOpenSessionRepo, nil, nil, nil, nil)

		mockOpenSessionRepo.EXPECT().DeleteOpenSession("session-id").Return(assert.AnError)
		err := server.DeleteSession("session-id")

		assert.Equal(t, &InternalServerError{originalError: fmt.Errorf("failed to delete session: %w", assert.AnError)}, err)
	})

	t.Run("Returns nil when successfully deletes open session", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
		server := New(nil, nil, nil, nil, nil, nil, mockOpenSessionRepo, nil, nil, nil, nil)

		mockOpenSessionRepo.EXPECT().DeleteOpenSession("session-id").Return(nil)
		err := server.DeleteSession("session-id")

		assert.NoError(t, err)
	})
}

func TestApplyAttendanceByFormSubmissions(t *testing.T) {
	t.Run("Failures", func(t *testing.T) {
		t.Run("Returns not found error when session is not found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{}, errors.New("session not found"))
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			var notFoundError *NotFoundError
			assert.ErrorAs(t, err, &notFoundError)
			assert.EqualError(t, notFoundError.originalError, "failed to get session: session not found")
		})

		t.Run("Returns bad request error when session is already closed", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				AttendanceStatus: session.AttendanceStatusApplied,
			}, nil)
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			var badRequestError *BadRequestError
			assert.ErrorAs(t, err, &badRequestError)
			assert.EqualError(t, badRequestError.originalError, "session is already closed")
		})

		t.Run("Returns bad request error when session cannot apply attendance by form submissions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, nil, nil, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				GoogleFormId:     "",
				GoogleFormUri:    "",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
			}, nil)
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			var badRequestError *BadRequestError
			assert.ErrorAs(t, err, &badRequestError)
			assert.EqualError(t, badRequestError.originalError, "session cannot apply attendance by form submissions")
		})

		t.Run("Returns internal server error when failed to get form submissions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			mockAttendanceFormHandler := NewMockattendanceFormHandler(ctrl)
			server := New(nil, nil, nil, nil, nil, mockSessionRepo, nil, mockAttendanceFormHandler, nil, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				GoogleFormId:     "form-id",
				GoogleFormUri:    "form-uri",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
			}, nil)
			mockAttendanceFormHandler.EXPECT().GetSubmissions("form-id").Return(nil, errors.New("failed to get form submissions"))
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			var internalServerError *InternalServerError
			assert.ErrorAs(t, err, &internalServerError)
			assert.EqualError(t, internalServerError.originalError, "failed to get form submissions: failed to get form submissions")
		})

		t.Run("Marks attendance as ignored when there are no submissions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			mockAttendanceFormHandler := NewMockattendanceFormHandler(ctrl)
			mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
			server := New(nil, nil, nil, nil, nil, mockSessionRepo, mockOpenSessionRepo,
				mockAttendanceFormHandler, nil, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				GoogleFormId:     "form-id",
				GoogleFormUri:    "form-uri",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
			}, nil)
			mockAttendanceFormHandler.EXPECT().GetSubmissions("form-id").Return([]attendance.FormSubmission{}, nil)
			mockOpenSessionRepo.EXPECT().MarkAttendanceIsIgnored("session-id", "no form submissions").Return(nil)
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			assert.NoError(t, err)
		})

		t.Run("Returns internal server error when failed to mark attendance as ignored", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			mockAttendanceFormHandler := NewMockattendanceFormHandler(ctrl)
			mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
			server := New(nil, nil, nil, nil, nil, mockSessionRepo, mockOpenSessionRepo,
				mockAttendanceFormHandler, nil, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				GoogleFormId:     "form-id",
				GoogleFormUri:    "form-uri",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
			}, nil)
			mockAttendanceFormHandler.EXPECT().GetSubmissions("form-id").Return([]attendance.FormSubmission{}, nil)
			mockOpenSessionRepo.EXPECT().MarkAttendanceIsIgnored("session-id", "no form submissions").Return(errors.New("failed to mark attendance as ignored"))
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			var internalServerError *InternalServerError
			assert.ErrorAs(t, err, &internalServerError)
			assert.EqualError(t, internalServerError.originalError, "failed to mark the session's attendance as ignored: failed to mark attendance as ignored")
		})

		t.Run("Returns internal server error when failed to get users by external names", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			mockAttendanceFormHandler := NewMockattendanceFormHandler(ctrl)
			mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)
			server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, mockOpenSessionRepo,
				mockAttendanceFormHandler, nil, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				GoogleFormId:     "form-id",
				GoogleFormUri:    "form-uri",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
				StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			}, nil)
			mockAttendanceFormHandler.EXPECT().GetSubmissions("form-id").Return([]attendance.FormSubmission{
				{
					UserExternalName: "user-external-name-1",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 0, 0, 0, time.UTC),
				},
				{
					UserExternalName: "user-external-name-2",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 59, 59, 0, time.UTC),
				},
			}, nil)
			mockUserRepo.EXPECT().GetAllByExternalNames([]string{"user-external-name-1", "user-external-name-2"}).
				Return(nil, errors.New("failed to get users by external names"))
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			var internalServerError *InternalServerError
			assert.ErrorAs(t, err, &internalServerError)
			assert.EqualError(t, internalServerError.originalError, "failed to get users by external names: failed to get users by external names")
		})

		t.Run("Returns internal server error when there are users not found but failed to mark attendance as ignored", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			mockAttendanceFormHandler := NewMockattendanceFormHandler(ctrl)
			mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)
			server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, mockOpenSessionRepo,
				mockAttendanceFormHandler, nil, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				GoogleFormId:     "form-id",
				GoogleFormUri:    "form-uri",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
				StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			}, nil)
			mockAttendanceFormHandler.EXPECT().GetSubmissions("form-id").Return([]attendance.FormSubmission{
				{
					UserExternalName: "user-external-name-1",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 0, 0, 0, time.UTC),
				},
				{
					UserExternalName: "user-external-name-2",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 59, 59, 0, time.UTC),
				},
			}, nil)
			mockUserRepo.EXPECT().GetAllByExternalNames([]string{"user-external-name-1", "user-external-name-2"}).
				Return([]user.User{
					// Only user-external-name-1 is found.
					{
						Id:           "user-id-1",
						ExternalName: "user-external-name-1",
					},
				}, nil)
			mockOpenSessionRepo.EXPECT().MarkAttendanceIsIgnored("session-id", "some users (user-external-name-2) were not found although there are form submissions").
				Return(errors.New("failed to mark attendance as ignored"))
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			var internalServerError *InternalServerError
			assert.ErrorAs(t, err, &internalServerError)
			assert.EqualError(t, internalServerError.originalError, "some users (user-external-name-2) were not found although there are form submissions and it has failed to mark the session's attendance as ignored: failed to mark attendance as ignored")
		})

		t.Run("Returns internal server error when there are not-found users and mark attendance as ignored successfully", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			mockAttendanceFormHandler := NewMockattendanceFormHandler(ctrl)
			mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)
			server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, mockOpenSessionRepo,
				mockAttendanceFormHandler, nil, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				GoogleFormId:     "form-id",
				GoogleFormUri:    "form-uri",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
				StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			}, nil)
			mockAttendanceFormHandler.EXPECT().GetSubmissions("form-id").Return([]attendance.FormSubmission{
				{
					UserExternalName: "user-external-name-1",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 0, 0, 0, time.UTC),
				},
				{
					UserExternalName: "user-external-name-2",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 59, 59, 0, time.UTC),
				},
			}, nil)
			mockUserRepo.EXPECT().GetAllByExternalNames([]string{"user-external-name-1", "user-external-name-2"}).
				Return([]user.User{
					{
						Id:           "user-id-1",
						ExternalName: "user-external-name-1",
					},
				}, nil)
			mockOpenSessionRepo.EXPECT().MarkAttendanceIsIgnored("session-id", "some users (user-external-name-2) were not found although there are form submissions").
				Return(nil)
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			var internalServerError *InternalServerError
			assert.ErrorAs(t, err, &internalServerError)
			assert.EqualError(t, internalServerError.originalError, "some users (user-external-name-2) were not found although there are form submissions")
		})

		t.Run("Returns internal server error when failed to bulk insert attendances", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			mockAttendanceFormHandler := NewMockattendanceFormHandler(ctrl)
			mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)
			mockAttendanceRepo := NewMockattendanceRepo(ctrl)
			server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, mockOpenSessionRepo,
				mockAttendanceFormHandler, mockAttendanceRepo, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				GoogleFormId:     "form-id",
				GoogleFormUri:    "form-uri",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
				StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			}, nil)
			mockAttendanceFormHandler.EXPECT().GetSubmissions("form-id").Return([]attendance.FormSubmission{
				{
					UserExternalName: "user-external-name-1",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 0, 0, 0, time.UTC),
				},
			}, nil)
			mockUserRepo.EXPECT().GetAllByExternalNames([]string{"user-external-name-1"}).
				Return([]user.User{
					{
						Id:           "user-id-1",
						ExternalName: "user-external-name-1",
					},
				}, nil)
			mockAttendanceRepo.EXPECT().BulkInsert(gomock.Any()).Return(errors.New("failed to bulk insert attendances"))
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			var internalServerError *InternalServerError
			assert.ErrorAs(t, err, &internalServerError)
			assert.EqualError(t, internalServerError.originalError, "failed to bulk insert attendance: failed to bulk insert attendances")
		})

		t.Run("Returns internal server error when failed to mark open session as attendance applied", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			mockAttendanceFormHandler := NewMockattendanceFormHandler(ctrl)
			mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)
			mockAttendanceRepo := NewMockattendanceRepo(ctrl)
			server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, mockOpenSessionRepo,
				mockAttendanceFormHandler, mockAttendanceRepo, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				GoogleFormId:     "form-id",
				GoogleFormUri:    "form-uri",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
				StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			}, nil)
			mockAttendanceFormHandler.EXPECT().GetSubmissions("form-id").Return([]attendance.FormSubmission{
				{
					UserExternalName: "user-external-name-1",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 0, 0, 0, time.UTC),
				},
			}, nil)
			mockUserRepo.EXPECT().GetAllByExternalNames([]string{"user-external-name-1"}).
				Return([]user.User{
					{
						Id:           "user-id-1",
						ExternalName: "user-external-name-1",
					},
				}, nil)
			mockAttendanceRepo.EXPECT().BulkInsert(gomock.Any()).Return(nil)
			mockOpenSessionRepo.EXPECT().MarkAsAttendanceApplied("session-id").Return(errors.New("failed to close open session"))
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			var internalServerError *InternalServerError
			assert.ErrorAs(t, err, &internalServerError)
			assert.EqualError(t, internalServerError.originalError, "failed to close open session: failed to close open session")
		})
	})

	t.Run("Success", func(t *testing.T) {
		t.Run("When there is zero submission", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			mockAttendanceFormHandler := NewMockattendanceFormHandler(ctrl)
			mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)
			mockAttendanceRepo := NewMockattendanceRepo(ctrl)
			server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, mockOpenSessionRepo,
				mockAttendanceFormHandler, mockAttendanceRepo, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				Name:             "session-name",
				GoogleFormId:     "form-id",
				GoogleFormUri:    "form-uri",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
				Score:            2,
				StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			}, nil)
			mockAttendanceFormHandler.EXPECT().GetSubmissions("form-id").Return([]attendance.FormSubmission{}, nil)
			mockOpenSessionRepo.EXPECT().MarkAttendanceIsIgnored("session-id", "no form submissions").Return(nil)
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			assert.NoError(t, err)
		})

		t.Run("When there is one submission", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			mockAttendanceFormHandler := NewMockattendanceFormHandler(ctrl)
			mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)
			mockAttendanceRepo := NewMockattendanceRepo(ctrl)
			server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, mockOpenSessionRepo,
				mockAttendanceFormHandler, mockAttendanceRepo, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				Name:             "session-name",
				GoogleFormId:     "form-id",
				GoogleFormUri:    "form-uri",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
				Score:            2,
				StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			}, nil)
			mockAttendanceFormHandler.EXPECT().GetSubmissions("form-id").Return([]attendance.FormSubmission{
				{
					UserExternalName: "user-external-name-1",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 0, 0, 0, time.UTC),
				},
			}, nil)
			mockUserRepo.EXPECT().GetAllByExternalNames([]string{"user-external-name-1"}).
				Return([]user.User{
					{
						Id:           "user-id-1",
						ExternalName: "user-external-name-1",
						Generation:   1,
					},
				}, nil)
			mockAttendanceRepo.EXPECT().BulkInsert([]attendance.AddAttendanceReq{
				{
					SessionId:        "session-id",
					SessionName:      "session-name",
					SessionScore:     2,
					SessionStartedAt: time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
					UserId:           "user-id-1",
					UserExternalName: "user-external-name-1",
					UserGeneration:   1,
					UserJoinedAt:     time.Date(2024, 1, 1, 19, 0, 0, 0, time.UTC),
					CreatedBy:        "caller-id",
				},
			}).Return(nil)
			mockOpenSessionRepo.EXPECT().MarkAsAttendanceApplied("session-id").Return(nil)
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			assert.NoError(t, err)
		})

		t.Run("When there are multiple submissions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			mockAttendanceFormHandler := NewMockattendanceFormHandler(ctrl)
			mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)
			mockAttendanceRepo := NewMockattendanceRepo(ctrl)
			server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, mockOpenSessionRepo,
				mockAttendanceFormHandler, mockAttendanceRepo, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				Name:             "session-name",
				GoogleFormId:     "form-id",
				GoogleFormUri:    "form-uri",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
				Score:            2,
				StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			}, nil)
			mockAttendanceFormHandler.EXPECT().GetSubmissions("form-id").Return([]attendance.FormSubmission{
				{
					UserExternalName: "user-external-name-1",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 0, 0, 0, time.UTC),
				},
				{
					UserExternalName: "user-external-name-2",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 59, 59, 0, time.UTC),
				},
				{
					UserExternalName: "user-external-name-3",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 59, 59, 0, time.UTC),
				},
			}, nil)
			mockUserRepo.EXPECT().GetAllByExternalNames([]string{"user-external-name-1", "user-external-name-2", "user-external-name-3"}).
				Return([]user.User{
					{
						Id:           "user-id-1",
						ExternalName: "user-external-name-1",
						Generation:   1,
					},
					{
						Id:           "user-id-2",
						ExternalName: "user-external-name-2",
						Generation:   1,
					},
					{
						Id:           "user-id-3",
						ExternalName: "user-external-name-3",
						Generation:   1.5,
					},
				}, nil)
			mockAttendanceRepo.EXPECT().BulkInsert([]attendance.AddAttendanceReq{
				{
					SessionId:        "session-id",
					SessionName:      "session-name",
					SessionScore:     2,
					SessionStartedAt: time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
					UserId:           "user-id-1",
					UserExternalName: "user-external-name-1",
					UserGeneration:   1,
					UserJoinedAt:     time.Date(2024, 1, 1, 19, 0, 0, 0, time.UTC),
					CreatedBy:        "caller-id",
				},
				{
					SessionId:        "session-id",
					SessionName:      "session-name",
					SessionScore:     2,
					SessionStartedAt: time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
					UserId:           "user-id-2",
					UserExternalName: "user-external-name-2",
					UserGeneration:   1,
					UserJoinedAt:     time.Date(2024, 1, 1, 19, 59, 59, 0, time.UTC),
					CreatedBy:        "caller-id",
				},
				{
					SessionId:        "session-id",
					SessionName:      "session-name",
					SessionScore:     2,
					SessionStartedAt: time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
					UserId:           "user-id-3",
					UserExternalName: "user-external-name-3",
					UserGeneration:   1.5,
					UserJoinedAt:     time.Date(2024, 1, 1, 19, 59, 59, 0, time.UTC),
					CreatedBy:        "caller-id",
				},
			}).Return(nil)
			mockOpenSessionRepo.EXPECT().MarkAsAttendanceApplied("session-id").Return(nil)
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			assert.NoError(t, err)
		})

		t.Run("Ignore submissions that are after the session starts", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := NewMocksessionRepo(ctrl)
			mockAttendanceFormHandler := NewMockattendanceFormHandler(ctrl)
			mockOpenSessionRepo := NewMockopenSessionRepo(ctrl)
			mockUserRepo := NewMockuserRepo(ctrl)
			mockAttendanceRepo := NewMockattendanceRepo(ctrl)
			server := New(nil, nil, mockUserRepo, nil, nil, mockSessionRepo, mockOpenSessionRepo,
				mockAttendanceFormHandler, mockAttendanceRepo, nil, nil)

			// Do.
			mockSessionRepo.EXPECT().Get("session-id").Return(session.Session{
				Id:               "session-id",
				Name:             "session-name",
				GoogleFormId:     "form-id",
				GoogleFormUri:    "form-uri",
				AttendanceStatus: session.AttendanceStatusNotAppliedYet,
				Score:            2,
				StartsAt:         time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
			}, nil)
			mockAttendanceFormHandler.EXPECT().GetSubmissions("form-id").Return([]attendance.FormSubmission{
				{
					UserExternalName: "user-external-name-1",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 0, 0, 0, time.UTC),
				},
				{
					UserExternalName: "user-external-name-2",
					SubmissionTime:   time.Date(2024, 1, 1, 19, 59, 0, 0, time.UTC),
				},
				{
					UserExternalName: "user-external-name-3",
					SubmissionTime:   time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
				},
			}, nil)
			mockUserRepo.EXPECT().GetAllByExternalNames([]string{"user-external-name-1", "user-external-name-2"}).
				Return([]user.User{
					{
						Id:           "user-id-1",
						ExternalName: "user-external-name-1",
						Generation:   1,
					},
					{
						Id:           "user-id-2",
						ExternalName: "user-external-name-2",
						Generation:   1,
					},
				}, nil)
			mockAttendanceRepo.EXPECT().BulkInsert([]attendance.AddAttendanceReq{
				{
					SessionId:        "session-id",
					SessionName:      "session-name",
					SessionScore:     2,
					SessionStartedAt: time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
					UserId:           "user-id-1",
					UserExternalName: "user-external-name-1",
					UserGeneration:   1,
					UserJoinedAt:     time.Date(2024, 1, 1, 19, 0, 0, 0, time.UTC),
					CreatedBy:        "caller-id",
				},
				{
					SessionId:        "session-id",
					SessionName:      "session-name",
					SessionScore:     2,
					SessionStartedAt: time.Date(2024, 1, 1, 20, 0, 0, 0, time.UTC),
					UserId:           "user-id-2",
					UserExternalName: "user-external-name-2",
					UserGeneration:   1,
					UserJoinedAt:     time.Date(2024, 1, 1, 19, 59, 0, 0, time.UTC),
					CreatedBy:        "caller-id",
				},
			}).Return(nil)
			mockOpenSessionRepo.EXPECT().MarkAsAttendanceApplied("session-id").Return(nil)
			err := server.ApplyAttendanceByFormSubmissions("session-id", "caller-id")

			// Assert.
			assert.NoError(t, err)
		})
	})
}
