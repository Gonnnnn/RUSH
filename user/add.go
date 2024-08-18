package user

import (
	"fmt"
)

// It is separated from userRepo because it requires additional logic and it should be centralized in one place.
type adder struct {
	userRepo userRepo
}

func NewAdder(userRepo userRepo) *adder {
	return &adder{
		userRepo: userRepo,
	}
}

func (a *adder) Add(name string, university string, phone string, generation float64, isActive bool) error {
	count, err := a.userRepo.CountByName(name)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	err = a.userRepo.Add(User{
		Name:         name,
		University:   university,
		Phone:        phone,
		Generation:   generation,
		IsActive:     isActive,
		ExternalName: fmt.Sprintf("name%d", count+1),
	})
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}
	return nil
}

type userRepo interface {
	CountByName(name string) (int, error)
	Add(u User) error
}
