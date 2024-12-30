package service

import (
	"github.com/LainInTheWired/ctf_backend/team/model"
	"github.com/LainInTheWired/ctf_backend/team/repository"
	"github.com/cockroachdb/errors"
)

type teamService struct {
	repo repository.MysqlRepository
}

type TeamService interface {
	DeleteTeam(t model.Team) error
	CreateTeam(t model.Team) error
	ListTeamByContest(cid int) ([]model.Team, error)
	JoinContest(ct model.ContestTeams) error
	ListTeamUsersByContest(cid int) ([]model.Team, error)
	ListTeamInContestByUserID(cid, uid int) ([]model.Team, error)
	ListUsers() ([]model.User, error)
}

func NewTeamService(r repository.MysqlRepository) TeamService {
	return &teamService{
		repo: r,
	}
}

func (s *teamService) CreateTeam(t model.Team) error {
	if err := s.repo.InsertTeam(t); err != nil {
		return errors.Wrap(err, "can't create team")
	}
	return nil
}

func (s *teamService) DeleteTeam(t model.Team) error {
	if err := s.repo.DeleteTeam(t); err != nil {
		return errors.Wrap(err, "can't delete team")
	}
	return nil
}

func (s *teamService) ListTeamByContest(cid int) ([]model.Team, error) {
	teams, err := s.repo.SelectTeamInContest(cid)
	if err != nil {
		return nil, errors.Wrap(err, "can't delete team")
	}
	return teams, nil
}

func (s *teamService) JoinContest(ct model.ContestTeams) error {
	if err := s.repo.InsertContestTeams(ct); err != nil {
		return errors.Wrap(err, "can't join consert team")
	}
	return nil
}

func (s *teamService) ListTeamUsersByContest(cid int) ([]model.Team, error) {
	teams, err := s.repo.SelectTeamUsersInContest(cid)
	if err != nil {
		return nil, errors.Wrap(err, "can't delete team")
	}
	return teams, nil
}

func (s *teamService) ListTeamInContestByUserID(cid, uid int) ([]model.Team, error) {
	teams, err := s.repo.SelectTeamUsersInContestByUserID(cid, uid)
	if err != nil {
		return nil, errors.Wrap(err, "can't delete team")
	}
	return teams, nil
}

func (s *teamService) ListUsers() ([]model.User, error) {
	users, err := s.repo.SelectUsers()
	if err != nil {
		return nil, errors.Wrap(err, "error select User")
	}
	return users, nil
}
