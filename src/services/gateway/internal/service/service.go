package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/LainInTheWired/ctf_backend/gateway/model"
	"github.com/LainInTheWired/ctf_backend/gateway/repository"
	"github.com/cockroachdb/errors"
)

type GatewayService interface {
	GetRoles(userid int) ([]model.Role, error)
	GetUserID(sessionid string) (int, error)
}
type gatewayService struct {
	rerepo repository.RedisRepository
}

func NewGatewayService(redisrepository repository.RedisRepository) GatewayService {
	return &gatewayService{
		rerepo: redisrepository,
	}
}

func (s *gatewayService) CheckAuthz(sessionid string) error {
	return nil
}

func (s *gatewayService) GetUserID(sessionid string) (int, error) {
	suserid, err := s.rerepo.Get(fmt.Sprintf("session:%s", sessionid))
	if err != nil {
		return 0, errors.Wrap(err, "")
	}
	if err := s.rerepo.Expire(fmt.Sprintf("session:%s", sessionid), time.Hour); err != nil {
		return 0, errors.Wrap(err, "")
	}
	if suserid == "" {
		return 0, errors.Wrap(err, "not found session")

	}
	userid, err := strconv.Atoi(suserid)
	if err != nil {
		return 0, errors.Wrap(err, "")
	}
	return userid, nil
}

func (s *gatewayService) GetRoles(userid int) ([]model.Role, error) {
	jusers, err := s.rerepo.Get(fmt.Sprintf("user:%d", userid))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if err := s.rerepo.Expire(fmt.Sprintf("user:%d", userid), time.Hour); err != nil {
		return nil, errors.Wrap(err, "")
	}
	var users []model.User
	if err := json.Unmarshal([]byte(jusers), &users); err != nil {
		return nil, errors.Wrap(err, "")
	}
	var roles []model.Role
	for _, user := range users {
		jrole, err := s.rerepo.Get(fmt.Sprintf("role:%d", user.ID))
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		var role model.Role
		if err := json.Unmarshal([]byte(jrole), &role); err != nil {
			return nil, errors.Wrap(err, "")
		}
		roles = append(roles, role)

	}
	return roles, nil
}
