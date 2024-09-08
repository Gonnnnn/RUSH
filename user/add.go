package user

import (
	"fmt"
)

// It is separated from userRepo because it requires additional logic and it should be centralized in one place.
type adder struct {
	userRepo UserRepo
}

func NewAdder(userRepo UserRepo) *adder {
	return &adder{
		userRepo: userRepo,
	}
}

func (a *adder) Add(name string, generation float64, isActive bool, email string) error {
	count, err := a.userRepo.CountByName(name)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	err = a.userRepo.Add(User{
		Name:       name,
		Generation: generation,
		IsActive:   isActive,
		Email:      email,
		ExternalName: func() string {
			if count == 0 {
				return name
			}
			return fmt.Sprintf("name%d", count+1)
		}(),
	})

	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}
	return nil
}

//go:generate mockgen -source=add.go -destination=add_mock.go -package=user
type UserRepo interface {
	CountByName(name string) (int, error)
	Add(u User) error
}
