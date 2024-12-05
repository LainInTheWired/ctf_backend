package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/LainInTheWired/ctf_backend/user/model"
	"github.com/LainInTheWired/ctf_backend/user/repository"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Signup(model.User) error
	Login(u model.User) (string, error)
	CheckSession(string) (string, error)
	AddRole(*model.Role) (int, error)
	AddPermission(*model.Permission) error
	BindRolePermissions(int, int) error
	BindUserRoles(uid int, rid int) error
	SetInitRolePermissionsToRedis() error
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
		return errors.Wrap(err, "can't create user")
	}
	return nil
}

func (s *userService) Login(u model.User) (string, error) {
	getUser, err := s.usrepo.GetUserByEmail(u.Email)
	if err != nil {
		return "", errors.Wrap(err, "can't get user by email")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(getUser.Password), []byte(u.Password)); err != nil {
		return "", errors.Wrap(err, "authentication failed")

	}
	sessionID, err := NewSession(getUser.ID, s)
	if err != nil {
		return "", errors.Wrap(err, "authentication failed")
	}
	if err := s.SetRolePermissionsToRedis(getUser.ID); err != nil {
		return "", errors.Wrap(err, "can't set role permission to redis")
	}
	return sessionID, nil
}

func (s *userService) CheckSession(sessionID string) (string, error) {
	sessionData, err := s.rerepo.Get(fmt.Sprintf("session:%s", sessionID))
	if err != nil {
		return "", errors.Wrap(err, "can't get sessionID")
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

	err := s.rerepo.Set(fmt.Sprintf("session:%s", sessionID), uid, time.Hour)
	if err != nil {
		return "", errors.Wrap(err, "redis can't set sesstionID")
	}
	return sessionID, nil
}

func (s *userService) AddRole(role *model.Role) (int, error) {
	rid, err := s.usrepo.CreateRole(role)
	if err != nil {
		return 0, errors.Wrap(err, "can't insert roles")
	}
	return rid, nil
}

func (s *userService) BindRolePermissions(rid int, pid int) error {
	err := s.usrepo.BindRolePermissions(rid, pid)
	if err != nil {
		return errors.Wrap(err, "can't bind role permission")
	}
	return nil
}
func (s *userService) BindUserRoles(uid int, rid int) error {
	err := s.usrepo.BindUserRoles(uid, rid)
	if err != nil {
		return errors.Wrap(err, "can't bind user roles")
	}
	return nil
}
func (s *userService) AddPermission(permission *model.Permission) error {
	err := s.usrepo.CreatePermission(permission)
	if err != nil {
		return errors.Wrap(err, "can't insert roles")
	}
	return nil
}

func (s *userService) Logout(sessionID string) error {
	err := s.rerepo.Delete(sessionID)
	if err != nil {
		return errors.Wrap(err, "redis can't delete key")
	}
	return nil
}

func (s *userService) SetRolePermissionsToRedis(userid int) error {
	ur, err := s.usrepo.GetUserWithRoles(userid)
	if err != nil {
		return errors.Wrap(err, "can't get user role")
	}
	jur, err := json.Marshal(ur.Role)
	if err != nil {
		return errors.Wrap(err, "can't json marshal")
	}
	if err := s.rerepo.Set(fmt.Sprintf("user:%d", ur.ID), jur, time.Hour); err != nil {
		return errors.Wrap(err, "can't set redis")
	}
	return nil
}

func (s *userService) SetInitRolePermissionsToRedis() error {
	rps, err := s.usrepo.GetRolePermissions()
	if err != nil {
		return errors.Wrap(err, "can't get role permission")
	}

	for _, rp := range rps {
		jrp, err := json.Marshal(rp)
		if err != nil {
			return errors.Wrap(err, "can't marshal json")
		}

		s.rerepo.Set(fmt.Sprintf("role:%d", rp.ID), jrp, 0)
	}
	return nil
}
