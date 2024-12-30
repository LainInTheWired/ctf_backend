package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LainInTheWired/ctf_backend/question/handler"
	"github.com/LainInTheWired/ctf_backend/question/repository"
	"github.com/LainInTheWired/ctf_backend/question/service"
	myvalidator "github.com/LainInTheWired/ctf_backend/shared/pkg/validator"
	"github.com/joho/godotenv"
	"golang.org/x/xerrors"

	_ "github.com/go-sql-driver/mysql" // 空のインポートを追加

	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/echo-swagger/example/docs" // docs is generated by Swag CLI, you have to import it.

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

// logBody はリクエストボディをログに出力するミドルウェアです
func logBody(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		buf, _ := io.ReadAll(c.Request().Body)
		c.Logger().Info(string(buf))
		c.Request().Body = io.NopCloser(bytes.NewBuffer(buf))
		return next(c)
	}
}
func NewDBClient() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("user:user@tcp(%s:3306)/ctf", os.Getenv("MYSQL_URL")))
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
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
	db, err := NewDBClient()
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
	e.Use(logBody)
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

	s := service.NewQuestionService(mr, pr, ter)
	h := handler.NewQuestionHander(s)

	fmt.Println(client)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/question", h.CreateQuestion)
	e.DELETE("/question", h.DeleteQuestion)

	e.GET("/question/:id", h.GetQuestionsInContest)
	e.GET("/question", h.GetQuestions)

	e.POST("/question/clone", h.CloneQuestion)
	e.DELETE("/question/clone", h.DeleteVM)

	// e.POST("/contest", h.CreateContest)
	// e.DELETE("/contest", h.DeleteContest)
	// e.POST("/team_contests", h.JoinTeamsinContest)
	// e.DELETE("/team_contests", h.DeleteTeamContest)

	e.Start(":8000")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
