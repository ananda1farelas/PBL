package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"latihan-fiber/app/model"
)

var (
	ErrNotFound  = errors.New("data tidak ditemukan")
	ErrDuplicate = errors.New("data sudah ada")
)

// Interface utama
type UserRepository interface {
	List(ctx context.Context) ([]model.User, error)
	FindByID(ctx context.Context, id int) (model.User, error)
	FindByUsername(ctx context.Context, username string) (model.User, error)
	Create(ctx context.Context, u model.User) (model.User, error)
	Update(ctx context.Context, id int, u model.User) (model.User, error)
	Delete(ctx context.Context, id int) error
}

type userPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userPostgresRepository{pool: pool}
}

func (r *userPostgresRepository) List(ctx context.Context) ([]model.User, error) {
	query := `SELECT id, username, email, password, role, is_active, created_at FROM users ORDER BY id`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("mengambil daftar user: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role,
			&u.IsActive, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("membaca baris user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userPostgresRepository) FindByID(ctx context.Context, id int) (model.User, error) {
	var u model.User

	err := r.pool.QueryRow(ctx,
		`SELECT id, username, email, password, role, is_active, created_at
         FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role,
		&u.IsActive, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, fmt.Errorf("mengambil user: %w", err)
	}

	return u, nil
}

func (r *userPostgresRepository) FindByUsername(
	ctx context.Context, username string,
) (model.User, error) {
	var u model.User

	err := r.pool.QueryRow(ctx,
		`SELECT id, username, email, password, role, is_active, created_at
         FROM users WHERE LOWER(username) = LOWER($1)`, username,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role,
		&u.IsActive, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, fmt.Errorf("mengambil user: %w", err)
	}

	return u, nil
}

func (r *userPostgresRepository) Create(ctx context.Context, u model.User) (model.User, error) {
	query := `INSERT INTO users (username, email, password, role, is_active)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, query,
		u.Username, u.Email, u.Password, u.Role, u.IsActive,
	).Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.User{}, ErrDuplicate
		}
		return model.User{}, fmt.Errorf("membuat user: %w", err)
	}

	return u, nil
}

func (r *userPostgresRepository) Update(ctx context.Context, id int, u model.User) (model.User, error) {
	query := `UPDATE users
         SET username = $1, email = $2, role = $3, is_active = $4
         WHERE id = $5
         RETURNING id, username, email, password, role, is_active, created_at`

	err := r.pool.QueryRow(ctx, query,
		u.Username, u.Email, u.Role, u.IsActive, id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role,
		&u.IsActive, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.User{}, ErrDuplicate
		}
		return model.User{}, fmt.Errorf("mengubah user: %w", err)
	}

	return u, nil
}

func (r *userPostgresRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"
	ct, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("menghapus user: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

var _ = strings.TrimSpace // hapus baris ini kalau strings sudah dipakai di tempat lain
