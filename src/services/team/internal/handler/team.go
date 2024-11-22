package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LainInTheWired/ctf_backend/team/model"
	"github.com/LainInTheWired/ctf_backend/team/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/xerrors"
)

type TeamHander interface {
	CreateTeam(c echo.Context) error
	DeleteTeam(c echo.Context) error
	ListTeamByContest(c echo.Context) error
	ListTeamUserByContest(c echo.Context) error
}

type teamHander struct {
	serv service.TeamService
}

type createTeamRequest struct {
	Name string `json:"name" validate:"required"`
}

type deleteTeamRequest struct {
	ID int `json:"id" validate:"required"`
}

type listTeamByContestRequest struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type listTeamByContestResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type joinContestRequest struct {
	ContestID int `json:"contest_id" validate:"required"`
	TeamID    int `json:"team_id" validate:"required"`
}

func NewTeamHandler(sv service.TeamService) TeamHander {
	return &teamHander{
		serv: sv,
	}
}

func (t *teamHander) CreateTeam(c echo.Context) error {
	var req createTeamRequest
	if err := c.Bind(&req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	m := model.Team{
		Name: req.Name,
	}
	if err := t.serv.CreateTeam(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "create teams "))
}

func (t *teamHander) DeleteTeam(c echo.Context) error {
	var req deleteTeamRequest
	if err := c.Bind(&req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	m := model.Team{
		ID: req.ID,
	}

	if err := t.serv.DeleteTeam(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "Delete teams "))
}

func (t *teamHander) ListTeamByContest(c echo.Context) error {
	sid := c.QueryParam("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}

	teams, err := t.serv.ListTeamByContest(id)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusOK, teams)
}

func (t *teamHander) JoinContest(c echo.Context) error {
	var req joinContestRequest
	if err := c.Bind(&req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	m := model.ContestTeams{
		ContestID: req.ContestID,
		TeamID:    req.TeamID,
	}
	if err := t.serv.JoinContest(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "Join contest_teams"))
}

func (t *teamHander) ListTeamUserByContest(c echo.Context) error {
	sid := c.QueryParam("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}

	teams, err := t.serv.ListTeamUsersByContest(id)

	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	fmt.Printf("%+v", teams[0])
	return c.JSON(http.StatusOK, teams)
}
