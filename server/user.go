package server

import (
	"fmt"
	"rush/golang/array"
)

func (s *Server) GetAllUsers() ([]*User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, newInternalServerError(fmt.Errorf("failed to get users: %w", err))
	}

	converted := []*User{}
	for _, user := range users {
		converted = append(converted, fromUser(&user))
	}
	return converted, nil
}

type ListUsersResult struct {
	Users      []User `json:"users"`
	IsEnd      bool   `json:"is_end"`
	TotalCount int    `json:"total_count"`
}

func (s *Server) ListUsers(offset int, pageSize int, onlyActive bool, all bool) (*ListUsersResult, error) {
	if all {
		if onlyActive {
			users, err := s.userRepo.GetAllActive()
			if err != nil {
				return nil, err
			}

			converted := []User{}
			for _, user := range users {
				converted = append(converted, *fromUser(&user))
			}
			return &ListUsersResult{
				Users:      converted,
				IsEnd:      true,
				TotalCount: len(users),
			}, nil
		}
		users, err := s.userRepo.GetAll()
		if err != nil {
			return nil, err
		}

		converted := []User{}
		for _, user := range users {
			converted = append(converted, *fromUser(&user))
		}
		return &ListUsersResult{
			Users:      converted,
			IsEnd:      true,
			TotalCount: len(users),
		}, nil
	}

	listResult, err := s.userRepo.List(offset, pageSize)
	if err != nil {
		return nil, newInternalServerError(fmt.Errorf("failed to list users: %w", err))
	}

	converted := []User{}
	for _, user := range listResult.Users {
		converted = append(converted, *fromUser(&user))
	}

	if onlyActive {
		return &ListUsersResult{
			Users: array.Filter(converted, func(user User) bool {
				return user.IsActive
			}),
			IsEnd:      listResult.IsEnd,
			TotalCount: listResult.TotalCount,
		}, nil
	}

	return &ListUsersResult{
		Users:      converted,
		IsEnd:      listResult.IsEnd,
		TotalCount: listResult.TotalCount,
	}, nil
}

func (s *Server) GetUser(id string) (*User, error) {
	user, err := s.userRepo.Get(id)
	if err != nil {
		return nil, newNotFoundError(fmt.Errorf("failed to get user: %w", err))
	}
	return fromUser(user), nil
}

func (s *Server) AddUser(name string, generation float64, isActive bool, email string) error {
	if err := s.userAdder.Add(name, generation, isActive, email); err != nil {
		return newInternalServerError(fmt.Errorf("failed to add user: %w", err))
	}
	return nil
}
