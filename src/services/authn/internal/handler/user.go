package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/LainInTheWired/ctf_backend/user/model"
	"github.com/LainInTheWired/ctf_backend/user/service"
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

// 依存関係用の構造体
type UserHandler struct {
	serv service.UserService
}

func NewHanderSignup(service service.UserService) *UserHandler {
	return &UserHandler{
		serv: service,
	}
}

// json用の構造体
type SignupRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type CreateRoleRequest struct {
	Name          string `json:"name" validate:"required,max=100"`
	Namespace     string `json:"namespace" validate:"required,max=100"`
	PermissionsID []int  `json:"permissions_id"`
}
type CreatePermission struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"required,max=100" `
}
type BindRolePermissions struct {
	RoleID        int   `json:"role_id" validate:"required,max=100"`
	PermissionsID []int `json:"permissions_id" validate:"required,max=100" `
}
type BindUserRoles struct {
	UserID int   `json:"user_id" validate:"required,max=100"`
	RoleID []int `json:"role_id" validate:"required,max=100" `
}

func (h *UserHandler) Signup(c echo.Context) error {
	// リクエストから構造体にデータをコピー
	var req SignupRequest
	if err := c.Bind(&req); err != nil {
		wrappedErr := errors.Wrap(err, "request bind error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := errors.Wrap(err, "validation error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	// serviceロジックに入る
	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.serv.Signup(user); err != nil {
		wrappedErr := errors.Wrap(err, "signup error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "Service registered successfully"})
}

func (h *UserHandler) Login(c echo.Context) error {
	var req LoginRequest

	// リクエストから構造体にデータをコピー
	if err := c.Bind(&req); err != nil {
		wrappedErr := errors.Wrap(err, "request bind error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := errors.Wrap(err, "validation error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	// model Userに移し替え
	u := model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	// serviceの処理
	id, err := c.Cookie("session")
	if err == nil {
		_, err := h.serv.CheckSession(id.Value)
		if err == nil {
			return c.JSON(http.StatusAccepted, map[string]string{"message": "already login"})
		} else if errors.Is(err, redis.Nil) {

		} else {
			wrappedErr := errors.Wrap(err, "redis get error")
			fmt.Println(redis.Nil == err)
			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
		}
	}
	//
	sessionID, err := h.serv.Login(u)
	if err != nil {
		wrappedErr := errors.Wrap(err, "login error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = sessionID
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	// res := LoginResponse{
	// 	Name:  user.Name,
	// 	Email: user.Email,
	// }

	return c.JSON(http.StatusAccepted, map[string]interface{}{"message": "Login Successful"})
}

func (h *UserHandler) Logout(c echo.Context) error {
	// serviceの処理
	id, err := c.Cookie("session")
	if err != nil {
		h.serv.CheckSession(id.Value)
		return c.JSON(http.StatusAccepted, map[string]string{"message": "already login"})
	}
	return c.JSON(http.StatusAccepted, map[string]string{"message": "Successfuly Logout"})
}

func (h *UserHandler) CreateRole(c echo.Context) error {
	var req CreateRoleRequest

	// リクエストから構造体にデータをコピー
	if err := c.Bind(&req); err != nil {
		wrappedErr := errors.Wrap(err, "request bind error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := errors.Wrap(err, "validation error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	fmt.Println(req.PermissionsID[0])
	role := model.Role{
		Name:      req.Name,
		Namespace: req.Namespace,
	}
	rid, err := h.serv.AddRole(&role)
	if err != nil {
		wrappedErr := errors.Wrap(err, "add role error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	for _, pid := range req.PermissionsID {
		err = h.serv.BindRolePermissions(rid, pid)
		if err != nil {
			wrappedErr := errors.Wrap(err, "bind role permission error")
			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
		}
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "Create role"})
}
func (h *UserHandler) BindRolePermissions(c echo.Context) error {
	var req BindRolePermissions

	// リクエストから構造体にデータをコピー
	if err := c.Bind(&req); err != nil {
		wrappedErr := errors.Wrap(err, "request bind error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := errors.Wrap(err, "validation error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	for _, pid := range req.PermissionsID {
		err := h.serv.BindRolePermissions(req.RoleID, pid)
		if err != nil {
			wrappedErr := errors.Wrap(err, "bind role permission error")
			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
		}
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "success bind role permissiosn"})
}
func (h *UserHandler) BindUserRoles(c echo.Context) error {
	var req BindUserRoles

	// リクエストから構造体にデータをコピー
	if err := c.Bind(&req); err != nil {
		wrappedErr := errors.Wrap(err, "request bind error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := errors.Wrap(err, "validation error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	for _, rid := range req.RoleID {
		err := h.serv.BindUserRoles(req.UserID, rid)
		if err != nil {
			wrappedErr := errors.Wrap(err, "can't bind user role")
			log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
		}

	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "success bind user roles"})
}

// func (h *UserHandler) Bind
func (h *UserHandler) CreatePermission(c echo.Context) error {
	var req CreatePermission

	// リクエストから構造体にデータをコピー
	if err := c.Bind(&req); err != nil {
		wrappedErr := errors.Wrap(err, "request bind error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	// データをバリデーションにかける
	if err := c.Validate(req); err != nil {
		wrappedErr := errors.Wrap(err, "validation error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}

	permission := model.Permission{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := h.serv.AddPermission(&permission); err != nil {
		wrappedErr := errors.Wrap(err, "add permission error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("error:", wrappedErr)})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "Service  permission"})
}

func getRootWrappedError(err error) error {
	if err == nil {
		return nil
	}

	// エラーのリストを作成して全てのラップされたエラーを集める
	var errorsList []error
	for err != nil {
		errorsList = append(errorsList, err)
		err = errors.Unwrap(err)
	}

	// エラースタックが1つしかない場合は空文字を返す
	if len(errorsList) < 2 {
		return nil
	}

	// 下から2番目のエラーメッセージを取得
	secondLastError := errorsList[len(errorsList)-2]
	return secondLastError
}

// func rapErrorPrint(msg string, err error) error {
// 	wrappedErr := errors.Wrap(err, msg)
// 	log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
// 	return wrappedErr
// }

func getSecondLastErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	// エラーを格納するリスト
	var errorsList []error
	for err != nil {
		// 再帰的にエラーをアンラップしてリストに追加
		errorsList = append(errorsList, err)
		err = errors.Unwrap(err)
	}

	// エラースタックが2未満の場合は空文字を返す
	if len(errorsList) < 2 {
		return ""
	}

	// 2番目のエラーのみのメッセージを返す
	secondLastError := errorsList[len(errorsList)-2]

	// もしそのエラーメッセージが ":" で複数のメッセージが含まれている場合
	// 最初のメッセージだけを取り出す
	parts := strings.Split(secondLastError.Error(), ":")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[0]) // 最初の部分だけを返す
	}

	// 通常のエラーメッセージを返す
	return secondLastError.Error()
}
