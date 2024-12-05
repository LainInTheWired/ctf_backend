package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type AuthzHander interface {
}

type authzHander struct {
}

func NewHanderSignup(service service.UserService) *UserHandler {
	return &UserHandler{
		serv: service,
	}
}

func (h *AuthzHander) CreateRole(c echo.Context) error {
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
func (h *AuthzHander) BindRolePermissions(c echo.Context) error {
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
func (h *AuthzHander) BindUserRoles(c echo.Context) error {
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

// func (h *AuthzHander) Bind
func (h *AuthzHander) CreatePermission(c echo.Context) error {
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
