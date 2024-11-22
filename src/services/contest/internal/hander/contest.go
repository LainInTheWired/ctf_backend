package hander

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/LainInTheWired/ctf_backend/contest/model"
	"github.com/LainInTheWired/ctf_backend/contest/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/xerrors"
)

type ContestHander interface {
	CreateContest(c echo.Context) error
	DeleteContest(c echo.Context) error
	JoinTeamsinContest(c echo.Context) error
	DeleteTeamContest(c echo.Context) error
	ListContest(c echo.Context) error
	ListContestByTeams(c echo.Context) error
	StartContest(c echo.Context) error
}

type contestHander struct {
	serv service.ContestService
}

type createContestRequest struct {
	Name      string `json"name" validate:"required"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type deleteContestRequest struct {
	ID int `json:"id" validate:"required"`
}
type contestJoinTeamRquest struct {
	ContestID int `json:"contest_id" validate:"required"`
	TeamID    int `json:"team_id" validate:"required"`
}
type listContestResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
type joinContestQuesiontsRequest struct {
	QID int `json:"qid"`
	CID int `json:"cid"`
}
type startContestRequest struct {
	contestID int `json:"contest_id"`
}

func NewContestHander(cr service.ContestService) ContestHander {
	return &contestHander{
		serv: cr,
	}
}

func (h *contestHander) CreateContest(c echo.Context) error {
	// リクエストから構造体にデータをコピー
	var req createContestRequest
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

	layout := "2006-01-02 15:04:05"

	st, err := time.Parse(layout, req.StartDate)
	et, err := time.Parse(layout, req.EndDate)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	m := model.Contest{
		Name:      req.Name,
		StartDate: st,
		EndDate:   et,
	}
	fmt.Println(m)

	if err := h.serv.CreateContest(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "create contest"))
}

func (h *contestHander) DeleteContest(c echo.Context) error {
	var req deleteContestRequest
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

	m := model.Contest{
		ID: req.ID,
	}
	if err := h.serv.DeleteContest(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "delete contest "))

}

func (h *contestHander) JoinTeamsinContest(c echo.Context) error {
	var req contestJoinTeamRquest
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

	m := model.ContestsTeam{
		ContestID: req.ContestID,
		TeamID:    req.TeamID,
	}

	if err := h.serv.CreateTeamContest(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "join contest_teams"))
}
func (h *contestHander) DeleteTeamContest(c echo.Context) error {
	var req contestJoinTeamRquest
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

	m := model.ContestsTeam{
		ContestID: req.ContestID,
		TeamID:    req.TeamID,
	}

	if err := h.serv.DeleteTeamContest(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "delete contest_teams"))
}

func (h *contestHander) ListContest(c echo.Context) error {

	contests, err := h.serv.ListContest()
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, contests)
}

func (h *contestHander) ListContestByTeams(c echo.Context) error {
	sid := c.QueryParam("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}

	contests, err := h.serv.ListContestByTeams(id)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, contests)
}

func (h *contestHander) JoinContestQuestions(c echo.Context) error {
	reqs := new([]joinContestQuesiontsRequest)
	if err := c.Bind(&reqs); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	if err := c.Validate(reqs); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	qcids := []map[string]int{}
	for _, req := range *reqs {
		t := map[string]int{"qid": req.QID, "cid": req.CID}
		qcids = append(qcids, t)
	}

	if err := h.serv.JoinContestQuesionts(qcids); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "join contests_quesions"))
}

func (h *contestHander) StartContest(c echo.Context) error {
	var req startContestRequest
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

	if err := h.serv.StartContest(req.contestID); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "join contests_quesions"))
}
