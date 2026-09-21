package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository interface {
	Save(ctx context.Context, studentID int, tokenHash string, expiresAt time.Time) error
	FindByHash(ctx context.Context, tokenHash string) (studentID int, err error)
	DeleteByHash(ctx context.Context, tokenHash string) error
	DeleteByStudentID(ctx context.Context, studentID int) error
}

type tokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Save(ctx context.Context, studentID int, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (student_id, token_hash, expires_at, created_at) VALUES ($1, $2, $3, NOW())`
	_, err := r.db.Exec(ctx, query, studentID, tokenHash, expiresAt)
	return err
}

func (r *tokenRepository) FindByHash(ctx context.Context, tokenHash string) (int, error) {
	query := `SELECT student_id FROM refresh_tokens WHERE token_hash = $1 AND expires_at > NOW()`
	var studentID int
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(&studentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("refresh token tidak ditemukan atau sudah kedaluwarsa")
		}
		return 0, err
	}
	return studentID, nil
}

func (r *tokenRepository) DeleteByHash(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := r.db.Exec(ctx, query, tokenHash)
	return err
}

func (r *tokenRepository) DeleteByStudentID(ctx context.Context, studentID int) error {
	query := `DELETE FROM refresh_tokens WHERE student_id = $1`
	_, err := r.db.Exec(ctx, query, studentID)
	return err
}