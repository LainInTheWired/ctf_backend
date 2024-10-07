package service

import (
	"strconv"
	"time"

	"github.com/LainInTheWired/ctf_backend/user/model"
	"github.com/LainInTheWired/ctf_backend/user/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
)

type UserService interface {
	Signup(model.User) error
	Login(u model.User) (string, error)
	CheckSession(sessionID string) (string, error)
}
type userService struct {
	usrepo repository.UserRepository
	rerepo repository.RedisRepository
}

func NewUserService(userrepository repository.UserRepository, redisrepository repository.RedisRepository) UserService {
	return &userService{
		usrepo: userrepository,
		rerepo: redisrepository,
	}
}

func (s *userService) Signup(u model.User) error {
	// モデルの構造体に移し替えてから、repositoryに渡す
	if err := s.usrepo.CreateUser(u); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}

func (s *userService) Login(u model.User) (string, error) {
	getUser, err := s.usrepo.GetUserByEmail(u.Email)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(getUser.Password), []byte(u.Password)); err != nil {
		return "", xerrors.Errorf("authentication failed: %w", err)
	}
	sessionID, err := NewSession(getUser.Id, s)
	if err != nil {
		return "", xerrors.Errorf("authentication failed: %w", err)
	}
	return sessionID, nil
}

func (s *userService) CheckSession(sessionID string) (string, error) {
	sessionData, err := s.rerepo.Get(sessionID)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}

	return sessionData, nil
}

func NewSession(userID int, s *userService) (string, error) {
	sessionID := uuid.New().String()

	uid := strconv.Itoa(userID)

	// session := map[string]string{
	// 	"userid":    uid,
	// 	"CreatedAt": time.Now().Format(time.RFC3339),
	// }

	err := s.rerepo.Set(sessionID, uid, time.Hour)
	if err != nil {
		return "", xerrors.Errorf("redis can't set sesstionID: %w", err)
	}
	return sessionID, nil
}

func (s *userService) GetSession(sessionID string) (string, error) {
	r, err := s.rerepo.Get(sessionID)
	if err != nil {
		return "", xerrors.Errorf("redis can't get : %w", err)
	}
	return r, nil
}
