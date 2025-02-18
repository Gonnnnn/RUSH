package server

import (
	"fmt"
	"rush/permission"
	"rush/user"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestGetAllActiveUsers(t *testing.T) {
	t.Run("Returns internal server error when failed to get users", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepo := NewMockuserRepo(ctrl)
		server := New(nil, nil, mockUserRepo, nil, nil, nil, nil, nil, nil, nil, nil)

		mockUserRepo.EXPECT().GetAll().Return(nil, assert.AnError)
		users, err := server.GetAllActiveUsers()

		assert.Nil(t, users)
		assert.Equal(t, &InternalServerError{originalError: fmt.Errorf("failed to get users: %w", assert.AnError)}, err)
	})

	t.Run("Returns only active users when successfully gets users", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepo := NewMockuserRepo(ctrl)
		server := New(nil, nil, mockUserRepo, nil, nil, nil, nil, nil, nil, nil, nil)

		mockUserRepo.EXPECT().GetAll().Return([]user.User{
			{
				Id:           "user-id",
				Name:         "user-name",
				Role:         permission.RoleMember,
				Generation:   9.5,
				IsActive:     true,
				Email:        "user-email",
				ExternalName: "user-external-name",
			},
			{
				Id:           "user-id2",
				Name:         "user-name2",
				Role:         permission.RoleMember,
				Generation:   9.5,
				IsActive:     false,
				Email:        "user-email2",
				ExternalName: "user-external-name2",
			},
			{
				Id:           "user-id3",
				Name:         "user-name3",
				Role:         permission.RoleMember,
				Generation:   9.5,
				IsActive:     true,
				Email:        "user-email3",
				ExternalName: "user-external-name3",
			},
		}, nil)
		users, err := server.GetAllActiveUsers()

		assert.Equal(t, []*User{
			{
				Id:           "user-id",
				Name:         "user-name",
				Generation:   9.5,
				IsActive:     true,
				Email:        "user-email",
				ExternalName: "user-external-name",
			},
			{
				Id:           "user-id3",
				Name:         "user-name3",
				Generation:   9.5,
				IsActive:     true,
				Email:        "user-email3",
				ExternalName: "user-external-name3",
			},
		}, users)
		assert.NoError(t, err)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("Returns not found error when user is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepo := NewMockuserRepo(ctrl)
		server := New(nil, nil, mockUserRepo, nil, nil, nil, nil, nil, nil, nil, nil)

		mockUserRepo.EXPECT().Get("user-id").Return(nil, user.ErrNotFound)
		dbUser, err := server.GetUser("user-id")

		assert.Nil(t, dbUser)
		assert.Equal(t, &NotFoundError{originalError: fmt.Errorf("failed to get user: %w", user.ErrNotFound)}, err)
	})

	t.Run("Returns internal server error when failed to get user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepo := NewMockuserRepo(ctrl)
		server := New(nil, nil, mockUserRepo, nil, nil, nil, nil, nil, nil, nil, nil)

		mockUserRepo.EXPECT().Get("user-id").Return(nil, assert.AnError)
		dbUser, err := server.GetUser("user-id")

		assert.Nil(t, dbUser)
		assert.Equal(t, &InternalServerError{originalError: fmt.Errorf("failed to get user: %w", assert.AnError)}, err)
	})

	t.Run("Returns user when successfully gets user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserRepo := NewMockuserRepo(ctrl)
		server := New(nil, nil, mockUserRepo, nil, nil, nil, nil, nil, nil, nil, nil)

		mockUserRepo.EXPECT().Get("user-id").Return(&user.User{
			Id:           "user-id",
			Name:         "user-name",
			Role:         permission.RoleMember,
			Generation:   9.5,
			IsActive:     true,
			Email:        "user-email",
			ExternalName: "user-external-name",
		}, nil)
		dbUser, err := server.GetUser("user-id")

		assert.Equal(t, &User{
			Id:           "user-id",
			Name:         "user-name",
			Generation:   9.5,
			IsActive:     true,
			Email:        "user-email",
			ExternalName: "user-external-name",
		}, dbUser)
		assert.NoError(t, err)
	})
}

func TestAddUser(t *testing.T) {
	t.Run("Returns internal server error when failed to add user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserAdder := NewMockuserAdder(ctrl)
		server := New(nil, nil, nil, mockUserAdder, nil, nil, nil, nil, nil, nil, nil)

		mockUserAdder.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)
		err := server.AddUser("user-name", 9.5, true, "user-email")

		assert.Equal(t, &InternalServerError{originalError: fmt.Errorf("failed to add user: %w", assert.AnError)}, err)
	})

	t.Run("Returns nil when successfully adds user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserAdder := NewMockuserAdder(ctrl)
		server := New(nil, nil, nil, mockUserAdder, nil, nil, nil, nil, nil, nil, nil)

		mockUserAdder.EXPECT().Add("user-name", 9.5, true, "user-email").Return(nil)
		err := server.AddUser("user-name", 9.5, true, "user-email")

		assert.NoError(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Returns bad request error when external name is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUpdater := NewMockuserUpdater(ctrl)
		server := New(nil, nil, nil, nil, mockUserUpdater, nil, nil, nil, nil, nil, nil)

		externalName := ""
		err := server.UpdateUser("user-id", &externalName, nil)

		assert.Equal(t, &BadRequestError{originalError: fmt.Errorf("external name is required")}, err)
	})

	t.Run("Returns bad request error when generation is 0", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUpdater := NewMockuserUpdater(ctrl)
		server := New(nil, nil, nil, nil, mockUserUpdater, nil, nil, nil, nil, nil, nil)

		generation := 0.0
		err := server.UpdateUser("user-id", nil, &generation)

		assert.Equal(t, &BadRequestError{originalError: fmt.Errorf("generation is required")}, err)
	})

	t.Run("Returns internal server error when failed to update user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUpdater := NewMockuserUpdater(ctrl)
		server := New(nil, nil, nil, nil, mockUserUpdater, nil, nil, nil, nil, nil, nil)

		externalName := "user-external-name"
		generation := 9.5
		mockUserUpdater.EXPECT().Update("user-id", user.UpdateForm{
			ExternalName: &externalName,
			Generation:   &generation,
		}).Return(assert.AnError)
		err := server.UpdateUser("user-id", &externalName, &generation)

		assert.Equal(t, &InternalServerError{originalError: fmt.Errorf("failed to update user: %w", assert.AnError)}, err)
	})

	t.Run("Returns nil when successfully updates user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserUpdater := NewMockuserUpdater(ctrl)
		server := New(nil, nil, nil, nil, mockUserUpdater, nil, nil, nil, nil, nil, nil)

		externalName := "user-external-name"
		generation := 9.5
		mockUserUpdater.EXPECT().Update("user-id", user.UpdateForm{
			ExternalName: &externalName,
			Generation:   &generation,
		}).Return(nil)

		err := server.UpdateUser("user-id", &externalName, &generation)

		assert.NoError(t, err)
	})
}
