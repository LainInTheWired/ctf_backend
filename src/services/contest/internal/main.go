package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LainInTheWired/ctf_backend/contest/hander"
	"github.com/LainInTheWired/ctf_backend/contest/repository"
	"github.com/LainInTheWired/ctf_backend/contest/service"
	myvalidator "github.com/LainInTheWired/ctf_backend/shared/pkg/validator"
	"github.com/joho/godotenv"
	"golang.org/x/xerrors"

	_ "github.com/go-sql-driver/mysql" // 空のインポートを追加

	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"

	// 正しいモジュールパス
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log" // エイリアスを付ける
)

// 認証ヘッダーを自動的に付与するカスタムクライアントを作成
type MiddlewareTransport struct {
	Transport http.RoundTripper
}

type httpClientRequestLog struct {
	method  string
	url     string
	headers map[string][]string
	body    string
}
type httpClinetResponseLog struct {
	status string
	header http.Header
	body   string
}

// リクエストのたびに Authorization ヘッダーを自動で追加
func (a *MiddlewareTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var reqbody []byte
	if req.Body != nil {
		reqbody, err := io.ReadAll(req.Body)
		if err != nil {
			xerrors.Errorf("can't read reqest body : %w", err)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(reqbody))

	}

	reqlog := httpClientRequestLog{
		method:  req.Method,
		url:     req.URL.String(),
		headers: req.Header,
		body:    string(reqbody),
	}

	echoLog.Debug(reqlog)
	res, err := a.Transport.RoundTrip(req)
	if err != nil {
		xerrors.Errorf("http clinet send error : %w", err)
	}

	if res != nil && res.Body != nil {

		resbody, err := io.ReadAll(res.Body)
		if err != nil {
			xerrors.Errorf("can't read response body : %w", err)
		}
		res.Body = io.NopCloser(bytes.NewBuffer(resbody))
		reslog := httpClinetResponseLog{
			status: res.Status,
			header: res.Header,
			body:   string(resbody),
		}
		echoLog.Debug(reslog)
	}
	return res, err
}

func NewRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "user",
		DB:       0,
	})
	// 接続確認
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func main() {
	// .envファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// mysql初期化処理
	db, err := repository.NewDBClient()
	if err != nil {
		xerrors.Errorf("mysql connection error: %w", err.Error())
	}
	defer db.Close()

	reddb, err := NewRedisClient()
	if err != nil {
		xerrors.Errorf("redis connetciono error: %w", err.Error())
	}
	defer reddb.Close()

	e := echo.New()
	// セッション
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	// CORSミドルウェアの設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"}, // フロントエンドのURL
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		ExposeHeaders:    []string{echo.HeaderContentLength},
		AllowCredentials: true, // クッキーや認証情報を許可する場合
	}))

	// ログの表示方式の切り替え
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	// ログフォーマットの設定
	if env == "production" {
		// 本番環境ではJSON形式のログを出力
		echoLog.SetHeader(`{"time":"${time_rfc3339}","level":"${level}","file":"${short_file}","line":"${line}","message":"${message}"}`)
		echoLog.SetLevel(echoLog.ERROR) // 必要に応じてログレベルを調整
	} else {
		// 開発環境ではテキスト形式のログを出力
		echoLog.SetHeader("${time_rfc3339} [${level}] ${short_file}:${line} ${message}")
		echoLog.SetLevel(echoLog.DEBUG) // 開発中は詳細なログを出力
	}

	// アクセスログを出力
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))
	e.Validator = myvalidator.NewValidator()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 本番では false にする
	}
	// カスタム AuthTransport を作成
	authTransport := &MiddlewareTransport{
		Transport: tr,
	}
	// カスタム HTTP クライアントの作成
	client := &http.Client{
		Transport: authTransport,
		Timeout:   60 * time.Second,
	}

	mr := repository.NewMysqlRepository(db)
	pr := repository.NewPVEAPIRepository(client, os.Getenv("PVEAPI_URL"))
	ter := repository.NewTeamRepository(client, os.Getenv("TEAM_URL"))
	qr := repository.NewQuestionRepository(client, os.Getenv("QUESTION_URL"))

	s := service.NewContestService(pr, mr, ter, qr)
	h := hander.NewContestHander(s)

	fmt.Println(h)
	e.POST("/contest", h.CreateContest)
	e.DELETE("/contest/:contestID", h.DeleteContest)
	e.POST("/contest/:contestID/team/:teamID", h.JoinTeamsinContest)
	e.POST("/contest/:contestID/team", h.JoinTeamsinContest)
	// e.DELETE("/team_contests", h.DeleteTeamContest)
	e.DELETE("/contest/:contestID/team/:teamID", h.DeleteTeamContest)

	e.GET("/contest", h.ListContest)
	e.GET("/contest/:contestID/team", h.ListContestByTeams)
	e.GET("/contest/:contestID/point", h.GetPoints)
	// e.POST("/start", h.StartContest)
	e.POST("/contest/:contestID/start", h.StartContest)
	e.POST("/contest/:contestID/stop", h.StopContest)
	// e.DELETE("/contest/:contestID/vm")

	e.POST("/contest/:contestID/answer", h.CheckAnswer)
	e.GET("/contest/:contestID", h.ListQuestionsByContestID)
	e.POST("/contest/:contestID/question", h.JoinContestQuestions)
	e.PUT("/contest/:contestID/question/:questionID", h.UpdateContestQuestions)
	e.GET("/contest/:contestID/cloudinit/:questionID", h.GetCloudinit)
	e.GET("/contest/cluster", h.GetClusterResource)

	e.DELETE("/contest/vm", h.AllVMDelete)
	// e.PUT("/contest/:contestID/question/:questionID",h.)

	e.Start(":8000")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
