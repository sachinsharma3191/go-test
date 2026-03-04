// Package repository provides the repository pattern implementation for data access.
//
// This package implements a generic repository pattern that provides a clean abstraction
// layer between the service layer and the data store. It uses Go's generics to provide
// type-safe CRUD operations for any entity type.
//
// Architecture Benefits:
//   - Type Safety: Generic interface ensures compile-time type checking
//   - Testability: Easy to mock for unit testing
//   - Consistency: Standardized CRUD operations across all entities
//   - Extensibility: Easy to add new repository implementations
//
// Repository Pattern:
//   - Repository: Generic interface defining CRUD operations
//   - BaseRepository: Common functionality using generic store
//   - UserRepository: User-specific repository implementation
//   - TaskRepository: Task-specific repository implementation
//
// Usage Example:
//   repo := repository.NewUserRepository(store)
//   users, err := repo.FindAll()
//   user, err := repo.FindByID(1)
package repository

import "go-backend/store"

// Repository is a generic CRUD abstraction for any model type T.
// This interface defines the standard contract for all repository implementations,
// ensuring consistent behavior across different entity types.
//
// Interface Methods:
//   - FindAll(): Retrieves all entities of type T
//   - FindByID(id): Retrieves a single entity by its ID
//   - Create(entity): Creates a new entity and returns it with generated ID
//   - Update(entity): Updates an existing entity by ID
//   - DeleteByID(id): Deletes an entity by its ID
//
// Type Parameters:
//   - T: The entity type (must be a struct type)
//
// Error Handling:
//   - All methods return errors for failure cases
//   - FindByID returns nil when entity not found
//   - Create returns the created entity with assigned ID
//
// Thread Safety:
//   - Thread safety depends on the underlying store implementation
//   - JSON store implementation is thread-safe with proper locking
//
// Implemented by:
//   - UserRepository: Handles User entity operations
//   - TaskRepository: Handles Task entity operations
type Repository[T any] interface {
	// FindAll retrieves all entities of type T from the data store.
	//
	// Returns:
	//   - []T: Slice of all entities
	//   - error: Any error encountered during retrieval
	FindAll() ([]T, error)

	// FindByID retrieves a single entity by its ID.
	//
	// Parameters:
	//   - id: The unique identifier of the entity
	//
	// Returns:
	//   - *T: Pointer to the found entity, or nil if not found
	//   - error: Any error encountered during retrieval
	FindByID(id int) (*T, error)

	// Create creates a new entity in the data store.
	// The entity should not have an ID set; it will be assigned by the repository.
	//
	// Parameters:
	//   - entity: Pointer to the entity to create (without ID)
	//
	// Returns:
	//   - *T: Pointer to the created entity with assigned ID
	//   - error: Any error encountered during creation
	Create(*T) (*T, error)

	// Update updates an existing entity in the data store.
	// The entity must have a valid ID set.
	//
	// Parameters:
	//   - entity: Pointer to the entity to update (with valid ID)
	//
	// Returns:
	//   - *T: Pointer to the updated entity
	//   - error: Any error encountered during update
	Update(*T) (*T, error)

	// DeleteByID removes an entity from the data store by its ID.
	//
	// Parameters:
	//   - id: The unique identifier of the entity to delete
	//
	// Returns:
	//   - error: Any error encountered during deletion
	DeleteByID(id int) error
}

// BaseRepository provides common functionality using the generic store.
// This struct implements the shared functionality that all repositories need,
// such as store access and common error handling patterns.
//
// Design Decisions:
//   - Uses composition instead of inheritance for flexibility
//   - Store is embedded for direct access to underlying storage
//   - Provides a foundation for specific repository implementations
//
// Thread Safety:
//   - Thread safety depends on the injected store implementation
//   - All operations are delegated to the store
//
// Usage:
//   // Create a base repository
//   base := NewBaseRepository(store)
//   // Use base.store in specific repository implementations
type BaseRepository struct {
	store store.Store
}

// NewBaseRepository creates a new base repository with the given store.
// This function provides a factory for creating base repositories with
// proper dependency injection.
//
// Parameters:
//   - store: The data store implementation (JSON, memory, etc.)
//
// Returns:
//   - *BaseRepository: A new base repository instance
//
// Example:
//   store := json.NewStore("data.json")
//   baseRepo := NewBaseRepository(store)
//
// Design Notes:
//   - Store is injected to enable testing with mock implementations
//   - No nil checking is done here - specific repositories should handle nil stores
//   - Store lifecycle management is the responsibility of the caller
func NewBaseRepository(store store.Store) *BaseRepository {
	return &BaseRepository{store: store}
}
