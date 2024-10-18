package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator
func NewValidator() echo.Validator {
	return &CustomValidator{validator: validator.New()}
}

// カスタムヴァリデータを編集
func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
			fieldName := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", fieldName))
			case "email":
				errorMessages = append(errorMessages, fmt.Sprintf("%s isn't email format.", fieldName))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s characters long.", fieldName, err.Param()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("%s is fail validation", fieldName))
			}
		}
		return fmt.Errorf(strings.Join(errorMessages, ", "))
	}
	return nil
}
