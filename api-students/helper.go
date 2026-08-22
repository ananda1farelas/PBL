package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Helper untuk respons sukses 200 OK
func ok(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Helper untuk respons daftar data dengan paginasi (200 OK)
func okList(c *fiber.Ctx, message string, data any, meta *Meta) error {
	return c.Status(fiber.StatusOK).JSON(WebResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Helper untuk respons 201 Created (Wajib menyertakan Header Location)
func created(c *fiber.Ctx, message string, data any, location string) error {
	c.Set("Location", location)
	return c.Status(fiber.StatusCreated).JSON(WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Helper untuk respons 204 No Content (Biasanya untuk DELETE)
func noContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Helper untuk respons gagal umum (400, 404, 409, 415, 500)
func fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(WebResponse{
		Success: false,
		Message: message,
	})
}

// Helper untuk respons gagal validasi (422 Unprocessable Entity)
func failValidation(c *fiber.Ctx, errs map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(WebResponse{
		Success: false,
		Message: "validasi gagal",
		Errors:  errs,
	})
}

// Whitelist kolom yang boleh diurutkan (Mencegah SQL Injection & Error)
var allowedSort = map[string]bool{
	"id":         true,
	"nim":        true,
	"name":       true,
	"grade":      true,
	"created_at": true,
}

// parseListQuery membaca query string dari URL dan memberikan nilai bawaan yang aman
func parseListQuery(c *fiber.Ctx) ListQuery {
	q := ListQuery{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: strings.TrimSpace(c.Query("search")),
		Sort:   c.Query("sort", "id"),
		Order:  strings.ToLower(c.Query("order", "asc")),
	}

	// Validasi halaman minimal 1
	if q.Page < 1 {
		q.Page = 1
	}

	// Validasi limit halaman (Minimal 1, Batas Atas Maksimal 100)
	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}

	// Cek whitelist untuk sort
	if !allowedSort[q.Sort] {
		q.Sort = "id"
	}

	// Validasi order hanya 'asc' atau 'desc'
	if q.Order != "desc" {
		q.Order = "asc"
	}

	// Parse filter is_active jika ada di query parameter
	if raw := c.Query("is_active"); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			q.IsActive = &v
		}
	}

	return q
}
