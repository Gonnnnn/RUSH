package server

import (
	"fmt"
	"rush/golang/array"
	"rush/user"
)

func (s *Server) GetAllActiveUsers() ([]*User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, newInternalServerError(fmt.Errorf("failed to get users: %w", err))
	}

	activeUsers := array.Filter(users, func(user user.User) bool {
		return user.IsActive
	})

	converted := []*User{}
	for _, user := range activeUsers {
		converted = append(converted, fromUser(&user))
	}
	return converted, nil
}

// Returns the user by the given ID.
func (s *Server) GetUser(id string) (*User, error) {
	user, err := s.userRepo.Get(id)
	if err != nil {
		return nil, newNotFoundError(fmt.Errorf("failed to get user: %w", err))
	}
	return fromUser(user), nil
}

// Adds a new user.
func (s *Server) AddUser(name string, generation float64, isActive bool, email string) error {
	if err := s.userAdder.Add(name, generation, isActive, email); err != nil {
		return newInternalServerError(fmt.Errorf("failed to add user: %w", err))
	}
	return nil
}

// Updates the user.
func (s *Server) UpdateUser(id string, externalName *string, generation *float64) error {
	if externalName != nil && *externalName == "" {
		return newBadRequestError(fmt.Errorf("external name is required"))
	}
	if generation != nil && *generation == 0 {
		return newBadRequestError(fmt.Errorf("generation is required"))
	}

	if err := s.userUpdater.Update(id, user.UpdateForm{
		ExternalName: externalName,
		Generation:   generation,
	}); err != nil {
		return newInternalServerError(fmt.Errorf("failed to update user: %w", err))
	}
	return nil
}
