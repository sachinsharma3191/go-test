package service

import (
	"strings"

	"go-backend/errors"
	"go-backend/model"
	"go-backend/repository"
	"go-backend/validation"
)

// UserService contains business logic for user-related operations.
type UserService struct {
	userRepository *repository.UserRepository
}

// NewUserService creates a UserService backed by the given repository.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: repo,
	}
}

// ListUsers returns all users.
func (svc *UserService) ListUsers() ([]model.User, error) {
	return svc.userRepository.FindAll()
}

// CreateUser creates a new user after validation and uniqueness checks.
func (svc *UserService) CreateUser(name, email, role string) (*model.User, error) {

	// Syntactic validation
	if err := validation.ValidateUser(name, email, role); err != nil {
		return nil, err
	}

	// Business rule: name and email must be unique (case-insensitive)
	existingUsers, err := svc.ListUsers()
	if err != nil {
		return nil, err
	}

	for _, existingUser := range existingUsers {
		if strings.EqualFold(existingUser.Name, name) {
			return nil, errors.NewDuplicateError("user", "name", nil)
		}
		if strings.EqualFold(existingUser.Email, email) {
			return nil, errors.NewDuplicateError("user", "email", nil)
		}
	}

	newUser := &model.User{
		Name:  name,
		Email: email,
		Role:  role,
	}

	return svc.userRepository.Create(newUser)
}

// GetUserByID returns a user by ID.
func (svc *UserService) GetUserByID(id int) (*model.User, error) {

	user, err := svc.userRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.NewNotFoundError("user", nil)
	}

	return user, nil
}
