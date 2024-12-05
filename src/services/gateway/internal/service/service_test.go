package service

import (
	"context"
	"testing"

	"github.com/LainInTheWired/ctf_backend/gateway/repository"
	"golang.org/x/xerrors"
)

func TestGetAuthz(t *testing.T) {
	reddb, err := repository.NewRedis()
	if err != nil {
		xerrors.Errorf("redis connetciono error: %w", err.Error())
	}
	defer reddb.Close()
	rr := repository.NewRedisClient(reddb, context.Background())
	s := NewGatewayService(rr)
	roles, err := s.GetRoles("3a0256c6-891b-4903-a963-abffffdaee64")
	if err != nil {
		t.Error(err)
	}
	t.Log(roles)
}
