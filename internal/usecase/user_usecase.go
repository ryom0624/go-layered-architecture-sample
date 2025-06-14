package usecase

import (
	"context"
	"errors"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, name, email string) (*entity.User, error)
	GetUser(ctx context.Context, id uint) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id uint) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, name, email string) (*entity.User, error) {
	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}

	existingUser, _ := u.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	user := &entity.User{
		Name:  name,
		Email: email,
	}

	err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.userRepo.GetByEmail(ctx, email)
}

func (u *userUsecase) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	return u.userRepo.GetAll(ctx)
}

func (u *userUsecase) UpdateUser(ctx context.Context, user *entity.User) error {
	if user.ID == 0 {
		return errors.New("user ID is required")
	}
	return u.userRepo.Update(ctx, user)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("user ID is required")
	}
	return u.userRepo.Delete(ctx, id)
}