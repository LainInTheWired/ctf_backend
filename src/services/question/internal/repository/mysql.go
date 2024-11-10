package repository

import (
	"database/sql"

	"github.com/cockroachdb/errors"

	"github.com/LainInTheWired/ctf_backend/question/model"
)

type mysqlRepository struct {
	DB *sql.DB
}
type MysqlRepository interface {
	InsertQuestion(q model.Question) error
	DeleteQuestion(qid int) error
}

func NewMysqlRepository(db *sql.DB) MysqlRepository {
	return &mysqlRepository{
		DB: db,
	}
}
func (m *mysqlRepository) InsertQuestion(q model.Question) error {
	// emailが登録されているかチェック
	ins, err := m.DB.Prepare("INSERT INTO questions (name,env) VALUES(?,?)")
	if err != nil {
		return errors.Wrap(err, "question insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(q.Name, q.Env)
	if err != nil {
		return errors.Wrap(err, "can't insert question")
	}
	return nil
}
func (m *mysqlRepository) DeleteQuestion(qid int) error {
	// emailが登録されているかチェック
	ins, err := m.DB.Prepare("DELETE FROM questions WHERE id = ?")
	if err != nil {
		return errors.Wrap(err, "question insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(qid)
	if err != nil {
		return errors.Wrap(err, "can't insert question")
	}
	return nil
}
