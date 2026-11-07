package route

import (
	"github.com/gofiber/fiber/v2"

	"api-students/app/handler"
	"api-students/app/service"
	"api-students/helper"
	"api-students/middleware"
)

func RegisterRoutes(app *fiber.App, studentService *service.StudentService, supervisionService *service.SupervisionService, portfolioService *service.PortfolioService, authHandler *handler.AuthHandler, jwtManager *helper.JWTManager) {
	api := app.Group("/api", middleware.RequireJSON)

	students := api.Group("/students")
	students.Get("/", studentService.List)
	students.Get("/:id", studentService.GetByID)
	students.Post("/", studentService.Create)
	students.Put("/:id", studentService.Replace)
	students.Patch("/:id", studentService.UpdatePartial)
	students.Delete("/:id", studentService.Delete)

	// --- portofolio mahasiswa (MongoDB) ---
	students.Get("/:id/portfolio", portfolioService.Get)
	students.Put("/:id/portfolio", portfolioService.Put)
	students.Post("/:id/portfolio/projects", portfolioService.AddProject)
	students.Delete("/:id/portfolio/projects/:projectId", portfolioService.RemoveProject)

	portfolios := api.Group("/portfolios")
	portfolios.Get("/stats/technologies", portfolioService.TechnologyStats)

	// Permission supervision:manage diberikan kepada role admin.
	permissions := helper.NewPermissionSet(map[string][]string{
		"admin": {"supervision:manage"},
	})
	supervisors := api.Group("/supervisors", middleware.JWTAuth(jwtManager), middleware.RequirePermission(permissions, "supervision:manage"))
	supervisors.Post("/:id/assign", supervisionService.Assign)
	supervisors.Patch("/:id", supervisionService.UpdateSupervisor)

	supervisions := api.Group("/supervisions", middleware.JWTAuth(jwtManager), middleware.RequirePermission(permissions, "supervision:manage"))
	supervisions.Post("/switch", supervisionService.Switch)

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/me", middleware.JWTAuth(jwtManager), authHandler.Me)
}
