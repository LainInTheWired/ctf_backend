package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/LainInTheWired/ctf_backend/user/handler" // 正しいモジュールパス
	"github.com/LainInTheWired/ctf_backend/user/repository"
	"github.com/LainInTheWired/ctf_backend/user/service"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"
	"golang.org/x/xerrors"

	// 正しいモジュールパス
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log" // エイリアスを付ける
)

// CustomValidator
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator
func NewValidator() echo.Validator {
	return &CustomValidator{validator: validator.New()}
}

// カスタムヴァリデータを編集
func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
			fieldName := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", fieldName))
			case "email":
				errorMessages = append(errorMessages, fmt.Sprintf("%s isn't email format.", fieldName))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s characters long.", fieldName, err.Param()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("%s is fail validation", fieldName))
			}
		}
		return fmt.Errorf(strings.Join(errorMessages, ", "))
	}
	return nil

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
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil

}

func main() {
	db, err := NewDBClient()
	if err != nil {
		xerrors.Errorf("mysql connection error: %w", err.Error())
	}

	reddb, err := NewRedisClient()
	if err != nil {
		xerrors.Errorf("redis connetciono error: %w", err.Error())
	}

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
	e.Validator = NewValidator()

	usrep := repository.NewUserRepository(db)
	rerep := repository.NewRedisClient(reddb, context.Background())

	s := service.NewUserService(usrep, rerep)
	h := handler.NewHanderSignup(s)

	e.GET("/", hello)
	e.POST("/signup", h.Signup)
	e.POST("/login", h.Login)

	e.Start(":8000")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
