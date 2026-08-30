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
	Patch(ctx context.Context, id int64, req model.PatchStudentRequest) (*model.Student, error)
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

	// 1. Pencarian Nama / NIM (ILIKE)
	if q.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(nim ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+q.Search+"%")
		argIdx++
	}

	// 2. Penyaringan Field Status Aktif (is_active)
	if q.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *q.IsActive)
		argIdx++
	}

	// Penggabungan Klausa WHERE
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 3. Menghitung Metadata Total Records (meta.total)
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM students %s", whereClause)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 4. Whitelist Pengurutan Kolom (ORDER BY)
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
	if strings.ToUpper(q.Order) == "DESC" || strings.ToUpper(q.Sort) == "DESC" {
		order = "DESC"
	}

	// 5. Normalisasi Paginasi (LIMIT & OFFSET)
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}
	offset := (q.Page - 1) * q.Limit

	// Query Utama
	query := fmt.Sprintf(
		"SELECT id, nim, name, grade, is_active, created_at FROM students %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, sortCol, order, argIdx, argIdx+1,
	)
	args = append(args, q.Limit, offset)

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

func (r *postgresStudentRepository) Patch(ctx context.Context, id int64, req model.PatchStudentRequest) (*model.Student, error) {
	var sets []string
	var args []any
	argIdx := 1

	if req.NIM != nil {
		nim := strings.TrimSpace(*req.NIM)
		if nim == "" {
			return nil, errors.New("nim tidak boleh kosong")
		}
		sets = append(sets, fmt.Sprintf("nim = $%d", argIdx))
		args = append(args, nim)
		argIdx++
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("name tidak boleh kosong")
		}
		sets = append(sets, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, name)
		argIdx++
	}

	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 4.0 {
			return nil, errors.New("grade harus di antara 0.0 sampai 4.0")
		}
		sets = append(sets, fmt.Sprintf("grade = $%d", argIdx))
		args = append(args, *req.Grade)
		argIdx++
	}

	if req.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *req.IsActive)
		argIdx++
	}

	if len(sets) == 0 {
		return r.FindByID(ctx, id)
	}

	query := fmt.Sprintf(
		"UPDATE students SET %s WHERE id = $%d RETURNING id, nim, name, grade, is_active, created_at",
		strings.Join(sets, ", "), argIdx,
	)
	args = append(args, id)

	var s model.Student
	err := r.pool.QueryRow(ctx, query, args...).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)
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
