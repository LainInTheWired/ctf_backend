package repository

import (
	"database/sql"

	"github.com/LainInTheWired/ctf_backend/user/model"
	"github.com/cockroachdb/errors"
	_ "github.com/go-sql-driver/mysql" // 空のインポートを追加
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(user model.User) error
	GetUserByEmail(email string) (model.User, error)
	GetUserByID(id int) (model.User, error)
	CreatePermission(permission *model.Permission) error
	BindRolePermissions(rid int, pid int) error
	BindUserRoles(uid int, rid int) error
	CreateRole(role *model.Role) (int, error)
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (u *userRepository) CreateUser(user model.User) error {
	// emailが登録されているかチェック
	if _, err := u.GetUserByEmail(user.Email); err == nil {
		return errors.Wrap(err, "already regist email")
	}
	// usersテーブルにinsertさせる
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "can't generate password hash")
		// return xerrors.Errorf(": %w", err)
	}
	ins, err := u.DB.Prepare("INSERT INTO users (name,email,password) VALUES(? ,?, ?)")
	if err != nil {
		return errors.Wrap(err, "user insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(user.Name, user.Email, hashPassword)
	if err != nil {
		return errors.Wrap(err, "can't insert user")
	}

	return nil
}

func (u *userRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	//  emailよりユーザ情報を取得
	if err := u.DB.QueryRow("SELECT id,name,email,password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, errors.Wrap(err, "not exist this email")
		}
		return model.User{}, errors.Wrap(err, "can't select user by email")

	}
	return user, nil
}
func (u *userRepository) GetUserByID(id int) (model.User, error) {
	var user model.User
	//  emailよりユーザ情報を取得
	if err := u.DB.QueryRow("SELECT id,name,email,password FROM users WHERE ID = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, errors.Wrap(err, "not exist this id")
		}
		return model.User{}, errors.Wrap(err, "can't select user by id")

	}
	return user, nil
}

func (u *userRepository) CreatePermission(permission *model.Permission) error {
	ins, err := u.DB.Prepare("INSERT INTO permissions (name,description)  VALUES(? ,?)")
	if err != nil {
		return errors.Wrap(err, "user insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(permission.Name, permission.Description)
	if err != nil {
		return errors.Wrap(err, "can't insert permissions")
	}
	return nil
}
func (u *userRepository) CreateRole(role *model.Role) (int, error) {
	ins, err := u.DB.Prepare("INSERT INTO roles (name,namespace)  VALUES(? ,?)")
	if err != nil {
		return 0, errors.Wrap(err, "user insert error")
	}
	defer ins.Close()

	r, err := ins.Exec(role.Name, role.Namespace)
	if err != nil {
		return 0, errors.Wrap(err, "can't insert roles")
	}
	lastInsertID, err := r.LastInsertId()
	if err != nil {
		// エラーハンドリング
		return 0, errors.Wrap(err, "failed to retrieve last insert ID")
	}

	return int(lastInsertID), nil
}
func (u *userRepository) BindRolePermissions(rid int, pid int) error {
	ins, err := u.DB.Prepare("INSERT INTO role_permissions ( role_id, permission_id)  VALUES(? ,?)")
	if err != nil {
		return errors.Wrap(err, "user insert error")
	}
	defer ins.Close()

	_, err = ins.Exec(rid, pid)
	if err != nil {
		return errors.Wrap(err, "can't bind role_permissions")
	}
	return nil
}
func (u *userRepository) BindUserRoles(uid int, rid int) error {
	ins, err := u.DB.Prepare("INSERT INTO user_roles (user_id,role_id) VALUES(?,?)")
	if err != nil {
		return errors.Wrap(err, "can't bind user_roles")
	}
	defer ins.Close()

	_, err = ins.Exec(uid, rid)
	if err != nil {
		return errors.Wrap(err, "can't bind user_roles")
	}
	return nil
}

func (u *userRepository) GetUserWithRolesAndPermissions(id int) (*model.User, error) {
	query := `
    SELECT
        u.id 
        u.name 
        r.id 
        r.name 
        p.id 
        p.name 
    FROM
        users AS u
        JOIN user_roles AS ur ON u.id = ur.user_id
        JOIN roles AS r ON ur.role_id = r.id
        JOIN role_permissions AS rp ON r.id = rp.role_id
        JOIN permissions AS p ON rp.permission_id = p.id
    WHERE
        u.id = $1
    ORDER BY
        r.id, p.id;
    `
	if _, err := u.GetUserByID(id); err != nil {
		return nil, errors.Wrap(err, "can't get user by id")
	}
	user := &model.User{}
	// rows, err := u.DB.Query("SELECT r.id ,r.name,r.namespace FROM roles AS r INNER JOIN user_roles ON r.id = user_roles.role_id WHERE user_roles.user_id = ?", id)
	rows, err := u.DB.Query(query, id)
	if err != nil {
		return nil, errors.Wrap(err, "can't select user_roles")
	}
	defer rows.Close()

	for rows.Next() {
		role := model.Role{}
		if err := rows.Scan(&role.ID, &role.Name, &role.Namespace); err != nil {
			return nil, errors.Wrap(err, "can't scan role")
		}
		user.Role = append(user.Role, role)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "can't scan role")
	}

	return nil, nil
}
