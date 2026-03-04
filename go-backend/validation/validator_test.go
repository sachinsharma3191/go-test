package validation

import (
	"strings"
	"testing"

	apperrors "go-backend/errors"
)

func TestValidateUser_Success(t *testing.T) {
	if err := ValidateUser("Alice", "alice@example.com", "developer"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateUser_Errors(t *testing.T) {
	if err := ValidateUser("", "bad-email", "unknown"); err == nil {
		t.Fatalf("expected validation error, got nil")
	} else if !apperrors.IsErrorCode(err, apperrors.ErrCodeValidation) {
		t.Fatalf("expected validation error code, got %v", err)
	}
}

func TestValidateTask_InvalidStatusOrUser(t *testing.T) {
	if err := ValidateTask("", "pending", 1); err == nil {
		t.Errorf("expected error for empty title")
	}
	if err := ValidateTask("Title", "not-a-status", 1); err == nil {
		t.Errorf("expected error for invalid status")
	}
}

func TestValidateTaskUpdate_StatusEnumAndOptional(t *testing.T) {
	// nil fields should be OK
	if err := ValidateTaskUpdate(nil, nil, nil); err != nil {
		t.Fatalf("expected no error for empty update, got %v", err)
	}

	// invalid status should fail
	title := "Updated"
	status := "bogus"
	if err := ValidateTaskUpdate(&title, &status, nil); err == nil {
		t.Fatalf("expected validation error for bad status")
	}
}

func TestValidateTaskUpdate_OptionalBranches(t *testing.T) {
	// title only (status and userID nil)
	title := "Valid"
	if err := ValidateTaskUpdate(&title, nil, nil); err != nil {
		t.Errorf("title only: %v", err)
	}
	// title too long
	longTitle := string(make([]byte, 201))
	if err := ValidateTaskUpdate(&longTitle, nil, nil); err == nil {
		t.Error("expected error for title > 200 chars")
	}

	// status only (title and userID nil)
	status := "pending"
	if err := ValidateTaskUpdate(nil, &status, nil); err != nil {
		t.Errorf("status only: %v", err)
	}

	// userID only (title and status nil) - RequiredRule for userID expects > 0
	uid := 1
	if err := ValidateTaskUpdate(nil, nil, &uid); err != nil {
		t.Errorf("userID only valid: %v", err)
	}
	zeroUID := 0
	if err := ValidateTaskUpdate(nil, nil, &zeroUID); err == nil {
		t.Error("expected error for userID 0")
	}
}

func TestValidateID(t *testing.T) {
	if _, err := ValidateID(""); err == nil {
		t.Errorf("expected error for empty id")
	}
	if _, err := ValidateID("abc"); err == nil {
		t.Errorf("expected error for non-integer id")
	}
	if _, err := ValidateID("0"); err == nil {
		t.Errorf("expected error for non-positive id")
	}
	id, err := ValidateID("42")
	if err != nil {
		t.Fatalf("expected no error for valid id, got %v", err)
	}
	if id != 42 {
		t.Errorf("expected id 42, got %d", id)
	}
}

func TestValidator_CollectsMultipleErrors(t *testing.T) {
	v := NewValidator()
	v.AddField("name", "", RequiredRule{})
	v.AddField("email", "bad", EmailRule{})
	err := v.Validate()
	if err == nil {
		t.Fatalf("expected aggregated validation error")
	}
	if !apperrors.IsErrorCode(err, apperrors.ErrCodeValidation) {
		t.Fatalf("expected validation error code, got %v", err)
	}
	if !strings.Contains(err.Error(), "name") || !strings.Contains(err.Error(), "email") {
		t.Errorf("expected error message to mention both fields, got %v", err)
	}
}

func TestRequiredRule_PointerTypes(t *testing.T) {
	var s *string
	var i *int
	v := NewValidator()
	v.AddField("s", s, RequiredRule{})
	v.AddField("i", i, RequiredRule{})
	if err := v.Validate(); err == nil {
		t.Error("expected error for nil *string and *int")
	}
}

func TestMaxLengthRule_PointerString(t *testing.T) {
	s := "hello"
	v := NewValidator()
	v.AddField("s", &s, MaxLengthRule{MaxLength: 3})
	if err := v.Validate(); err == nil {
		t.Error("expected error for string exceeding max length")
	}
}

func TestEmailRule_EmptySkipsValidation(t *testing.T) {
	v := NewValidator()
	v.AddField("email", "", EmailRule{})
	if err := v.Validate(); err != nil {
		t.Errorf("empty string should skip email rule: %v", err)
	}
}

func TestEnumRule_EmptySkipsValidation(t *testing.T) {
	v := NewValidator()
	v.AddField("role", "", EnumRule{AllowedValues: []string{"a", "b"}})
	if err := v.Validate(); err != nil {
		t.Errorf("empty string should skip enum rule: %v", err)
	}
}

func TestEnumRule_ValidValue(t *testing.T) {
	v := NewValidator()
	v.AddField("role", "a", EnumRule{AllowedValues: []string{"a", "b"}})
	if err := v.Validate(); err != nil {
		t.Errorf("valid enum value should pass: %v", err)
	}
}

func TestValidateTask_InvalidUserID(t *testing.T) {
	if err := ValidateTask("Title", "pending", 0); err == nil {
		t.Error("expected error for userID 0")
	}
}

func TestValidateTaskUpdate_UserIDRequired(t *testing.T) {
	uid := 0
	if err := ValidateTaskUpdate(nil, nil, &uid); err == nil {
		t.Error("expected error for userID 0")
	}
}

func TestRequiredRule_DefaultCase(t *testing.T) {
	v := NewValidator()
	v.AddField("x", nil, RequiredRule{})
	if err := v.Validate(); err == nil {
		t.Error("expected error for nil in default case")
	}
}

func TestMaxLengthRule_DefaultCase(t *testing.T) {
	v := NewValidator()
	v.AddField("x", 123, MaxLengthRule{MaxLength: 10})
	if err := v.Validate(); err != nil {
		t.Errorf("MaxLengthRule should skip non-string types: %v", err)
	}
}

func TestMaxLengthRule_NilPointer(t *testing.T) {
	var s *string
	v := NewValidator()
	v.AddField("s", s, MaxLengthRule{MaxLength: 5})
	if err := v.Validate(); err != nil {
		t.Errorf("MaxLengthRule should skip nil *string: %v", err)
	}
}

func TestEmailRule_DefaultCase(t *testing.T) {
	v := NewValidator()
	v.AddField("x", 42, EmailRule{})
	if err := v.Validate(); err != nil {
		t.Errorf("EmailRule should skip non-string types: %v", err)
	}
}

func TestEmailRule_NilPointer(t *testing.T) {
	var s *string
	v := NewValidator()
	v.AddField("email", s, EmailRule{})
	if err := v.Validate(); err != nil {
		t.Errorf("EmailRule should skip nil *string: %v", err)
	}
}

func TestEnumRule_DefaultCase(t *testing.T) {
	v := NewValidator()
	v.AddField("x", true, EnumRule{AllowedValues: []string{"a", "b"}})
	if err := v.Validate(); err != nil {
		t.Errorf("EnumRule should skip non-string types: %v", err)
	}
}

func TestEnumRule_NilPointer(t *testing.T) {
	var s *string
	v := NewValidator()
	v.AddField("role", s, EnumRule{AllowedValues: []string{"a", "b"}})
	if err := v.Validate(); err != nil {
		t.Errorf("EnumRule should skip nil *string: %v", err)
	}
}

func TestRequiredRule_IntZero(t *testing.T) {
	v := NewValidator()
	v.AddField("id", 0, RequiredRule{})
	if err := v.Validate(); err == nil {
		t.Error("expected error for int 0")
	}
}

func TestRequiredRule_PointerIntZero(t *testing.T) {
	zero := 0
	v := NewValidator()
	v.AddField("id", &zero, RequiredRule{})
	if err := v.Validate(); err == nil {
		t.Error("expected error for *int 0")
	}
}
