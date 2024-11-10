package repository

import (
	"database/sql"

	"github.com/LainInTheWired/ctf_backend/team/model"
	"github.com/cockroachdb/errors"
)

type mysqlRepository struct {
	DB *sql.DB
}
type MysqlRepository interface {
	InsertTeam(team model.Team) error
	DeleteTeam(team model.Team) error
	SelectContestByTeam(cid int) ([]model.Team, error)
	InsertContestTeams(ct model.ContestTeams) error
}

func NewMysqlRepository(db *sql.DB) MysqlRepository {
	return &mysqlRepository{
		DB: db,
	}
}

func (m *mysqlRepository) InsertTeam(team model.Team) error {
	ins, err := m.DB.Prepare("INSERT INTO teams (name)  VALUES(?)")
	if err != nil {
		return errors.Wrap(err, "team insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(team.Name)
	if err != nil {
		return errors.Wrap(err, "can't insert team")
	}
	return nil
}
func (m *mysqlRepository) DeleteTeam(team model.Team) error {
	// DELETE文を直接実行（Prepareは必要に応じて使用）
	result, err := m.DB.Exec("DELETE FROM teams WHERE id = ?", team.ID)
	if err != nil {
		return errors.Wrap(err, "team delete error")
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

func (m *mysqlRepository) SelectContestByTeam(cid int) ([]model.Team, error) {
	var teams []model.Team
	//  emailよりユーザ情報を取得
	rows, err := m.DB.Query("SELECT t.id,t.name FROM contest_teams AS ct JOIN teams AS t ON t.id = ct.team_id WHERE ct.contest_id = ?", cid)
	if err != nil {
		return nil, errors.Wrap(err, "error select contest_teams")
	}

	for rows.Next() {
		t := model.Team{}
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		teams = append(teams, t)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "select errors")
	}

	// if err := m.DB.QueryRow("SELECT t.id,t.name FROM contest_teams JOIN teams AS t ON t.ID = ct.team_ID WHERE ct.team_id = ?", cid).Scan(&team.ID, &team.Name); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return model.Team{}, errors.Wrap(err, "not exist this team")
	// 	}
	// 	return model.Team{}, errors.Wrap(err, "can't select user by email")

	// }
	return teams, nil
}

func (m *mysqlRepository) InsertContestTeams(ct model.ContestTeams) error {
	ins, err := m.DB.Prepare("INSERT INTO contest_teams (contest_id,team_id)  VALUES(?,?)")
	if err != nil {
		return errors.Wrap(err, "team insert contest_teams")
	}
	defer ins.Close()

	_, err = ins.Exec(ct.TeamID)
	if err != nil {
		return errors.Wrap(err, "can't insert contest_teams")
	}
	return nil
}
