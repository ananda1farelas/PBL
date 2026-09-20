package route

import (
	"github.com/gofiber/fiber/v2"

	"api-students/app/service"
	"api-students/middleware"
)

func RegisterRoutes(app *fiber.App, studentService *service.StudentService) {
	api := app.Group("/api", middleware.RequireJSON)

	students := api.Group("/students")
	students.Get("/", studentService.List)
	students.Get("/:id", studentService.GetByID)
	students.Post("/", studentService.Create)
	students.Put("/:id", studentService.Replace)
	students.Patch("/:id", studentService.UpdatePartial)
	students.Delete("/:id", studentService.Delete)
}
