package pkg

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationErrors formats validator.ValidationErrors into a slice of ValidationError
// Uses JSON field names for better API responses
func FormatValidationErrors(err error) []ValidationError {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return []ValidationError{{Message: err.Error()}}
	}

	errArr := make([]ValidationError, 0, len(ve))
	for _, fe := range ve {
		errArr = append(errArr, ValidationError{
			Field:   toJSONFieldName(fe),
			Message: getErrorMessage(fe),
		})
	}
	return errArr
}

// toJSONFieldName converts struct field name to JSON field name
func toJSONFieldName(fe validator.FieldError) string {
	// Get the field name in lowercase snake_case style for API consistency
	field := fe.Field()
	return toSnakeCase(field)
}

// toSnakeCase converts PascalCase/camelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// getErrorMessage returns a human-readable error message for a validation error
func getErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fe.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "dive":
		return fmt.Sprintf("%s contains invalid items", field)
	default:
		return fe.Error()
	}
}

// ConfigureValidator configures Gin's validator to use JSON field names in error messages
func ConfigureValidator(engine *gin.Engine) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register a function to get JSON field names
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return fld.Name
			}
			if name == "" {
				return fld.Name
			}
			return name
		})
	}
}

// HandleValidationError sends a standardized validation error response
func HandleValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		c.JSON(400, gin.H{"errors": FormatValidationErrors(err)})
		return
	}
	c.JSON(400, gin.H{"error": "Invalid request body"})
}
