package route

import (
	"github.com/gofiber/fiber/v2"

	"api-students/app/handler"
	"api-students/app/service"
	"api-students/helper"
	"api-students/middleware"
)

func RegisterRoutes(app *fiber.App, studentService *service.StudentService, authHandler *handler.AuthHandler, jwtManager *helper.JWTManager) {
	api := app.Group("/api", middleware.RequireJSON)

	students := api.Group("/students")
	students.Get("/", studentService.List)
	students.Get("/:id", studentService.GetByID)
	students.Post("/", studentService.Create)
	students.Put("/:id", studentService.Replace)
	students.Patch("/:id", studentService.UpdatePartial)
	students.Delete("/:id", studentService.Delete)

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/me", middleware.JWTAuth(jwtManager), authHandler.Me)
}