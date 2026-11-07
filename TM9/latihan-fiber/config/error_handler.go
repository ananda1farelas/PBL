package config

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"latihan-fiber/app/model"
	"latihan-fiber/helper"
)

func newErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var appErr *helper.AppError

		// Jika sudah berupa AppError, gunakan langsung.
		if !errors.As(err, &appErr) {
			switch {
			case errors.Is(err, fiber.ErrRequestEntityTooLarge):
				appErr = &helper.AppError{
					Status:  fiber.StatusRequestEntityTooLarge,
					Code:    "PAYLOAD_TOO_LARGE",
					Message: "ukuran body melebihi batas yang diizinkan",
				}

			default:
				var fiberErr *fiber.Error

				if errors.As(err, &fiberErr) {
					appErr = &helper.AppError{
						Status:  fiberErr.Code,
						Code:    "HTTP_ERROR",
						Message: fiberErr.Message,
					}
				} else {
					appErr = helper.Internal(err)
				}
			}
		}

		requestID := helper.RequestID(c)

		// Error 4xx dicatat sebagai peringatan.
		// Error 5xx dicatat sebagai kesalahan server.
		if logger != nil {
			if appErr.Status >= fiber.StatusInternalServerError {
				logger.Error(
					"request_failed",
					slog.String("request_id", requestID),
					slog.String("path", c.Path()),
					slog.String("code", appErr.Code),
					slog.Int("status", appErr.Status),
					slog.String("error", err.Error()),
				)
			} else {
				logger.Warn(
					"request_rejected",
					slog.String("request_id", requestID),
					slog.String("path", c.Path()),
					slog.String("code", appErr.Code),
					slog.Int("status", appErr.Status),
				)
			}
		}

		return c.Status(appErr.Status).JSON(model.ErrorResponse{
			Success:   false,
			Code:      appErr.Code,
			Message:   appErr.Message,
			Fields:    appErr.Fields,
			RequestID: requestID,
		})
	}
}
