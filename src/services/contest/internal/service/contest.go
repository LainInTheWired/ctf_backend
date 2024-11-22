package service

import (
	"fmt"

	"github.com/LainInTheWired/ctf_backend/contest/model"
	"github.com/LainInTheWired/ctf_backend/contest/repository"
	"github.com/cockroachdb/errors"
)

type ContestService interface {
	CreateContest(c model.Contest) error
	DeleteContest(c model.Contest) error
	CreateTeamContest(c model.ContestsTeam) error
	DeleteTeamContest(c model.ContestsTeam) error
	ListContest() ([]model.Contest, error)
	ListContestByTeams(tid int) ([]model.Contest, error)
	JoinContestQuesionts(ids []map[string]int) error
	StartContest(cid int) error
}

type contestService struct {
	pveRepo   repository.PVEAPIRepository
	mysqlRepo repository.MysqlRepository
	teamRepo  repository.TeamRepository
	quesRepo  repository.QuestionRepository
}

func NewContestService(pveRepo repository.PVEAPIRepository, mysqlRepo repository.MysqlRepository, teamRepo repository.TeamRepository, quesRepo repository.QuestionRepository) ContestService {
	return &contestService{
		pveRepo:   pveRepo,
		mysqlRepo: mysqlRepo,
		teamRepo:  teamRepo,
		quesRepo:  quesRepo,
	}
}

func (r *contestService) CreateContest(c model.Contest) error {
	// モデルの構造体に移し替えてから、repositoryに渡す
	if err := r.mysqlRepo.InsertContest(c); err != nil {
		return errors.Wrap(err, "can't create contest")
	}
	return nil
}

func (r *contestService) DeleteContest(c model.Contest) error {
	if err := r.mysqlRepo.DeleteContest(c); err != nil {
		return errors.Wrap(err, "can't delete contest")
	}
	return nil
}
func (r *contestService) CreateTeamContest(c model.ContestsTeam) error {
	if err := r.mysqlRepo.InsertTeamContests(c); err != nil {
		return errors.Wrap(err, "can't create team_contests")
	}
	return nil
}
func (r *contestService) DeleteTeamContest(c model.ContestsTeam) error {
	if err := r.mysqlRepo.DeleteTeamContests(c); err != nil {
		return errors.Wrap(err, "can't delete team_contests")
	}
	return nil
}

func (r *contestService) ListContest() ([]model.Contest, error) {
	contests, err := r.mysqlRepo.SelectContest()
	if err != nil {
		return nil, errors.Wrap(err, "can't delete team_contests")
	}
	return contests, nil
}
func (r *contestService) ListContestByTeams(tid int) ([]model.Contest, error) {
	contests, err := r.mysqlRepo.SelectTeamsByContest(tid)
	if err != nil {
		return nil, errors.Wrap(err, "can't delete team_contests")
	}
	return contests, nil
}

func (r *contestService) JoinContestQuesionts(ids []map[string]int) error {
	for _, id := range ids {
		err := r.mysqlRepo.InsertContestsQuestions(id["qid"], id["cid"])
		if err != nil {
			return errors.Wrap(err, "can't delete team_contests")
		}
	}
	return nil
}

func (r *contestService) StartContest(cid int) error {
	teams, err := r.teamRepo.ListTeamUsersByContest(1)
	if err != nil {
		return errors.Wrap(err, "can't get ListTeamUsers")
	}
	fmt.Printf("%+v", teams)
	questions, err := r.quesRepo.GetListQuestionsByContest(1)
	if err != nil {
		return errors.Wrap(err, "can't get ListQuestions")
	}

	for _, team := range teams {
		for _, ques := range questions {
			name := fmt.Sprintf("%d-%d-%d", cid, team.ID, ques.ID)
			m := model.QuesionRequest{
				ID:   ques.VMID,
				Name: name,
			}
			if err := r.quesRepo.CloneQuestion(m); err != nil {
				return errors.Wrap(err, "can't get ListQuestions")
			}
			// cloudinit := model.Cloudinit{
			// 	ContestQuestionsID: ,
			// }
			// r.mysqlRepo.InsertCloudinit()
		}
	}
	return nil
}
