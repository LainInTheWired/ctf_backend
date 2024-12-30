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

	"github.com/LainInTheWired/ctf-backend/pveapi/handler"
	"github.com/LainInTheWired/ctf-backend/pveapi/model"
	"github.com/LainInTheWired/ctf-backend/pveapi/repository"
	"github.com/LainInTheWired/ctf-backend/pveapi/service"
	myvalidator "github.com/LainInTheWired/ctf_backend/shared/pkg/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log" // エイリアスを付ける
	"github.com/redis/go-redis/v9"
	"golang.org/x/xerrors"
)

// 認証ヘッダーを自動的に付与するカスタムクライアントを作成
type MiddlewareTransport struct {
	Transport http.RoundTripper
	Token     string
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
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s", a.Token))
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

func NewDBClient() (*sql.DB, error) {
	db, err := sql.Open("mysql", "user:user@tcp(db:3306)/ctf")
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
	// ctx := context.Background()
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

	// redis初期化処理
	reddb, err := NewRedisClient()
	if err != nil {
		xerrors.Errorf("redis connetciono error: %w", err.Error())
	}
	defer reddb.Close()

	// proxmoxapi初期化処理
	config := &model.PVEConfig{
		APIURL:        os.Getenv("PROXMOX_API_URL"),
		Authorization: os.Getenv("PROXMOX_API_TOKEN"),
	}
	// httpclinet auth middleware
	// カスタムトランスポートを作成
	// カスタム Transport を設定（InsecureSkipVerify は本番環境では false に）
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 本番では false にする
	}

	// カスタム AuthTransport を作成
	authTransport := &MiddlewareTransport{
		Transport: tr,
		Token:     config.Authorization,
	}
	// カスタム HTTP クライアントの作成
	client := &http.Client{
		Transport: authTransport,
		Timeout:   60 * time.Second,
	}

	// 設定のバリデーション
	if config.APIURL == "" || config.Authorization == "" {
		log.Fatal("必要な環境変数が設定されていません。PROXMOX_API_URL および PROXMOX_AUTHORIZATION を設定してください。")
	}

	// echo初期化処理
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
	e.Validator = myvalidator.NewValidator()

	// p := service.NewPVEClient(config)
	r := repository.NewPVERepository(config, client)
	s := service.NewPVEService(r)
	h := handler.NewPVEAPI(s)

	e.GET("/", hello)
	e.POST("/vm", h.CreateCloudinitVM)
	e.DELETE("/vm", h.DeleteVM)
	e.POST("/cloudinit", h.Cloudinit)
	e.DELETE("/cloudinit", h.DeleteCloudinit)
	e.POST("/template", h.ToTemplate)
	e.GET("/vm/:vmid/ips", h.GetIps)
	e.GET("/cluster", h.GetClusterResource)

	// e.GET("/vm", h.GetVM)
	// e.PUT("/vm", h.GETTestHander)
	// e.POST("/getnode", h.GetNodeTestHander)
	// e.DELETE("/deletevm", h.DeleteVM)
	// e.POST("/cloneque", h.CloneQuestions)

	e.Start(":8000")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
