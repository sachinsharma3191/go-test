package json

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go-backend/store"
)

// Store implements a JSON file-based store
type Store struct {
	filePath string
	mu       sync.RWMutex
}

// NewStore creates a new JSON store
func NewStore(filePath string) (*Store, error) {
	store := &Store{
		filePath: filePath,
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	return store, nil
}

// Load loads data from the JSON file
func (s *Store) Load(ctx context.Context) (store.AppData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var data store.AppData

	// Check if file exists
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// Return zero value if file doesn't exist
		return data, nil
	}

	// Read file
	fileData, err := os.ReadFile(s.filePath)
	if err != nil {
		return data, fmt.Errorf("failed to read file: %w", err)
	}

	// Check if file is empty
	if len(fileData) == 0 {
		return data, nil
	}

	// Parse JSON
	if err := json.Unmarshal(fileData, &data); err != nil {
		return data, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return data, nil
}

// Save saves data to the JSON file
func (s *Store) Save(ctx context.Context, data store.AppData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to temporary file first
	tempFile := s.filePath + ".tmp"
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Rename temp file to actual file (atomic operation)
	if err := os.Rename(tempFile, s.filePath); err != nil {
		// Clean up temp file if rename fails
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}

// Health checks if the store is healthy
func (s *Store) Health(ctx context.Context) error {
	// Try to load data to check if file is accessible
	_, err := s.Load(ctx)
	return err
}

// Close cleans up resources
func (s *Store) Close() error {
	// JSON store doesn't need to close anything
	return nil
}
