package repository

import (
	"database/sql"
	"log"

	"github.com/LainInTheWired/ctf_backend/contest/model"
	"github.com/cockroachdb/errors"
	_ "github.com/go-sql-driver/mysql" // 空のインポートを追加
)

type MysqlRepository interface {
	InsertContest(contest model.Contest) error
	DeleteContest(contest model.Contest) error
	InsertTeamContests(ct model.ContestsTeam) error
	DeleteTeamContests(ct model.ContestsTeam) error
	SelectContest() ([]model.Contest, error)
	// SelectTeamsByContest(tid int) ([]model.Contest, error)
	SelectContestsByTeamID(tid int) ([]model.Contest, error)
	InsertContestsQuestions(qid, cid int) error
	InsertCloudinit(contest model.Cloudinit) error
	SelectPoint(cid int) ([]model.Point, error)
	InsertPoint(tid int, qid int, cid int, point int) error
	SelectContestQuestionsByContestID(cid int) (model.Contest, error)
	SelectPointByTeamidAndContestid(cid int, tid int) ([]model.Point, error)
}

func NewDBClient() (*sql.DB, error) {
	db, err := sql.Open("mysql", "user:user@tcp(db:3306)/ctf?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
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

func (r *mysqlRepository) InsertCloudinit(contest model.Cloudinit) error {
	// emailが登録されているかチェック
	ins, err := r.db.Prepare("INSERT INTO cloudinit (contest_questions_id,team_id,filename) VALUES(? ,?, ?)")
	if err != nil {
		return errors.Wrap(err, "contest insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(contest.ContestQuestionsID, contest.TeamID, contest.Filename)
	if err != nil {
		return errors.Wrap(err, "can't insert cloudinit")
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

func (m *mysqlRepository) SelectContestsByTeamID(tid int) ([]model.Contest, error) {
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

func (m *mysqlRepository) InsertPoint(tid int, qid int, cid int, point int) error {
	ins, err := m.db.Prepare("INSERT INTO points (team_id, question_id, contest_id,point) VALUES(?,?,?,?)")
	if err != nil {
		return errors.Wrap(err, "contest_teams insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(tid, qid, cid, point)
	if err != nil {
		return errors.Wrap(err, "can't insert contest_questions")
	}
	return nil
}

func (m *mysqlRepository) SelectPoint(cid int) ([]model.Point, error) {
	var points []model.Point
	//  emailよりユーザ情報を取得
	rows, err := m.db.Query("SELECT team_id,question_id,contest_id,point,insert_date FROM points WHERE contest_id = ?", cid)
	if err != nil {
		return nil, errors.Wrap(err, "error select contest_teams")
	}

	for rows.Next() {
		p := model.Point{}
		if err := rows.Scan(&p.TeamID, &p.QuestionID, &p.ContestID, &p.Point, &p.InsertDate); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		points = append(points, p)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "select errors")
	}
	return points, nil
}
func (m *mysqlRepository) SelectPointByTeamidAndContestid(cid int, tid int) ([]model.Point, error) {
	var points []model.Point
	//  emailよりユーザ情報を取得
	rows, err := m.db.Query("SELECT question_id,contest_id,point,insert_date FROM points WHERE contest_id = ? AND team_id = ?", cid, tid)
	if err != nil {
		return nil, errors.Wrap(err, "error select contest_teams")
	}

	for rows.Next() {
		p := model.Point{}
		if err := rows.Scan(&p.QuestionID, &p.ContestID, &p.Point, &p.InsertDate); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		points = append(points, p)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "select errors")
	}
	return points, nil
}

func (m *mysqlRepository) SelectContestQuestions() ([]model.Contest, error) {
	var contests []model.Contest
	//  emailよりユーザ情報を取得
	// rows, err := m.DB.Query("SELECT id,name,category_id,description,vmid FROM questions WEHERE id = ?", contestID)
	rows, err := m.db.Query("SELECT c.id,c.name,q.id,q.name,cq.point,cg.name,q.description,q.vmid FROM contest_questions as cq JOIN questions as q ON q.id = cq.question_id JOIN contests  AS c ON  c.id = cq.contest_id JOIN category AS cg ON cg.id = q.category_id;")
	if err != nil {
		return nil, errors.Wrap(err, "error select contest")
	}
	contestMap := make(map[int]*model.Contest)

	for rows.Next() {
		var (
			contestID    int
			contestName  string
			CategoryID   int
			CategoryName string
			questionID   int
			questionName string
			Point        int
			Description  string
			VMID         int
		)
		// すべてのカラムをスキャン
		if err := rows.Scan(&contestID, &contestName, &questionID, &questionName, &CategoryID, &CategoryName, &Point, &Description, &VMID); err != nil {
			return nil, errors.Wrap(err, "SelectTeamUsersInContest: failed to scan row")
		}
		// チームがマップに存在しない場合、追加
		contest, exists := contestMap[contestID]
		if !exists {
			contest = &model.Contest{
				ID:        contestID,
				Name:      contestName,
				Questions: []model.Question{},
			}
			contestMap[contestID] = contest
		}
		// ユーザーをチームに追加
		question := model.Question{
			ID:           questionID,
			Name:         questionName,
			Point:        Point,
			Description:  Description,
			VMID:         VMID,
			CategoryId:   CategoryID,
			CategoryName: contestName,
			// 必要に応じてPasswordフィールドも追加
		}
		contest.Questions = append(contest.Questions, question)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "select errors")
	}
	// マップからスライスへチームを追加
	for _, contest := range contestMap {
		contests = append(contests, *contest)
	}
	return contests, nil
}

func (m *mysqlRepository) SelectContestQuestionsByContestID(cid int) (model.Contest, error) {
	var contest model.Contest
	//  emailよりユーザ情報を取得
	// rows, err := m.DB.Query("SELECT id,name,category_id,description,vmid FROM questions WEHERE id = ?", contestID)
	rows, err := m.db.Query("SELECT c.id,c.name,q.id,q.name,cg.name,cq.point,q.description,q.vmid,q.answer FROM contest_questions as cq JOIN questions as q ON q.id = cq.question_id JOIN contests  AS c ON  c.id = cq.contest_id JOIN category AS cg ON cg.id = q.category_id WHERE c.id = ?;", cid)
	if err != nil {
		return model.Contest{}, errors.Wrap(err, "error select contest")
	}

	for rows.Next() {
		var (
			contestID    int
			contestName  string
			CategoryID   int
			CategoryName string
			questionID   int
			questionName string
			Point        int
			Description  string
			VMID         int
			Answer       string
		)
		// すべてのカラムをスキャン
		if err := rows.Scan(&contestID, &contestName, &questionID, &questionName, &CategoryName, &Point, &Description, &VMID, &Answer); err != nil {
			return model.Contest{}, errors.Wrap(err, "SelectTeamUsersInContest: failed to scan row")
		}
		contest.ID = contestID
		contest.Name = contestName

		// ユーザーをチームに追加
		question := model.Question{
			ID:           questionID,
			Name:         questionName,
			Point:        Point,
			Description:  Description,
			VMID:         VMID,
			CategoryId:   CategoryID,
			CategoryName: contestName,
			Answer:       Answer,
			// 必要に応じてPasswordフィールドも追加
		}
		contest.Questions = append(contest.Questions, question)
	}
	if err = rows.Err(); err != nil {
		return model.Contest{}, errors.Wrap(err, "select errors")
	}

	return contest, nil
}
