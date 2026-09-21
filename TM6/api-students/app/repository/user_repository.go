package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"api-students/app/model"
)

type UserRepository interface {
	Create(ctx context.Context, username string, hashedPassword string, role string) error
	FindByUsername(ctx context.Context, username string) (*model.Student, error)
	FindByID(ctx context.Context, id int) (*model.Student, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, username string, hashedPassword string, role string) error {
	query := `INSERT INTO students (username, password, role, nim, name, grade, is_active, created_at) VALUES ($1, $2, $3, '', '', NULL, true, NOW())`
	_, err := r.pool.Exec(ctx, query, username, hashedPassword, role)
	return err
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.Student, error) {
	query := `SELECT id, username, password, role, nim, name, grade, is_active, created_at FROM students WHERE username = $1`
	return r.findBy(ctx, query, username)
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*model.Student, error) {
	query := `SELECT id, username, password, role, nim, name, grade, is_active, created_at FROM students WHERE id = $1`
	return r.findBy(ctx, query, id)
}

func (r *userRepository) findBy(ctx context.Context, query string, args ...any) (*model.Student, error) {
	var s model.Student
	err := r.pool.QueryRow(ctx, query, args...).Scan(&s.ID, &s.Username, &s.Password, &s.Role, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	return &s, nil
}