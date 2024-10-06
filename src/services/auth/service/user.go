package service

import (
	"github.com/LainInTheWired/ctf_backend/user/model"
	"github.com/LainInTheWired/ctf_backend/user/repository"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
)

type UserService interface {
	Signup(model.User) error
	Login(model.User) (model.User, error)
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		repo: repository,
	}
}

func (s *userService) Signup(u model.User) error {
	// モデルの構造体に移し替えてから、repositoryに渡す
	if err := s.repo.CreateUser(u); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (s *userService) Login(u model.User) (model.User, error) {
	getUser, err := s.repo.GetUserByEmail(u.Email)
	if err != nil {
		return model.User{}, xerrors.Errorf(": %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(getUser.Password), []byte(u.Password)); err != nil {
		return model.User{}, xerrors.Errorf("not same password: %w", err)
	}
	return getUser, nil
}
