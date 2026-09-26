package route

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"latihan-fiber/app/service"
	"latihan-fiber/helper"
	"latihan-fiber/middleware"
)

type Dependencies struct {
	Pool        *pgxpool.Pool
	JWT         *helper.JWTManager
	Permissions *helper.PermissionSet
	UserService *service.UserService
	AuthService service.AuthService
}

func Register(app *fiber.App, deps Dependencies) {
	api := app.Group("/api/v1")

	// =========================================================
	// PUBLIC
	// =========================================================

	// Health check
	api.Get("/health", healthCheck(deps.Pool))

	// =========================================================
	// AUTHENTICATION
	// =========================================================

	auth := api.Group(
		"/auth",
		middleware.RequireJSON,
	)

	// Register
	auth.Post(
		"/register",
		deps.AuthService.Register,
	)
	auth.Post(
		"/login",
		middleware.LoginRateLimiter(),
		deps.AuthService.Login,
	)
	auth.Post(
		"/refresh",
		deps.AuthService.Refresh,
	)
	auth.Post(
		"/logout",
		deps.AuthService.Logout,
	)
	auth.Get(
		"/me",
		middleware.RequireAuth(deps.JWT),
		deps.AuthService.Me,
	)

	users := api.Group(
		"/users",
		middleware.RequireJSON,
		middleware.RequireAuth(deps.JWT),
	)

	perms := deps.Permissions

	users.Get(
		"/",
		middleware.RequirePermission(
			perms,
			"user:list",
		),
		deps.UserService.List,
	)

	users.Post(
		"/",
		middleware.RequirePermission(
			perms,
			"user:update:any",
		),
		deps.UserService.Create,
	)

	users.Delete(
		"/:id",
		middleware.RequirePermission(
			perms,
			"user:delete",
		),
		deps.UserService.Delete,
	)

	users.Patch(
		"/:id/role",
		middleware.RequirePermission(
			perms,
			"role:assign",
		),
		deps.UserService.AssignRole,
	)

	users.Get(
		"/:id",
		deps.UserService.Get,
	)

	// PUT /api/v1/users/:id
	users.Put(
		"/:id",
		deps.UserService.Replace,
	)

	// PATCH /api/v1/users/:id
	users.Patch(
		"/:id",
		deps.UserService.Patch,
	)
}

func healthCheck(pool *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(
			c.UserContext(),
			2*time.Second,
		)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			return helper.Fail(
				c,
				fiber.StatusServiceUnavailable,
				"database tidak dapat dihubungi",
			)
		}

		return helper.Success(
			c,
			fiber.StatusOK,
			"server dan database berjalan",
			nil,
		)
	}
}
