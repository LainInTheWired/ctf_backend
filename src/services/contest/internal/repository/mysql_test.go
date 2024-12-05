package repository

import (
	"testing"

	"golang.org/x/xerrors"
)

func TestSelectContestQuestionsByContestID(t *testing.T) {
	db, err := NewDBClient()
	if err != nil {
		xerrors.Errorf("mysql connection error: %w", err.Error())
	}
	defer db.Close()

	mr := NewMysqlRepository(db)
	points, err := mr.SelectPointByTeamidAndContestid(1, 1)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", points)
	t.Log()
}
