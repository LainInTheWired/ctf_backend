package repository

import (
	"database/sql"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/LainInTheWired/ctf_backend/question/model"
)

type mysqlRepository struct {
	DB *sql.DB
}
type MysqlRepository interface {
	InsertQuestion(q model.Question) error
	DeleteQuestion(qid int) error
	SelectContestQuestionsByContestID(contestID int) ([]model.Question, error)
	SelectContestQuestions() ([]model.Question, error)
}

func NewMysqlRepository(db *sql.DB) MysqlRepository {
	return &mysqlRepository{
		DB: db,
	}
}

func (m *mysqlRepository) InsertQuestion(q model.Question) error {
	// emailが登録されているかチェック
	ins, err := m.DB.Prepare("INSERT INTO questions (name,env,category_id,description,vmid) VALUES(?,?,?,?,?)")
	if err != nil {
		return errors.Wrap(err, "question insert error")
	}
	defer ins.Close()

	fmt.Printf("INSERT INTO questions (name,env,category_id,describe,vmid) VALUES(%s,%s,%d,%s,%d)", q.Name, q.Env, q.CategoryId, q.Description, q.VMID)

	_, err = ins.Exec(q.Name, q.Env, q.CategoryId, q.Description, q.VMID)
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

func (m *mysqlRepository) SelectContestQuestionsByContestID(contestID int) ([]model.Question, error) {
	var questions []model.Question
	//  emailよりユーザ情報を取得
	// rows, err := m.DB.Query("SELECT id,name,category_id,description,vmid FROM questions WEHERE id = ?", contestID)
	rows, err := m.DB.Query("SELECT q.id,q.name,c.name,q.description,q.vmid FROM questions as q JOIN category AS c ON c.id = q.category_id  WHERE c.id = ?", contestID)
	if err != nil {
		return nil, errors.Wrap(err, "error select contest")
	}

	for rows.Next() {
		q := model.Question{}
		if err := rows.Scan(&q.ID, &q.Name, &q.CategoryName, &q.Description, &q.VMID); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		questions = append(questions, q)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "select errors")
	}

	return questions, nil
}

func (m *mysqlRepository) SelectContestQuestions() ([]model.Question, error) {
	var questions []model.Question
	//  emailよりユーザ情報を取得
	// rows, err := m.DB.Query("SELECT id,name,category_id,description,vmid FROM questions WEHERE id = ?", contestID)
	rows, err := m.DB.Query("SELECT q.id,q.name,c.name,q.description,q.vmid FROM questions as q JOIN category AS c ON c.id = q.category_id")
	if err != nil {
		return nil, errors.Wrap(err, "error select contest")
	}

	for rows.Next() {
		q := model.Question{}
		if err := rows.Scan(&q.ID, &q.Name, &q.CategoryName, &q.Description, &q.VMID); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		questions = append(questions, q)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "select errors")
	}

	return questions, nil
}
