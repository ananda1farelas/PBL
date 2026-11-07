package middleware

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"latihan-fiber/helper"
)

func Register(app *fiber.App, logger *slog.Logger) {
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(RequestLogger(logger))
}

func RequestLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		requestID, _ := c.Locals("requestid").(string)

		// Status response mungkin belum diisi oleh ErrorHandler
		// ketika c.Next() mengembalikan error.
		status := c.Response().StatusCode()

		if err != nil {
			var appErr *helper.AppError
			if errors.As(err, &appErr) {
				status = appErr.Status
			} else {
				status = fiber.StatusInternalServerError
			}
		}

		attrs := []any{
			slog.String("request_id", requestID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.Duration("duration", time.Since(start)),
			slog.String("ip", c.IP()),
		}

		// Identitas ikut dicatat bila request sudah melewati RequireAuth.
		if user, ok := helper.CurrentUser(c); ok {
			attrs = append(
				attrs,
				slog.Int("user_id", user.UserID),
				slog.String("role", user.Role),
			)
		}

		logger.Info("http_request", attrs...)

		return err
	}
}

// ---------- REQUIRE JSON ----------

var methodsWithBody = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

func RequireJSON(c *fiber.Ctx) error {
	if methodsWithBody[c.Method()] {
		contentType := c.Get("Content-Type")

		if !strings.HasPrefix(
			strings.ToLower(contentType),
			fiber.MIMEApplicationJSON,
		) {
			return helper.UnsupportedMediaType(
				"Content-Type harus application/json",
			)
		}
	}

	return c.Next()
}

func RequireJSONAccept(c *fiber.Ctx) error {
	accept := strings.TrimSpace(c.Get("Accept"))

	// Header Accept kosong berarti klien tidak membatasi format.
	if accept == "" || accept == "*/*" {
		return c.Next()
	}

	// Terima JSON dan wildcard tipe media yang mencakup JSON.
	for _, mediaType := range strings.Split(accept, ",") {
		mediaType = strings.TrimSpace(strings.SplitN(mediaType, ";", 2)[0])

		if mediaType == fiber.MIMEApplicationJSON ||
			mediaType == "application/*" {
			return c.Next()
		}
	}

	return helper.NotAcceptable(
		"format respons yang diminta tidak didukung",
	)
}

// ---------- LOGIN RATE LIMITER ----------

func LoginRateLimiter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
