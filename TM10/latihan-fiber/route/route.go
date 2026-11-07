package route

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"latihan-fiber/app/service"
	"latihan-fiber/helper"
	"latihan-fiber/middleware"
)

type Dependencies struct {
	Pool          *pgxpool.Pool
	Mongo         *mongo.Client
	JWT           *helper.JWTManager
	Permissions   *helper.PermissionSet
	UserService   *service.UserService
	AuthService   service.AuthService
	CourseService *service.CourseService
	ReviewService *service.ReviewService
}

func Register(app *fiber.App, deps Dependencies) {
	api := app.Group("/api/v1")

	// =========================================================
	// PUBLIC
	// =========================================================

	// Health check
	api.Get("/health", healthCheck(deps.Pool, deps.Mongo))

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

	// =========================================================
	// USERS
	// =========================================================

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

	// =========================================================
	// COURSES
	// =========================================================

	courses := api.Group(
		"/courses",
		middleware.RequireJSON,
		middleware.RequireAuth(deps.JWT),
	)

	courses.Get(
		"/",
		deps.CourseService.List,
	)

	// WAJIB sebelum route /:id
	courses.Get(
		"/mine",
		deps.CourseService.Mine,
	)

	courses.Get(
		"/:id",
		deps.CourseService.Get,
	)

	courses.Post(
		"/:id/enroll",
		deps.CourseService.Enroll,
	)

	courses.Delete(
		"/:id/enroll",
		deps.CourseService.Drop,
	)

	courses.Patch(
		"/:id",
		middleware.RequirePermission(
			perms,
			"course:manage",
		),
		deps.CourseService.UpdateQuota,
	)

	courses.Get(
		"/:id/reviews",
		deps.ReviewService.ListByCourse,
	)

	courses.Get(
		"/:id/reviews/summary",
		deps.ReviewService.Summary,
	)

	courses.Post(
		"/:id/reviews",
		deps.ReviewService.Create,
	)

	// =========================================================
	// ENROLLMENTS
	// =========================================================

	enrollments := api.Group(
		"/enrollments",
		middleware.RequireJSON,
		middleware.RequireAuth(deps.JWT),
	)

	enrollments.Post(
		"/swap",
		deps.CourseService.Swap,
	)

	// =========================================================
	// REVIEWS
	// =========================================================
	reviews := api.Group("/reviews",
		middleware.RequireJSON,
		middleware.RequireAuth(deps.JWT))

	reviews.Get(
		"/:id",
		deps.ReviewService.Get,
	)

	reviews.Patch(
		"/:id",
		deps.ReviewService.Update,
	)

	reviews.Delete(
		"/:id",
		deps.ReviewService.Delete,
	)

	reviews.Post(
		"/:id/replies",
		deps.ReviewService.Reply,
	)

	reviews.Post(
		"/:id/helpful",
		deps.ReviewService.MarkHelpful,
	)
}

// healthCheck kini memeriksa DUA database. Aplikasi dianggap sehat hanya
// bila keduanya dapat dihubungi, dan pesan kegagalannya menyebut yang mana.
func healthCheck(pool *pgxpool.Pool, mongoClient *mongo.Client) fiber.Handler {
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
				"PostgreSQL tidak dapat dihubungi",
			)
		}

		if err := mongoClient.Ping(ctx, nil); err != nil {
			return helper.Fail(
				c,
				fiber.StatusServiceUnavailable,
				"MongoDB tidak dapat dihubungi",
			)
		}

		return helper.Success(
			c,
			fiber.StatusOK,
			"server, PostgreSQL, dan MongoDB berjalan",
			nil,
		)
	}
}
