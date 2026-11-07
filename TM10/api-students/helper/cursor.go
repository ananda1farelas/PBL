package helper

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
)

var ErrInvalidCursor = errors.New("cursor tidak valid")

// ParseCursorQuery membaca limit dan cursor dari query string.
func ParseCursorQuery(c *fiber.Ctx) (model.CursorQuery, error) {
	limit := c.QueryInt("limit", 5)
	if limit < 1 {
		limit = 5
	}
	if limit > 100 {
		limit = 100
	}

	q := model.CursorQuery{Limit: limit}
	rawCursor := strings.TrimSpace(c.Query("cursor"))
	if rawCursor == "" {
		return q, nil
	}

	cursor, err := DecodeCursor(rawCursor)
	if err != nil {
		return model.CursorQuery{}, err
	}
	q.After = &cursor
	return q, nil
}

// EncodeCursor mengubah posisi terakhir menjadi token Base64 URL-safe.
func EncodeCursor(createdAt time.Time, id int) string {
	value := createdAt.Format("2006-01-02 15:04:05.000000") + "|" + strconv.Itoa(id)
	return base64.RawURLEncoding.EncodeToString([]byte(value))
}

// DecodeCursor membaca token dan memvalidasi kedua komponennya.
func DecodeCursor(encoded string) (model.Cursor, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return model.Cursor{}, ErrInvalidCursor
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) != 2 {
		return model.Cursor{}, ErrInvalidCursor
	}

	if _, err := time.ParseInLocation("2006-01-02 15:04:05.000000", parts[0], time.UTC); err != nil {
		return model.Cursor{}, ErrInvalidCursor
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil || id < 1 {
		return model.Cursor{}, ErrInvalidCursor
	}

	return model.Cursor{CreatedAt: parts[0], ID: id}, nil
}
