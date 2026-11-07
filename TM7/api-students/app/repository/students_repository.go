package repository

import (
	"context"
	"errors"
	"fmt"

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
	FindAfterCursor(ctx context.Context, q model.CursorQuery) ([]model.Student, error)
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

// FindAfterCursor mengambil satu halaman dengan keyset pagination.
// Urutan created_at DESC, id DESC memakai id sebagai pemecah seri unik.
func (r *studentRepository) FindAfterCursor(ctx context.Context, q model.CursorQuery) ([]model.Student, error) {
	query := `SELECT id, nim, name, grade, is_active, created_at
		FROM students`
	args := make([]any, 0, 3)

	if q.After != nil {
		query += ` WHERE (created_at, id) < ($1::timestamp, $2)`
		args = append(args, q.After.CreatedAt, q.After.ID)
	}

	args = append(args, q.Limit+1)
	query += fmt.Sprintf(` ORDER BY created_at DESC, id DESC LIMIT $%d`, len(args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := make([]model.Student, 0, q.Limit+1)
	for rows.Next() {
		var student model.Student
		if err := rows.Scan(
			&student.ID, &student.NIM, &student.Name, &student.Grade,
			&student.IsActive, &student.CreatedAt,
		); err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
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
