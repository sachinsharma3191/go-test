package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"go-backend/model"
	"go-backend/repository"
	"go-backend/store"

	_ "go-backend/store/json"
)

type userSaveFailingStore struct {
	data store.AppData
}

func (s *userSaveFailingStore) Load(_ context.Context) (store.AppData, error) {
	return s.data, nil
}
func (s *userSaveFailingStore) Save(_ context.Context, _ store.AppData) error {
	return errors.New("save failed")
}
func (s *userSaveFailingStore) Health(_ context.Context) error { return nil }
func (s *userSaveFailingStore) Close() error                  { return nil }

func setupUserService(t *testing.T) *UserService {
	t.Helper()
	f, err := os.CreateTemp("", "user_svc_*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })

	s, err := store.QuickCreate("json", path)
	if err != nil {
		t.Fatalf("QuickCreate: %v", err)
	}
	// Seed a user
	data := store.AppData{
		Users: []model.User{{ID: 1, Name: "Alice", Email: "alice@example.com", Role: "developer"}},
		Tasks: []model.Task{},
	}
	if err := s.Save(context.Background(), data); err != nil {
		t.Fatalf("Save: %v", err)
	}
	repo := repository.NewUserRepository(s)
	return NewUserService(repo)
}

func TestUserService_ListUsers(t *testing.T) {
	svc := setupUserService(t)
	users, err := svc.ListUsers()
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("len(users) = %d, want 1", len(users))
	}
	if users[0].Name != "Alice" {
		t.Errorf("users[0].Name = %q", users[0].Name)
	}
}

func TestUserService_CreateUser(t *testing.T) {
	svc := setupUserService(t)
	u, err := svc.CreateUser("Bob", "bob@example.com", "designer")
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if u.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if u.Name != "Bob" || u.Email != "bob@example.com" || u.Role != "designer" {
		t.Errorf("user = %+v", u)
	}
}

func TestUserService_CreateUser_DuplicateName(t *testing.T) {
	svc := setupUserService(t)
	_, err := svc.CreateUser("Alice", "other@example.com", "designer")
	if err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	svc := setupUserService(t)
	_, err := svc.CreateUser("Other", "alice@example.com", "designer")
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestUserService_CreateUser_ValidationError(t *testing.T) {
	svc := setupUserService(t)
	_, err := svc.CreateUser("", "x@x.com", "developer")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	_, err = svc.CreateUser("x", "invalid-email", "developer")
	if err == nil {
		t.Fatal("expected error for invalid email")
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	svc := setupUserService(t)
	u, err := svc.GetUserByID(1)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if u == nil || u.Name != "Alice" {
		t.Errorf("user = %+v", u)
	}
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	svc := setupUserService(t)
	u, err := svc.GetUserByID(999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
	if u != nil {
		t.Error("expected nil user")
	}
}

func TestUserService_CreateUser_DataStoreError(t *testing.T) {
	st := &userSaveFailingStore{}
	repo := repository.NewUserRepository(st)
	svc := NewUserService(repo)
	_, err := svc.CreateUser("New", "new@example.com", "developer")
	if err == nil {
		t.Fatal("expected error when store Save fails")
	}
}

type loadFailStore struct{}

func (loadFailStore) Load(context.Context) (store.AppData, error) {
	return store.AppData{}, errors.New("load failed")
}
func (loadFailStore) Save(context.Context, store.AppData) error { return nil }
func (loadFailStore) Health(context.Context) error               { return nil }
func (loadFailStore) Close() error                               { return nil }

func TestUserService_GetUserByID_StoreError(t *testing.T) {
	repo := repository.NewUserRepository(loadFailStore{})
	svc := NewUserService(repo)
	_, err := svc.GetUserByID(1)
	if err == nil {
		t.Fatal("expected error when store Load fails")
	}
}
