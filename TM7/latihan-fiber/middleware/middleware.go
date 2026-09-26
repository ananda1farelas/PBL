package middleware

import (
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

		attrs := []any{
			slog.String("request_id", requestID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration", time.Since(start)),
			slog.String("ip", c.IP()),
		}

		// Identitas ikut dicatat bila request sudah melewati RequireAuth.
		// Tanpa ini, log sebuah 403 tidak berguna: kita tahu ada yang ditolak,
		// tetapi tidak tahu siapa dan mengapa.
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

		if !strings.HasPrefix(contentType, fiber.MIMEApplicationJSON) {
			return helper.Fail(
				c,
				fiber.StatusUnsupportedMediaType,
				"Content-Type harus application/json",
			)
		}
	}

	return c.Next()
}

// ---------- LOGIN RATE LIMITER ----------

func LoginRateLimiter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
