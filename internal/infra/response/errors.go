package response

import (
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

func ErrConflict(message ...string) *Err {
	return newError(fiber.StatusConflict, "Conflict", message...)
}

var validationMessages = map[string]string{
	"required": "The {field} field is required.",
	"email":    "The {field} field must be a valid email format.",
	"min":      "The {field} field must be at least {param} characters long.",
	"max":      "The {field} field must be at most {param} characters long.",
	"uuid":     "The {field} field must be a valid UUID format.",
	"numeric":  "The {field} field must be a number.",
}

func ErrValidation(errs validator.ValidationErrors) *Err {
	errorsMap := make(map[string]string)

	for _, err := range errs {
		field := strings.ToLower(err.Field())
		tag := err.Tag()
		param := err.Param()

		msg, exists := validationMessages[tag]
		if !exists {
			errorsMap[field] = err.Error()
			continue
		}

		msg = strings.Replace(msg, "{field}", field, -1)

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
