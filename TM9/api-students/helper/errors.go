package helper

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

// APIError adalah error aplikasi yang akan diubah menjadi HTTP response
// oleh satu ErrorHandler terpusat.
type APIError struct {
	Status  int
	Code    string
	Message string
	Fields  map[string]string
	Cause   error
}

func (e *APIError) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Message
}

func BadRequest(message string) error {
	return &APIError{Status: fiber.StatusBadRequest, Code: "BAD_REQUEST", Message: message}
}

func NotFound(message string) error {
	return &APIError{Status: fiber.StatusNotFound, Code: "NOT_FOUND", Message: message}
}

func Conflict(message string) error {
	return &APIError{Status: fiber.StatusConflict, Code: "CONFLICT", Message: message}
}

func NotAcceptable(message string) error {
	return &APIError{Status: fiber.StatusNotAcceptable, Code: "NOT_ACCEPTABLE", Message: message}
}

func UnsupportedMediaType(message string) error {
	return &APIError{Status: fiber.StatusUnsupportedMediaType, Code: "UNSUPPORTED_MEDIA_TYPE", Message: message}
}

func ValidationFailed(message string, fields map[string]string) error {
	return &APIError{
		Status: fiber.StatusUnprocessableEntity, Code: "VALIDATION_ERROR",
		Message: message, Fields: fields,
	}
}

// ErrorHandler adalah satu-satunya tempat yang membentuk response kegagalan.
func ErrorHandler(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	code := "INTERNAL_ERROR"
	message := "terjadi kesalahan pada server"
	var fields map[string]string

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		status = apiErr.Status
		code = apiErr.Code
		message = apiErr.Message
		fields = apiErr.Fields
		if apiErr.Cause != nil {
			log.Printf("request_id=%v cause=%v", c.Locals("request_id"), apiErr.Cause)
		}
	} else {
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			status = fiberErr.Code
			switch status {
			case fiber.StatusNotFound:
				code, message = "NOT_FOUND", "endpoint tidak ditemukan"
			case fiber.StatusRequestEntityTooLarge:
				code, message = "PAYLOAD_TOO_LARGE", "ukuran payload terlalu besar"
			case fiber.StatusUnsupportedMediaType:
				code, message = "UNSUPPORTED_MEDIA_TYPE", "jenis media tidak didukung"
			default:
				code, message = "BAD_REQUEST", fiberErr.Message
			}
		} else {
			log.Printf("request_id=%v internal_error=%v", c.Locals("request_id"), err)
		}
	}

	requestID, _ := c.Locals("request_id").(string)
	if requestID == "" {
		requestID = GenerateRandomString(16)
	}

	response := fiber.Map{
		"success":    false,
		"code":       code,
		"message":    message,
		"request_id": requestID,
	}
	if len(fields) > 0 {
		response["fields"] = fields
	}

	return c.Status(status).JSON(response)
}

func Forbidden(message string) error {
	return &APIError{Status: fiber.StatusForbidden, Code: "FORBIDDEN", Message: message}
}

func Unauthorized(message string) error {
	return &APIError{Status: fiber.StatusUnauthorized, Code: "UNAUTHORIZED", Message: message}
}
