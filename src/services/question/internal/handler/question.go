package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LainInTheWired/ctf_backend/question/model"
	"github.com/LainInTheWired/ctf_backend/question/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/xerrors"
)

type QuesionHander interface {
	CreateQuestion(c echo.Context) error
	DeleteQuestion(c echo.Context) error
	GetQuestionsInContest(c echo.Context) error
	GetQuestions(c echo.Context) error
	CloneQuestion(c echo.Context) error
}

type quesionHander struct {
	serv service.QuesionService
}
type quesionRequest struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  int      `json:"category_id"`
	Env         string   `json:"env"`
	Sshkeys     []string `json:"sshkeys"`
	Memory      int      `json"memory"`
	Password    string   `json:"password"`
	CPUs        int      `json"cpu"`
	Disk        int      `json"disk"`
	IP          string   `json:"ip,omitempty" validate:"omitempty,cidr"`
	Gateway     string   `json:"gateway,omitempty" validate:"omitempty,ip"`
	Filename    string   `json:"filename"`
}

type QuestionsInContestRequest struct {
	ContestID int `json:"contest_id"`
}

type CloneQuestionRequest struct {
	VMID        int      `json:"vmid"`
	ContestName string   `json:"contest_name"`
	Sshkeys     []string `json:"sshkeys"`
	Password    string   `json:"password"`
}

func NewQuestionHander(s service.QuesionService) QuesionHander {
	return &quesionHander{
		serv: s,
	}
}

func (h *quesionHander) CreateQuestion(c echo.Context) error {
	var req quesionRequest
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

	m := model.CreateQuestion{
		Name:        req.Name,
		Env:         req.Env,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Sshkeys:     req.Sshkeys,
		CPUs:        req.CPUs,
		Disk:        req.Disk,
		ID:          req.ID,
		Memory:      req.Memory,
		IP:          req.IP,
		Gateway:     req.Gateway,
	}

	if err := h.serv.CreateQuestion(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "delete contest_teams"))
}

func (h *quesionHander) DeleteQuestion(c echo.Context) error {
	var req quesionRequest
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

	if err := h.serv.DeleteQuestion(req.ID); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "delete contest_teams"))
}

func (h *quesionHander) GetQuestions(c echo.Context) error {
	q, err := h.serv.GetQuestions()
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, q)
}

func (h *quesionHander) GetQuestionsInContest(c echo.Context) error {
	sid := c.Param("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error: param")})
	}

	q, err := h.serv.GetQuestionsInContest(id)
	if err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, q)
}

func (h *quesionHander) QuestionHander(c echo.Context) error {

	return nil
}

func (h *quesionHander) CloneQuestion(c echo.Context) error {
	var req quesionRequest
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
	m := model.CreateQuestion{
		Name:        req.Name,
		Env:         req.Env,
		Description: req.Description,
		Sshkeys:     req.Sshkeys,
		CPUs:        req.CPUs,
		Disk:        req.Disk,
		ID:          req.ID,
		Memory:      req.Memory,
		IP:          req.IP,
		Gateway:     req.Gateway,
	}

	if err := h.serv.CloneQuestion(m); err != nil {
		wrappedErr := xerrors.Errorf(": %w", err)
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusAccepted, fmt.Sprintf("message", "delete contest_teams"))
}
