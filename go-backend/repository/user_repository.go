package repository

import (
	"context"
	"fmt"
	"go-backend/errors"
	"go-backend/model"
	"go-backend/store"
)

// UserRepository implements persistence for users using the generic store
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository returns a new UserRepository backed by the given store
func NewUserRepository(store store.Store) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(store),
	}
}

// FindAll returns all users
func (r *UserRepository) FindAll() ([]model.User, error) {
	if r.store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	data, err := r.store.Load(context.Background())
	if err != nil {
		return nil, err
	}

	users := make([]model.User, len(data.Users))
	copy(users, data.Users)
	return users, nil
}

// FindByID returns the user with the given ID, or nil if not found
func (r *UserRepository) FindByID(id int) (*model.User, error) {
	if r.store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	data, err := r.store.Load(context.Background())
	if err != nil {
		return nil, err
	}

	return data.GetUserByID(id), nil
}

// Create adds a new user, persists, and returns it
func (r *UserRepository) Create(u *model.User) (*model.User, error) {
	if r.store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	data, err := r.store.Load(context.Background())
	if err != nil {
		return nil, err
	}

	newID := 1
	if len(data.Users) > 0 {
		newID = data.Users[len(data.Users)-1].ID + 1
	}
	newUser := model.User{ID: newID, Name: u.Name, Email: u.Email, Role: u.Role}

	data.AddUser(newUser)
	if err := r.store.Save(context.Background(), data); err != nil {
		return nil, errors.NewDataStoreError("Failed to save user data", err)
	}

	return &newUser, nil
}

// Update updates an existing user by ID
func (r *UserRepository) Update(u *model.User) (*model.User, error) {
	if r.store == nil {
		return nil, fmt.Errorf("store is nil")
	}
	data, err := r.store.Load(context.Background())
	if err != nil {
		return nil, err
	}

	if !data.UpdateUser(u.ID, *u) {
		return nil, errors.NewNotFoundError("user", nil)
	}

	if err := r.store.Save(context.Background(), data); err != nil {
		return nil, errors.NewDataStoreError("Failed to update user data", err)
	}

	return u, nil
}

// DeleteByID removes a user by ID
func (r *UserRepository) DeleteByID(id int) error {
	if r.store == nil {
		return fmt.Errorf("store is nil")
	}
	data, err := r.store.Load(context.Background())
	if err != nil {
		return err
	}

	if !data.DeleteUser(id) {
		return errors.NewNotFoundError("user", nil)
	}

	return r.store.Save(context.Background(), data)
}
