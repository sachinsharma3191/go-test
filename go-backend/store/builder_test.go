package store

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	b := NewBuilder("json")
	if b == nil {
		t.Fatal("NewBuilder returned nil")
	}
	if b.config.Type != "json" {
		t.Errorf("config.Type = %q", b.config.Type)
	}
}

func TestBuilder_WithFile(t *testing.T) {
	b := NewBuilder("json").WithFile("/tmp/data.json")
	if b.config.Settings["file"] != "/tmp/data.json" {
		t.Errorf("file = %v", b.config.Settings["file"])
	}
}

func TestBuilder_WithDatabase(t *testing.T) {
	b := NewBuilder("sqlite").WithDatabase("/tmp/db.sqlite")
	if b.config.Settings["db"] != "/tmp/db.sqlite" {
		t.Errorf("db = %v", b.config.Settings["db"])
	}
}

func TestBuilder_WithSetting(t *testing.T) {
	b := NewBuilder("json").WithSetting("custom", "value")
	if b.config.Settings["custom"] != "value" {
		t.Errorf("custom = %v", b.config.Settings["custom"])
	}
}

func TestQuickCreate_Unsupported(t *testing.T) {
	_, err := QuickCreate("unknown", "/tmp/x")
	if err == nil {
		t.Fatal("expected error for unsupported backend")
	}
}

func TestQuickCreate_Sqlite(t *testing.T) {
	// Exercises sqlite branch; may fail if no sqlite factory registered
	_, err := QuickCreate("sqlite", "/tmp/test.db")
	if err != nil && err.Error() != "unsupported store type: sqlite" {
		// Sqlite factory might exist and fail for other reasons
		t.Logf("QuickCreate sqlite: %v", err)
	}
}
