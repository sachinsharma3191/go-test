package repository

import "go-backend/store"

// Repository is a generic CRUD abstraction for any model type T.
// Implemented by UserRepository and TaskRepository.
type Repository[T any] interface {
	FindAll() ([]T, error)
	FindByID(id int) (*T, error)
	Create(*T) (*T, error)
	Update(*T) (*T, error)
	DeleteByID(id int) error
}

// BaseRepository provides common functionality using the generic store
type BaseRepository struct {
	store store.Store
}

// NewBaseRepository creates a new base repository with the given store
func NewBaseRepository(store store.Store) *BaseRepository {
	return &BaseRepository{store: store}
}
