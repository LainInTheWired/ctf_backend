package repository

import (
	"database/sql"

	"github.com/LainInTheWired/ctf_backend/contest/model"
	"github.com/cockroachdb/errors"
)

type MysqlRepository interface {
	InsertContest(contest model.Contest) error
	DeleteContest(contest model.Contest) error
	InsertTeamContests(ct model.ContestsTeam) error
	DeleteTeamContests(ct model.ContestsTeam) error
	SelectContest() ([]model.Contest, error)
	SelectTeamsByContest(tid int) ([]model.Contest, error)
	InsertContestsQuestions(qid, cid int) error
}

type mysqlRepository struct {
	db *sql.DB
}

func NewMysqlRepository(db *sql.DB) MysqlRepository {
	return &mysqlRepository{
		db: db,
	}
}

func (r *mysqlRepository) InsertContest(contest model.Contest) error {
	// emailが登録されているかチェック
	ins, err := r.db.Prepare("INSERT INTO contests (name,start,end) VALUES(? ,?, ?)")
	if err != nil {
		return errors.Wrap(err, "contest insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(contest.Name, contest.StartDate, contest.EndDate)
	if err != nil {
		return errors.Wrap(err, "can't insert conster")
	}
	return nil
}

func (r *mysqlRepository) DeleteContest(contest model.Contest) error {

	// DELETE文を直接実行（Prepareは必要に応じて使用）
	result, err := r.db.Exec("DELETE FROM contests WHERE id = ?", contest.ID)
	if err != nil {
		return errors.Wrap(err, "contest delete error")
	}

	// 影響を受けた行数を取得
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve affected rows")
	}

	// 影響を受けた行数が0の場合、対象のIDが存在しない
	if rowsAffected == 0 {
		return errors.New("no contest found with the given ID")
	}

	return nil
}

func (r *mysqlRepository) InsertTeamContests(ct model.ContestsTeam) error {
	// emailが登録されているかチェック
	ins, err := r.db.Prepare("INSERT INTO contest_teams (contest_id,team_id) VALUES(?,?)")
	if err != nil {
		return errors.Wrap(err, "contest_teams insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(ct.ContestID, ct.TeamID)
	if err != nil {
		return errors.Wrap(err, "can't insert contest_teams")
	}
	return nil
}
func (r *mysqlRepository) DeleteTeamContests(ct model.ContestsTeam) error {
	// DELETE文を直接実行（Prepareは必要に応じて使用）
	result, err := r.db.Exec("DELETE FROM contest_teams WHERE contest_id = ? AND team_id = ?", ct.ContestID, ct.TeamID)
	if err != nil {
		return errors.Wrap(err, "contest_teams delete error")
	}

	// 影響を受けた行数を取得
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve affected rows")
	}

	// 影響を受けた行数が0の場合、対象のIDが存在しない
	if rowsAffected == 0 {
		return errors.New("no contest_teams found with the given ID")
	}
	return nil
}

func (m *mysqlRepository) SelectContest() ([]model.Contest, error) {
	var contests []model.Contest
	//  emailよりユーザ情報を取得
	rows, err := m.db.Query("SELECT id,name,start,end FROM contests")
	if err != nil {
		return nil, errors.Wrap(err, "error select contest")
	}

	for rows.Next() {
		c := model.Contest{}
		if err := rows.Scan(&c.ID, &c.Name, &c.StartDate, &c.EndDate); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		contests = append(contests, c)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "select errors")
	}

	return contests, nil
}

func (m *mysqlRepository) SelectTeamsByContest(tid int) ([]model.Contest, error) {
	var contests []model.Contest
	//  emailよりユーザ情報を取得
	rows, err := m.db.Query("SELECT c.id,c.name,c.start,c.end FROM contest_teams AS ct JOIN contests AS c ON c.id = ct.contest_id WHERE ct.team_id = ?", tid)
	if err != nil {
		return nil, errors.Wrap(err, "error select contest_teams")
	}

	for rows.Next() {
		c := model.Contest{}
		if err := rows.Scan(&c.ID, &c.Name, &c.StartDate, &c.EndDate); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		contests = append(contests, c)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "select errors")
	}
	return contests, nil
}

func (m *mysqlRepository) InsertContestsQuestions(qid, cid int) error {
	// emailが登録されているかチェック
	ins, err := m.db.Prepare("INSERT INTO contest_questions (contest_id,question_id) VALUES(?,?)")
	if err != nil {
		return errors.Wrap(err, "contest_teams insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(cid, qid)
	if err != nil {
		return errors.Wrap(err, "can't insert contest_questions")
	}
	return nil
}
