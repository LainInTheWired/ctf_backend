package repository

import (
	"testing"

	"github.com/LainInTheWired/ctf_backend/contest/model"
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

func TestInsertCloudinit(t *testing.T) {
	db, err := NewDBClient()
	if err != nil {
		xerrors.Errorf("mysql connection error: %w", err.Error())
	}
	defer db.Close()

	mr := NewMysqlRepository(db)
	cloudinit := model.Cloudinit{
		QuestionID: 8,
		ContestID:  1,
		Filename:   "",
		TeamID:     1,
		VMID:       115,
		Access:     "highlows",
	}
	err = mr.InsertCloudinit(cloudinit)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Log()
}
