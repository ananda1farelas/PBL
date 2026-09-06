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
	ErrNotFound  = errors.New("data tidak ditemukan")
	ErrDuplicate = errors.New("email sudah terdaftar")
)

type UserRepository interface {
	FindAll(ctx context.Context, q model.ListQuery) ([]model.User, int, error)
	FindByID(ctx context.Context, id int) (model.User, error)
	Create(ctx context.Context, user model.User) (model.User, error)
	Update(ctx context.Context, user model.User) (model.User, error)
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) FindAll(ctx context.Context, q model.ListQuery) ([]model.User, int, error) {
	var conditions []string
	var args []any
	argID := 1

	if q.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR email ILIKE $%d)", argID, argID))
		args = append(args, "%"+q.Search+"%")
		argID++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (q.Page - 1) * q.Limit
	query := fmt.Sprintf(
		"SELECT id, name, email FROM users %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, q.Sort, q.Order, argID, argID+1,
	)
	args = append(args, q.Limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int) (model.User, error) {
	var u model.User
	query := "SELECT id, name, email FROM users WHERE id = $1"
	err := r.pool.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return u, ErrNotFound
	}
	return u, err
}

func (r *userRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
	err := r.pool.QueryRow(ctx, query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return user, ErrDuplicate
		}
		return user, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user model.User) (model.User, error) {
	query := "UPDATE users SET name = $1, email = $2 WHERE id = $3 RETURNING id, name, email"
	err := r.pool.QueryRow(ctx, query, user.Name, user.Email, user.ID).Scan(&user.ID, &user.Name, &user.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return user, ErrNotFound
	}
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return user, ErrDuplicate
		}
		return user, err
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"
	ct, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
