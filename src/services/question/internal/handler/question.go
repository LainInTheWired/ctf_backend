package handler

import (
	"fmt"
	"net/http"

	"github.com/LainInTheWired/ctf_backend/question/model"
	"github.com/LainInTheWired/ctf_backend/question/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/xerrors"
)

type QuesionHander interface {
	CreateQuestion(c echo.Context) error
	DeleteQuestion(c echo.Context) error
}

type quesionHander struct {
	serv service.QuesionService
}

type quesionRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Env  string `json:"env"`
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

	m := model.Question{
		Name: req.Name,
		Env:  req.Env,
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
