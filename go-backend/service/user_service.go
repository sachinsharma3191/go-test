// Package service provides business logic layer for the REST API.
//
// This package contains the service layer that sits between the HTTP handlers
// and the repository layer. It implements business rules, validation, and
// orchestration of repository operations.
//
// Service Responsibilities:
//   - Business logic implementation
//   - Input validation and business rule enforcement
//   - Repository operation orchestration
//   - Error handling and business exception management
//
// Design Principles:
//   - Separation of concerns from HTTP and data layers
//   - Dependency injection for repositories
//   - Business rule validation
//   - Transaction management (when applicable)
//
// Service Layer Architecture:
//   - UserService: User business logic and operations
//   - TaskService: Task business logic and user relationships
//   - HealthService: System health monitoring logic
package service

import (
	"strings"

	"go-backend/errors"
	"go-backend/model"
	"go-backend/repository"
	"go-backend/validation"
)

// UserService contains business logic for user-related operations.
// This service encapsulates all user-related business rules and provides
// a clean interface for the HTTP layer to interact with user data.
//
// Business Logic:
//   - User creation with validation
//   - Email uniqueness checking
//   - User data updates with validation
//   - User deletion with dependency checking
//
// Dependencies:
//   - userRepository: Data access layer for user operations
//
// Thread Safety:
//   - Thread safe as long as the repository is thread safe
//   - No mutable state stored in the service
type UserService struct {
	userRepository *repository.UserRepository
}

// NewUserService creates a UserService backed by the given repository.
// This function implements dependency injection, making the service
// easy to test with mock repositories.
//
// Parameters:
//   - repo: Repository for user data access operations
//
// Returns:
//   - *UserService: A new user service instance
//
// Example:
//   service := NewUserService(userRepository)
//   user, err := service.CreateUser("John", "john@example.com", "developer")
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
