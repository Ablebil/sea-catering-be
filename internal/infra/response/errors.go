package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func newError(code int, defaultMsg string, message ...string) *Err {
	msg := defaultMsg
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	return &Err{
		Code:    code,
		Message: msg,
	}
}

func ErrInternalServerError(message ...string) *Err {
	return newError(fiber.StatusInternalServerError, "Internal Server Error", message...)
}

func ErrBadRequest(message ...string) *Err {
	return newError(fiber.StatusBadRequest, "Bad Request", message...)
}

func ErrNotFound(message ...string) *Err {
	return newError(fiber.StatusNotFound, "Not Found", message...)
}

func ErrUnauthorized(message ...string) *Err {
	return newError(fiber.StatusUnauthorized, "Unauthorized", message...)
}

func ErrForbidden(message ...string) *Err {
	return newError(fiber.StatusForbidden, "Forbidden", message...)
}

var validationMessages = map[string]string{
	"required": "This field is required.",
	"email":    "Invalid email format.",
	"min":      "This field must be at least %d characters long.",
	"max":      "This field must be at most %d characters long.",
	"uuid":     "Invalid UUID format.",
	"numeric":  "This field must be a number.",
}

func getValidationMessage(tag string, params ...interface{}) string {
	if msg, exists := validationMessages[tag]; exists {
		if len(params) > 0 {
			return fmt.Sprintf(msg, params...)
		}
		return msg
	}
	return "Invalid input"
}

func ErrValidation(errs validator.ValidationErrors) *Err {
	errorsMap := make(map[string]string)
	for _, err := range errs {
		field := strings.ToLower(err.Field())
		tag := err.Tag()
		param := err.Param()

		msgTemplates := getValidationMessage(tag)
		msg := strings.Replace(msgTemplates, "{field}", field, -1)
		msg = strings.Replace(msg, "{param}", param, -1)

		errorsMap[field] = msg
	}

	payload := map[string]interface{}{"errors": errorsMap}
	return &Err{
		Code:    fiber.StatusBadRequest,
		Message: "Validation Error",
		Payload: payload,
	}
}
