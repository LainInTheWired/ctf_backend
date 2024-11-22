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
	SelectTeamInContest(cid int) ([]model.Team, error)
	InsertContestTeams(ct model.ContestTeams) error
	SelectTeamUsersInContest(cid int) ([]model.Team, error)
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

func (m *mysqlRepository) SelectTeamInContest(cid int) ([]model.Team, error) {
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

func (m *mysqlRepository) SelectTeamUsersInContest(cid int) ([]model.Team, error) {
	var teams []model.Team
	//  emailよりユーザ情報を取得
	rows, err := m.DB.Query("SELECT t.id, t.name,tu.user_id,u.name,u.email  FROM contest_teams AS ct JOIN teams AS t ON t.id = ct.team_id JOIN team_users AS tu ON t.id = tu.team_id  JOIN users AS u ON u.id = tu.user_id   WHERE ct.contest_id = ?", cid)
	if err != nil {
		return nil, errors.Wrap(err, "error select contest_teams")
	}
	defer rows.Close()

	// チームIDをキーとするマップを作成
	teamMap := make(map[int]*model.Team)

	for rows.Next() {
		var (
			teamID    int
			teamName  string
			userID    int
			userName  string
			userEmail string
		)

		// すべてのカラムをスキャン
		if err := rows.Scan(&teamID, &teamName, &userID, &userName, &userEmail); err != nil {
			return nil, errors.Wrap(err, "SelectTeamUsersInContest: failed to scan row")
		}

		// チームがマップに存在しない場合、追加
		team, exists := teamMap[teamID]
		if !exists {
			team = &model.Team{
				ID:    teamID,
				Name:  teamName,
				Users: []model.User{},
			}
			teamMap[teamID] = team
		}

		// ユーザーをチームに追加
		user := model.User{
			ID:    userID,
			Name:  userName,
			Email: userEmail,
			// 必要に応じてPasswordフィールドも追加
		}
		team.Users = append(team.Users, user)
	}

	// エラーチェック
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "SelectTeamUsersInContest: rows error")
	}

	// マップからスライスへチームを追加
	for _, team := range teamMap {
		teams = append(teams, *team)
	}
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
