package shared

import "fmt"

// ValidationError はバリデーションエラーを表す
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError は新しいValidationErrorを作成
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// DomainError はドメインエラーを表す
type DomainError struct {
	Code    string
	Message string
}

func (e DomainError) Error() string {
	return fmt.Sprintf("domain error [%s]: %s", e.Code, e.Message)
}

// NewDomainError は新しいDomainErrorを作成
func NewDomainError(code, message string) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
	}
}