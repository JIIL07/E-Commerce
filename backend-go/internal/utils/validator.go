package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type ValidationError struct {
	Field   string
	Message string
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Message)
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

type Validator struct {
	errors ValidationErrors
}

func NewValidator() *Validator {
	return &Validator{
		errors: make(ValidationErrors, 0),
	}
}

func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

func (v *Validator) Validate() error {
	if len(v.errors) == 0 {
		return nil
	}
	return v.errors
}

func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

func (v *Validator) GetErrors() ValidationErrors {
	return v.errors
}

func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "is required")
	}
	return v
}

func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.AddError(field, fmt.Sprintf("must be at least %d characters long", min))
	}
	return v
}

func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.AddError(field, fmt.Sprintf("must be at most %d characters long", max))
	}
	return v
}

func (v *Validator) Email(field, value string) *Validator {
	if value == "" {
		return v
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.AddError(field, "must be a valid email address")
	}
	return v
}

func (v *Validator) Phone(field, value string) *Validator {
	if value == "" {
		return v
	}

	phoneRegex := regexp.MustCompile(`^\+?[\d\s\-\(\)]{10,}$`)
	if !phoneRegex.MatchString(value) {
		v.AddError(field, "must be a valid phone number")
	}
	return v
}

func (v *Validator) URL(field, value string) *Validator {
	if value == "" {
		return v
	}

	urlRegex := regexp.MustCompile(`^https?:\/\/[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(value) {
		v.AddError(field, "must be a valid URL")
	}
	return v
}

func (v *Validator) Numeric(field, value string) *Validator {
	if value == "" {
		return v
	}

	numericRegex := regexp.MustCompile(`^\d+$`)
	if !numericRegex.MatchString(value) {
		v.AddError(field, "must be numeric")
	}
	return v
}

func (v *Validator) Decimal(field, value string) *Validator {
	if value == "" {
		return v
	}

	decimalRegex := regexp.MustCompile(`^\d+(\.\d+)?$`)
	if !decimalRegex.MatchString(value) {
		v.AddError(field, "must be a valid decimal number")
	}
	return v
}

func (v *Validator) Alpha(field, value string) *Validator {
	if value == "" {
		return v
	}

	alphaRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	if !alphaRegex.MatchString(value) {
		v.AddError(field, "must contain only letters")
	}
	return v
}

func (v *Validator) AlphaNumeric(field, value string) *Validator {
	if value == "" {
		return v
	}

	alphaNumericRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphaNumericRegex.MatchString(value) {
		v.AddError(field, "must contain only letters and numbers")
	}
	return v
}

func (v *Validator) Password(field, value string) *Validator {
	if value == "" {
		return v
	}

	if len(value) < 8 {
		v.AddError(field, "must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range value {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		v.AddError(field, "must contain at least one uppercase letter")
	}
	if !hasLower {
		v.AddError(field, "must contain at least one lowercase letter")
	}
	if !hasDigit {
		v.AddError(field, "must contain at least one digit")
	}
	if !hasSpecial {
		v.AddError(field, "must contain at least one special character")
	}

	return v
}

func (v *Validator) In(field, value string, options []string) *Validator {
	if value == "" {
		return v
	}

	for _, option := range options {
		if value == option {
			return v
		}
	}

	v.AddError(field, fmt.Sprintf("must be one of: %s", strings.Join(options, ", ")))
	return v
}

func (v *Validator) NotIn(field, value string, options []string) *Validator {
	if value == "" {
		return v
	}

	for _, option := range options {
		if value == option {
			v.AddError(field, fmt.Sprintf("must not be one of: %s", strings.Join(options, ", ")))
			return v
		}
	}

	return v
}

func (v *Validator) Range(field string, value, min, max float64) *Validator {
	if value < min || value > max {
		v.AddError(field, fmt.Sprintf("must be between %f and %f", min, max))
	}
	return v
}

func (v *Validator) Custom(field string, value interface{}, fn func(interface{}) error) *Validator {
	if err := fn(value); err != nil {
		v.AddError(field, err.Error())
	}
	return v
}

func ValidateStruct(data interface{}) error {
	validator := NewValidator()

	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if str, ok := value.(string); ok {
				validator.Required(key, str)
			}
		}
	default:
		return errors.New("unsupported data type for validation")
	}

	return validator.Validate()
}

func ValidateEmail(email string) error {
	validator := NewValidator()
	validator.Email("email", email)
	return validator.Validate()
}

func ValidatePassword(password string) error {
	validator := NewValidator()
	validator.Password("password", password)
	return validator.Validate()
}

func ValidatePhone(phone string) error {
	validator := NewValidator()
	validator.Phone("phone", phone)
	return validator.Validate()
}

func ValidateURL(url string) error {
	validator := NewValidator()
	validator.URL("url", url)
	return validator.Validate()
}

func SanitizeString(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")
	input = strings.ReplaceAll(input, "\t", "")
	return input
}

func SanitizeHTML(input string) string {
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#x27;")
	input = strings.ReplaceAll(input, "&", "&amp;")
	return input
}
