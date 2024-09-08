package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestAdd(t *testing.T) {
	t.Run("Fails when the repo fails to count the docs with the same name", func(t *testing.T) {
		controller := gomock.NewController(t)
		repo := NewMockUserRepo(controller)

		repo.EXPECT().CountByName("name").Return(0, assert.AnError)
		adder := NewAdder(repo)

		err := adder.Add("name", 1, true, "email")
		assert.Error(t, err)
	})

	t.Run("Fails when the repo fails to add the user", func(t *testing.T) {
		controller := gomock.NewController(t)
		repo := NewMockUserRepo(controller)

		repo.EXPECT().CountByName("name").Return(0, nil)
		repo.EXPECT().Add(User{
			Name:         "name",
			Role:         RoleMember,
			Generation:   1,
			IsActive:     true,
			Email:        "email",
			ExternalName: "name",
		}).Return(assert.AnError)
		adder := NewAdder(repo)

		err := adder.Add("name", 1, true, "email")
		assert.Error(t, err)
	})

	t.Run("Succeeds", func(t *testing.T) {
		controller := gomock.NewController(t)
		repo := NewMockUserRepo(controller)

		// When there was no duplicate name.
		repo.EXPECT().CountByName("name").Return(0, nil)
		repo.EXPECT().Add(User{
			Name:         "name",
			Role:         RoleMember,
			Generation:   1,
			IsActive:     true,
			Email:        "email",
			ExternalName: "name",
		}).Return(nil)
		adder := NewAdder(repo)

		err := adder.Add("name", 1, true, "email")
		assert.NoError(t, err)

		// When there was a duplicate name.
		repo.EXPECT().CountByName("name").Return(1, nil)
		repo.EXPECT().Add(User{
			Name:         "name",
			Role:         RoleMember,
			Generation:   1,
			IsActive:     true,
			Email:        "email",
			ExternalName: "name2",
		}).Return(nil)

		err = adder.Add("name", 1, true, "email")
		assert.NoError(t, err)

		// When there was 2 duplicate names.
		repo.EXPECT().CountByName("name").Return(2, nil)
		repo.EXPECT().Add(User{
			Name:         "name",
			Role:         RoleMember,
			Generation:   1,
			IsActive:     true,
			Email:        "email",
			ExternalName: "name3",
		}).Return(nil)

		err = adder.Add("name", 1, true, "email")
		assert.NoError(t, err)
	})
}
