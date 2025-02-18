package server

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	controller := gomock.NewController(t)
	mockOauthClient := NewMockoauthClient(controller)
	mockAuthHandler := NewMockauthHandler(controller)
	mockUserRepo := NewMockuserRepo(controller)
	mockUserAdder := NewMockuserAdder(controller)
	mockUserUpdater := NewMockuserUpdater(controller)
	mockSessionRepo := NewMocksessionRepo(controller)
	mockOpenSessionRepo := NewMockopenSessionRepo(controller)
	mockAttendanceFormHandler := NewMockattendanceFormHandler(controller)
	mockAttendanceRepo := NewMockattendanceRepo(controller)
	formTimeLocation := time.UTC
	clock := clock.NewMock()

	server := New(mockOauthClient, mockAuthHandler, mockUserRepo, mockUserAdder, mockUserUpdater, mockSessionRepo, mockOpenSessionRepo, mockAttendanceFormHandler, mockAttendanceRepo, formTimeLocation, clock)

	assert.Equal(t, &Server{
		oauthClient:           mockOauthClient,
		authHandler:           mockAuthHandler,
		userRepo:              mockUserRepo,
		userAdder:             mockUserAdder,
		userUpdater:           mockUserUpdater,
		sessionRepo:           mockSessionRepo,
		openSessionRepo:       mockOpenSessionRepo,
		attendanceFormHandler: mockAttendanceFormHandler,
		attendanceRepo:        mockAttendanceRepo,
		formTimeLocation:      formTimeLocation,
		clock:                 clock,
	}, server)
}
