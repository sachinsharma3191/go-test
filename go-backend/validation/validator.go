package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go-backend/errors"
)

// FieldValidator handles validation for specific fields
type FieldValidator struct {
	fieldName string
	value     interface{}
	rules     []ValidationRule
}

// ValidationRule represents a single validation rule
type ValidationRule interface {
	Validate(value interface{}) error
}

// RequiredRule validates that a field is not empty
type RequiredRule struct{}

func (r RequiredRule) Validate(value interface{}) error {
	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return errors.NewValidationError("field is required", nil)
		}
	case *string:
		if v == nil || strings.TrimSpace(*v) == "" {
			return errors.NewValidationError("field is required", nil)
		}
	case int:
		if v == 0 {
			return errors.NewValidationError("field must be greater than 0", nil)
		}
	case *int:
		if v == nil || *v == 0 {
			return errors.NewValidationError("field must be greater than 0", nil)
		}
	default:
		if v == nil {
			return errors.NewValidationError("field is required", nil)
		}
	}
	return nil
}

// MaxLengthRule validates maximum string length
type MaxLengthRule struct {
	MaxLength int
}

func (r MaxLengthRule) Validate(value interface{}) error {
	var str string
	switch v := value.(type) {
	case string:
		str = v
	case *string:
		if v != nil {
			str = *v
		}
	default:
		return nil // Not applicable to non-string types
	}

	if len(str) > r.MaxLength {
		return errors.NewValidationError(
			fmt.Sprintf("field must be %d characters or less", r.MaxLength),
			nil,
		)
	}
	return nil
}

// EmailRule validates email format
type EmailRule struct{}

func (r EmailRule) Validate(value interface{}) error {
	var str string
	switch v := value.(type) {
	case string:
		str = v
	case *string:
		if v != nil {
			str = *v
		}
	default:
		return nil // Not applicable to non-string types
	}

	if str == "" {
		return nil // Skip email validation if empty (use RequiredRule for required emails)
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(str) {
		return errors.NewValidationError("invalid email format", nil)
	}
	return nil
}

// EnumRule validates that value is one of allowed values
type EnumRule struct {
	AllowedValues []string
}

func (r EnumRule) Validate(value interface{}) error {
	var str string
	switch v := value.(type) {
	case string:
		str = v
	case *string:
		if v != nil {
			str = *v
		}
	default:
		return nil // Not applicable to non-string types
	}

	if str == "" {
		return nil // Skip enum validation if empty (use RequiredRule for required fields)
	}

	for _, allowed := range r.AllowedValues {
		if str == allowed {
			return nil
		}
	}

	return errors.NewValidationError(
		fmt.Sprintf("value must be one of: %s", strings.Join(r.AllowedValues, ", ")),
		nil,
	)
}

// Validator manages multiple field validations
type Validator struct {
	fields []*FieldValidator
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		fields: make([]*FieldValidator, 0),
	}
}

// AddField adds a field to validate
func (v *Validator) AddField(fieldName string, value interface{}, rules ...ValidationRule) *Validator {
	v.fields = append(v.fields, &FieldValidator{
		fieldName: fieldName,
		value:     value,
		rules:     rules,
	})
	return v
}

// Validate executes all validations
func (v *Validator) Validate() error {
	var validationErrors []string

	for _, field := range v.fields {
		for _, rule := range field.rules {
			if err := rule.Validate(field.value); err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("%s: %v", field.fieldName, err))
			}
		}
	}

	if len(validationErrors) > 0 {
		return errors.NewValidationError(strings.Join(validationErrors, "; "), nil)
	}

	return nil
}

// ValidateUser validates user creation/update requests
func ValidateUser(name, email, role string) error {
	validator := NewValidator()

	validator.AddField("name", name,
		RequiredRule{},
		MaxLengthRule{MaxLength: 100},
	)

	validator.AddField("email", email,
		RequiredRule{},
		EmailRule{},
	)

	validator.AddField("role", role,
		RequiredRule{},
		EnumRule{AllowedValues: []string{"developer", "designer", "manager", "admin", "tester"}},
	)

	return validator.Validate()
}

// ValidateTask validates task creation requests
func ValidateTask(title, status string, userID int) error {
	validator := NewValidator()

	validator.AddField("title", title,
		RequiredRule{},
		MaxLengthRule{MaxLength: 200},
	)

	validator.AddField("status", status,
		RequiredRule{},
		EnumRule{AllowedValues: []string{"pending", "in-progress", "completed"}},
	)

	validator.AddField("userId", userID,
		RequiredRule{},
	)

	return validator.Validate()
}

// ValidateTaskUpdate validates task update requests (all fields optional)
func ValidateTaskUpdate(title, status *string, userID *int) error {
	validator := NewValidator()

	if title != nil {
		validator.AddField("title", *title,
			MaxLengthRule{MaxLength: 200},
		)
	}

	if status != nil {
		validator.AddField("status", *status,
			EnumRule{AllowedValues: []string{"pending", "in-progress", "completed"}},
		)
	}

	if userID != nil {
		validator.AddField("userId", *userID,
			RequiredRule{},
		)
	}

	return validator.Validate()
}

// ValidateID validates that an ID is a positive integer
func ValidateID(idStr string) (int, error) {
	if idStr == "" {
		return 0, errors.NewValidationError("ID is required", nil)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.NewValidationError("ID must be a valid integer", err)
	}

	if id <= 0 {
		return 0, errors.NewValidationError("ID must be a positive integer", nil)
	}

	return id, nil
}
