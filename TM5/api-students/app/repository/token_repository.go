package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository interface {
	Save(ctx context.Context, username string, tokenHash string, expiresAt time.Time) error
	FindByHash(ctx context.Context, tokenHash string) (username string, err error)
	DeleteByHash(ctx context.Context, tokenHash string) error
	DeleteByUsername(ctx context.Context, username string) error
}

type tokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Save(ctx context.Context, username string, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (username, token_hash, expires_at, created_at) VALUES ($1, $2, $3, NOW())`
	_, err := r.db.Exec(ctx, query, username, tokenHash, expiresAt)
	return err
}

func (r *tokenRepository) FindByHash(ctx context.Context, tokenHash string) (string, error) {
	query := `SELECT username FROM refresh_tokens WHERE token_hash = $1 AND expires_at > NOW()`
	var username string
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(&username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("refresh token tidak ditemukan atau sudah kedaluwarsa")
		}
		return "", err
	}
	return username, nil
}

func (r *tokenRepository) DeleteByHash(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := r.db.Exec(ctx, query, tokenHash)
	return err
}

func (r *tokenRepository) DeleteByUsername(ctx context.Context, username string) error {
	query := `DELETE FROM refresh_tokens WHERE username = $1`
	_, err := r.db.Exec(ctx, query, username)
	return err
}
