package service

import (
	"crypto/rand"
	"fmt"
	"math/big"

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
	GetPoints(cid int) ([]model.ResponsePoints, error)
	CheckQuestion(cid int, qid int, tid int, ans string) (bool, error)
	ListQuestionsByContestID(cid int, tid int) (*model.Contest, error)
	GetTeamByUserID(cid int, uid int) ([]model.Team, error)
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
	contests, err := r.mysqlRepo.SelectContestsByTeamID(tid)
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
	teams, err := r.teamRepo.ListTeamUsersByContest(cid, nil)
	if err != nil {
		return errors.Wrap(err, "can't get ListTeamUsers")
	}
	fmt.Printf("%+v", teams)
	questions, err := r.quesRepo.GetListQuestionsByContest(cid)
	if err != nil {
		return errors.Wrap(err, "can't get ListQuestions")
	}

	for _, team := range teams {
		for _, ques := range questions {
			password, err := generatePassword(16)
			if err != nil {
				return errors.Wrap(err, "can't generate password")
			}
			name := fmt.Sprintf("%d-%d-%d", cid, team.ID, ques.ID)
			m := model.QuesionRequest{
				ID:       ques.VMID,
				Name:     name,
				Password: password,
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

func generatePassword(length int) (string, error) {
	var (
		lowerLetters = "abcdefghijklmnopqrstuvwxyz"
		upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits       = "0123456789"
		symbols      = "!@#$%^&*"
		allChars     = lowerLetters + upperLetters + digits + symbols
	)

	password := make([]byte, length)
	for i := range password {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password[i] = allChars[index.Int64()]
	}

	return string(password), nil
}
func (r *contestService) GetPoints(cid int) ([]model.ResponsePoints, error) {
	var res []model.ResponsePoints
	teams, err := r.teamRepo.ListTeamUsersByContest(cid, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't get point")
	}

	points, err := r.mysqlRepo.SelectPoint(cid)
	if err != nil {
		return nil, errors.Wrap(err, "can't get point")
	}
	for _, team := range teams {
		var tpoints []model.Point
		for _, point := range points {
			if point.TeamID == team.ID {
				tpoint := model.Point{
					Point:      point.Point,
					InsertDate: point.InsertDate,
				}
				tpoints = append(tpoints, tpoint)
			}
		}
		t := model.ResponsePoints{
			TeamID: team.ID,
			Name:   team.Name,
			Points: tpoints,
		}
		res = append(res, t)
	}
	return res, nil
}

func (r *contestService) CheckQuestion(cid int, qid int, tid int, ans string) (bool, error) {
	contest, err := r.mysqlRepo.SelectContestQuestionsByContestID(cid)
	if err != nil {
		return false, errors.Wrap(err, "can't get Questions")
	}
	question := FilterQuestionsByID(contest.Questions, qid)
	if question == nil {
		return false, errors.Wrap(err, "can't filter quesion")
	}
	if question.Answer == ans {
		if err := r.mysqlRepo.InsertPoint(tid, qid, cid, question.Point); err != nil {
			return false, errors.Wrap(err, "can't get Questions")
		}
		return true, nil
	}

	return false, nil
}

// FilterQuestionsByID 指定されたIDでフィルタリングする関数
func FilterQuestionsByID(questions []model.Question, id int) *model.Question {
	for _, q := range questions {
		fmt.Printf("questions: %+v\n", q)
		if q.ID == id {
			return &q
		}
	}
	return nil
}

func (s *contestService) ListQuestionsByContestID(cid int, tid int) (*model.Contest, error) {
	contests, err := s.mysqlRepo.SelectContestQuestionsByContestID(cid)
	if err != nil {
		return nil, errors.Wrap(err, "get questions")
	}
	points, err := s.mysqlRepo.SelectPointByTeamidAndContestid(cid, tid)
	fmt.Printf("%+v", points)
	if err != nil {
		return nil, errors.Wrap(err, "get questions")
	}
	// pointsをマップに変換
	pointMap := make(map[int]int)
	for _, point := range points {
		if _, exists := pointMap[point.QuestionID]; !exists {
			pointMap[point.QuestionID] = point.Point
		}
	}
	// contests.Questionsを更新
	for i := range contests.Questions {
		if point, exists := pointMap[contests.Questions[i].ID]; exists {
			contests.Questions[i].CurrentPoint = point
		}
	}
	return &contests, nil
}

func (s *contestService) GetTeamByUserID(cid int, uid int) ([]model.Team, error) {
	fmt.Printf("%d\n", uid)
	teams, err := s.teamRepo.ListTeamUsersByContest(cid, &uid)
	if err != nil {
		return nil, nil
	}
	return teams, nil
}
