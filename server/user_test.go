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
