package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"api-students/app/model"
)

// Sentinel Error wajib sesuai spesifikasi modul
var (
	ErrNotFound  = errors.New("data mahasiswa tidak ditemukan")
	ErrDuplicate = errors.New("NIM sudah terdaftar di sistem")
)

// Interface StudentRepository dengan 5 method utama
type StudentRepository interface {
	FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error)
	FindByID(ctx context.Context, id int64) (*model.Student, error)
	Create(ctx context.Context, req model.CreateStudentRequest) (*model.Student, error)
	Update(ctx context.Context, id int64, req model.UpdateStudentRequest) (*model.Student, error)
	Delete(ctx context.Context, id int64) error
}

type postgresStudentRepository struct {
	pool *pgxpool.Pool
}

func NewStudentRepository(pool *pgxpool.Pool) StudentRepository {
	return &postgresStudentRepository{pool: pool}
}

func (r *postgresStudentRepository) FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error) {
	var conditions []string
	var args []any
	argIdx := 1

	// Parameterized Query untuk pencarian nama/NIM
	if q.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(nim ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+q.Search+"%")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM students %s", whereClause)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Whitelist pengurutan kolom untuk mencegah SQL Injection
	validColumns := map[string]string{
		"id":         "id",
		"nim":        "nim",
		"name":       "name",
		"grade":      "grade",
		"created_at": "created_at",
	}
	sortCol, ok := validColumns[q.SortBy]
	if !ok {
		sortCol = "id"
	}

	order := "ASC"
	if strings.ToUpper(q.Order) == "DESC" {
		order = "DESC"
	}

	query := fmt.Sprintf(
		"SELECT id, nim, name, grade, is_active, created_at FROM students %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, sortCol, order, argIdx, argIdx+1,
	)
	args = append(args, q.Limit, q.Offset())

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	students := make([]model.Student, 0)
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt); err != nil {
			return nil, 0, err
		}
		students = append(students, s)
	}

	return students, total, rows.Err()
}

func (r *postgresStudentRepository) FindByID(ctx context.Context, id int64) (*model.Student, error) {
	query := `SELECT id, nim, name, grade, is_active, created_at FROM students WHERE id = $1`
	var s model.Student
	err := r.pool.QueryRow(ctx, query, id).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *postgresStudentRepository) Create(ctx context.Context, req model.CreateStudentRequest) (*model.Student, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	query := `
		INSERT INTO students (nim, name, grade, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, nim, name, grade, is_active, created_at`

	var s model.Student
	err := r.pool.QueryRow(ctx, query, req.NIM, req.Name, req.Grade, isActive).
		Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicate
		}
		return nil, err
	}
	return &s, nil
}

func (r *postgresStudentRepository) Update(ctx context.Context, id int64, req model.UpdateStudentRequest) (*model.Student, error) {
	query := `
		UPDATE students
		SET nim = $1, name = $2, grade = $3, is_active = $4
		WHERE id = $5
		RETURNING id, nim, name, grade, is_active, created_at`

	var s model.Student
	err := r.pool.QueryRow(ctx, query, req.NIM, req.Name, req.Grade, req.IsActive, id).
		Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicate
		}
		return nil, err
	}
	return &s, nil
}

func (r *postgresStudentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM students WHERE id = $1`
	ct, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
