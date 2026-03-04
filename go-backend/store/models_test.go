package store

import (
	"testing"

	"go-backend/model"
)

func TestNewAppData(t *testing.T) {
	d := NewAppData()
	if d.Users == nil || len(d.Users) != 0 {
		t.Error("Users should be empty slice")
	}
	if d.Tasks == nil || len(d.Tasks) != 0 {
		t.Error("Tasks should be empty slice")
	}
}

func TestAppData_Clone(t *testing.T) {
	d := AppData{
		Users: []model.User{{ID: 1, Name: "A"}},
		Tasks: []model.Task{{ID: 1, Title: "T", Status: "pending", UserID: 1}},
	}
	c := d.Clone()
	if &c.Users[0] == &d.Users[0] {
		t.Error("Clone should copy users")
	}
	if &c.Tasks[0] == &d.Tasks[0] {
		t.Error("Clone should copy tasks")
	}
	if c.Users[0].Name != d.Users[0].Name {
		t.Error("Clone should preserve user data")
	}
}

func TestAppData_IsEmpty(t *testing.T) {
	if !NewAppData().IsEmpty() {
		t.Error("NewAppData should be empty")
	}
	d := AppData{Users: []model.User{{ID: 1}}}
	if d.IsEmpty() {
		t.Error("data with users should not be empty")
	}
}

func TestAppData_GetUserByID(t *testing.T) {
	d := AppData{Users: []model.User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}}
	u := d.GetUserByID(1)
	if u == nil || u.Name != "Alice" {
		t.Errorf("GetUserByID(1) = %v", u)
	}
	if d.GetUserByID(999) != nil {
		t.Error("GetUserByID(999) should be nil")
	}
}

func TestAppData_GetTaskByID(t *testing.T) {
	d := AppData{Tasks: []model.Task{{ID: 1, Title: "T1"}, {ID: 2, Title: "T2"}}}
	task := d.GetTaskByID(1)
	if task == nil || task.Title != "T1" {
		t.Errorf("GetTaskByID(1) = %v", task)
	}
	if d.GetTaskByID(999) != nil {
		t.Error("GetTaskByID(999) should be nil")
	}
}

func TestAppData_AddUser(t *testing.T) {
	var d AppData
	d.AddUser(model.User{ID: 1, Name: "A"})
	if len(d.Users) != 1 {
		t.Errorf("len(Users) = %d", len(d.Users))
	}
}

func TestAppData_AddTask(t *testing.T) {
	var d AppData
	d.AddTask(model.Task{ID: 1, Title: "T", Status: "pending", UserID: 1})
	if len(d.Tasks) != 1 {
		t.Errorf("len(Tasks) = %d", len(d.Tasks))
	}
}

func TestAppData_UpdateUser(t *testing.T) {
	d := AppData{Users: []model.User{{ID: 1, Name: "Old"}}}
	if !d.UpdateUser(1, model.User{ID: 1, Name: "New", Email: "e@e.com", Role: "dev"}) {
		t.Fatal("UpdateUser should succeed")
	}
	if d.Users[0].Name != "New" {
		t.Errorf("UpdateUser did not apply: %s", d.Users[0].Name)
	}
	if d.UpdateUser(999, model.User{ID: 999}) {
		t.Error("UpdateUser(999) should fail")
	}
}

func TestAppData_UpdateTask(t *testing.T) {
	d := AppData{Tasks: []model.Task{{ID: 1, Title: "Old", Status: "pending", UserID: 1}}}
	if !d.UpdateTask(1, model.Task{ID: 1, Title: "New", Status: "completed", UserID: 1}) {
		t.Fatal("UpdateTask should succeed")
	}
	if d.Tasks[0].Title != "New" {
		t.Errorf("UpdateTask did not apply: %s", d.Tasks[0].Title)
	}
	if d.UpdateTask(999, model.Task{ID: 999}) {
		t.Error("UpdateTask(999) should fail")
	}
}

func TestAppData_DeleteUser(t *testing.T) {
	d := AppData{Users: []model.User{{ID: 1}, {ID: 2}}}
	if !d.DeleteUser(1) {
		t.Fatal("DeleteUser should succeed")
	}
	if len(d.Users) != 1 || d.Users[0].ID != 2 {
		t.Errorf("DeleteUser failed: %v", d.Users)
	}
	if d.DeleteUser(999) {
		t.Error("DeleteUser(999) should fail")
	}
}

func TestAppData_DeleteTask(t *testing.T) {
	d := AppData{Tasks: []model.Task{{ID: 1}, {ID: 2}}}
	if !d.DeleteTask(1) {
		t.Fatal("DeleteTask should succeed")
	}
	if len(d.Tasks) != 1 || d.Tasks[0].ID != 2 {
		t.Errorf("DeleteTask failed: %v", d.Tasks)
	}
	if d.DeleteTask(999) {
		t.Error("DeleteTask(999) should fail")
	}
}
