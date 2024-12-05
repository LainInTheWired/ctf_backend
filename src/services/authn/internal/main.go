package main

import (
	"context"
	"log"
	"net/http"
	"os"

	myvalidator "github.com/LainInTheWired/ctf_backend/shared/pkg/validator"
	"github.com/LainInTheWired/ctf_backend/user/handler" // 正しいモジュールパス
	"github.com/LainInTheWired/ctf_backend/user/repository"
	"github.com/LainInTheWired/ctf_backend/user/service"

	"github.com/gorilla/sessions"
	"golang.org/x/xerrors"

	// 正しいモジュールパス
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log" // エイリアスを付ける
)

func main() {
	// mysql初期化処理
	db, err := repository.NewDBClient()
	if err != nil {
		xerrors.Errorf("mysql connection error: %w", err.Error())
	}
	defer db.Close()

	reddb, err := repository.NewRedis()
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
	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))
	e.Validator = myvalidator.NewValidator()

	usrep := repository.NewUserRepository(db)
	rerep := repository.NewRedisClient(reddb, context.Background())

	s := service.NewUserService(usrep, rerep)
	h := handler.NewHanderSignup(s)

	if err := s.SetInitRolePermissionsToRedis(); err != nil {
		log.Fatalf("can't role init %v", err)
	}

	e.GET("/", hello)
	e.POST("/", hello)

	e.POST("/signup", h.Signup)
	e.POST("/login", h.Login)
	e.POST("/createrole", h.CreateRole)
	e.POST("/createpermission", h.CreatePermission)
	e.POST("/bindpermission", h.BindRolePermissions)
	e.POST("/bindrole", h.BindUserRoles)

	e.DELETE("/logout", h.Logout)

	e.Start(":8000")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
