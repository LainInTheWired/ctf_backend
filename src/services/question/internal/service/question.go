package service

import (
	"github.com/LainInTheWired/ctf_backend/question/model"
	"github.com/LainInTheWired/ctf_backend/question/repository"
	"github.com/cockroachdb/errors"
)

type quesionService struct {
	myrepo repository.MysqlRepository
}

type QuesionService interface {
	CreateQuestion(q model.Question) error
	DeleteQuestion(qid int) error
}

func NewQuestionService(r repository.MysqlRepository) QuesionService {
	return &quesionService{
		myrepo: r,
	}

}

func (s *quesionService) CreateQuestion(q model.Question) error {
	// モデルの構造体に移し替えてから、repositoryに渡す
	if err := s.myrepo.InsertQuestion(q); err != nil {
		return errors.Wrap(err, "can't create contest")
	}
	return nil
}

func (s *quesionService) DeleteQuestion(qid int) error {
	// モデルの構造体に移し替えてから、repositoryに渡す
	if err := s.myrepo.DeleteQuestion(qid); err != nil {
		return errors.Wrap(err, "can't create contest")
	}
	return nil
}
