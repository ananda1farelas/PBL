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

var (
	ErrNotFound  = errors.New("data not found")
	ErrDuplicate = errors.New("duplicate data")
)

type StudentRepository interface {
	FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error)
	FindByID(ctx context.Context, id int) (model.Student, error)
	Create(ctx context.Context, student model.Student) (model.Student, error)
	Update(ctx context.Context, student model.Student) (model.Student, error)
	Delete(ctx context.Context, id int) error
}

type studentRepository struct {
	pool *pgxpool.Pool
}

func NewStudentRepository(pool *pgxpool.Pool) StudentRepository {
	return &studentRepository{pool: pool}
}

func (r *studentRepository) FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error) {
	var conditions []string
	var args []any
	argID := 1

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM students %s", whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (q.Page - 1) * q.Limit
	query := fmt.Sprintf(
		"SELECT id, nim, name, grade, is_active FROM students %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, q.Sort, q.Order, argID, argID+1,
	)
	args = append(args, q.Limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive); err != nil {
			return nil, 0, err
		}
		students = append(students, s)
	}

	return students, total, nil
}

func (r *studentRepository) FindByID(ctx context.Context, id int) (model.Student, error) {
	var s model.Student
	query := "SELECT id, nim, name, grade, is_active, created_at FROM students WHERE id = $1"
	err := r.pool.QueryRow(ctx, query, id).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return s, ErrNotFound
	}
	return s, err
}

func (r *studentRepository) Create(ctx context.Context, s model.Student) (model.Student, error) {
	query := "INSERT INTO students (nim, name, grade, is_active, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id"
	err := r.pool.QueryRow(ctx, query, s.NIM, s.Name, s.Grade, s.IsActive).Scan(&s.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return s, ErrDuplicate
		}
		return s, err
	}
	return s, nil
}

func (r *studentRepository) Update(ctx context.Context, s model.Student) (model.Student, error) {
	query := "UPDATE students SET nim = $1, name = $2, grade = $3, is_active = $4 WHERE id = $5 RETURNING id, nim, name, grade, is_active, created_at"
	err := r.pool.QueryRow(ctx, query, s.NIM, s.Name, s.Grade, s.IsActive, s.ID).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return s, ErrNotFound
	}
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return s, ErrDuplicate
		}
		return s, err
	}
	return s, nil
}

func (r *studentRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM students WHERE id = $1"
	ct, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
