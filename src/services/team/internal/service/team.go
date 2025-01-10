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
	CreateTeam(t model.Team, uids []int) (int, error)
	EditTeam(t model.Team, uids []int) error
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

func (s *teamService) CreateTeam(t model.Team, uids []int) (int, error) {
	tid, err := s.repo.InsertTeam(t)
	if err != nil {
		return 0, errors.Wrap(err, "can't create team")
	}
	for _, uid := range uids {
		if err := s.repo.InsertTeamUsers(tid, uid); err != nil {
			return 0, errors.Wrap(err, "can't create user_teams")
		}
	}
	return tid, nil
}
func (s *teamService) EditTeam(t model.Team, uids []int) error {
	if err := s.repo.UpdateTeam(t); err != nil {
		return errors.Wrap(err, "can't update teams")

	}

	users, err := s.repo.SelectUsersInTeamID(t.ID)
	if err != nil {
		return errors.Wrap(err, "can't edit team")
	}
	incominguser := map[int]string{}
	for _, uid := range uids {
		incominguser[uid] = ""
	}
	existusers := map[int]model.User{}
	for _, u := range users {
		existusers[u.ID] = u
	}
	for _, uid := range uids {
		if _, exist := existusers[uid]; !exist {
			if err := s.repo.InsertTeamUsers(t.ID, uid); err != nil {
				return errors.Wrap(err, "can't insert team users")
			}
		}
	}
	for exuid := range existusers {
		if _, exist := incominguser[exuid]; !exist {
			if err := s.repo.DeleteTeamUsers(t.ID, exuid); err != nil {
				return errors.Wrap(err, "can't delete team users")
			}
		}
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
