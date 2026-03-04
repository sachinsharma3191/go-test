package repository

import (
	"testing"

	"go-backend/model"
	"go-backend/store"
)

// storeSpy embeds store.Store so it automatically satisfies the interface
// without needing to implement its full method set.
type storeSpy struct {
	store.Store
}

func TestNewBaseRepository_WithNilStore(t *testing.T) {
	r := NewBaseRepository(nil)

	if r == nil {
		t.Fatalf("expected non-nil BaseRepository")
	}
	if r.store != nil {
		t.Fatalf("expected store to be nil, got %#v", r.store)
	}
}

func TestNewBaseRepository_AssignsProvidedStore(t *testing.T) {
	spy := &storeSpy{}
	var s store.Store = spy // non-nil interface value

	r := NewBaseRepository(s)

	if r == nil {
		t.Fatalf("expected non-nil BaseRepository")
	}
	if r.store != s {
		t.Fatalf("expected repository.store to be the same instance passed in")
	}
}

func TestNewBaseRepository_ReturnsDifferentInstances(t *testing.T) {
	r1 := NewBaseRepository(nil)
	r2 := NewBaseRepository(nil)

	if r1 == r2 {
		t.Fatalf("expected different instances, got same pointer")
	}
}

func TestUserRepository_NilStore(t *testing.T) {
	repo := NewUserRepository(nil)
	_, err := repo.FindAll()
	if err == nil {
		t.Error("FindAll with nil store should error")
	}
	_, err = repo.FindByID(1)
	if err == nil {
		t.Error("FindByID with nil store should error")
	}
	_, err = repo.Create(&model.User{Name: "x", Email: "x@x.com", Role: "developer"})
	if err == nil {
		t.Error("Create with nil store should error")
	}
	_, err = repo.Update(&model.User{ID: 1, Name: "x", Email: "x@x.com", Role: "developer"})
	if err == nil {
		t.Error("Update with nil store should error")
	}
	err = repo.DeleteByID(1)
	if err == nil {
		t.Error("DeleteByID with nil store should error")
	}
}

func TestTaskRepository_NilStore(t *testing.T) {
	repo := NewTaskRepository(nil)
	_, err := repo.FindAll()
	if err == nil {
		t.Error("FindAll with nil store should error")
	}
	_, err = repo.FindByID(1)
	if err == nil {
		t.Error("FindByID with nil store should error")
	}
	_, err = repo.Create(&model.Task{Title: "x", Status: "pending", UserID: 1})
	if err == nil {
		t.Error("Create with nil store should error")
	}
	_, err = repo.Update(&model.Task{ID: 1, Title: "x", Status: "pending", UserID: 1})
	if err == nil {
		t.Error("Update with nil store should error")
	}
	err = repo.DeleteByID(1)
	if err == nil {
		t.Error("DeleteByID with nil store should error")
	}
}
