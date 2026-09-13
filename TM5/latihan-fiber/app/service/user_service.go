package service

import (
	"context"
	"errors"

	"latihan-fiber/app/model"
	"latihan-fiber/app/repository"
)

type UserService interface {
	List(ctx context.Context) ([]model.User, error)
	FindByID(ctx context.Context, id int) (model.User, error)
	Create(ctx context.Context, req model.CreateUserRequest) (model.User, error)
	Update(ctx context.Context, id int, req model.ReplaceUserRequest) (model.User, error)
	Patch(ctx context.Context, id int, req model.PatchUserRequest) (model.User, error)
	Delete(ctx context.Context, id int) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) List(ctx context.Context) ([]model.User, error) {
	return s.repo.List(ctx) // Gunakan List, bukan FindAll
}

func (s *userService) FindByID(ctx context.Context, id int) (model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) Create(ctx context.Context, req model.CreateUserRequest) (model.User, error) {
	if err := ValidateCreate(req); err != nil {
		return model.User{}, errors.New("input tidak valid")
	}

	u := model.User{
		Username: req.Username, // Gunakan Username, bukan Name
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
	}

	return s.repo.Create(ctx, u)
}

func (s *userService) Update(ctx context.Context, id int, req model.ReplaceUserRequest) (model.User, error) {
	u := model.User{
		ID:       id,
		Username: req.Username, // Gunakan Username, bukan Name
		Email:    req.Email,
		Password: req.Password,
	}

	// Sertakan id sebagai argumen kedua
	return s.repo.Update(ctx, id, u)
}

func (s *userService) Patch(ctx context.Context, id int, req model.PatchUserRequest) (model.User, error) {
	current, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	updated, errs := ApplyPatch(current, req)
	if len(errs) > 0 {
		return model.User{}, errors.New("input tidak valid")
	}

	// Sertakan id sebagai argumen kedua
	return s.repo.Update(ctx, id, updated)
}

func (s *userService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
