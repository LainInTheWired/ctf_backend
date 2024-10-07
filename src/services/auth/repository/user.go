package repository

import (
	"database/sql"

	"github.com/LainInTheWired/ctf_backend/user/model"
	_ "github.com/go-sql-driver/mysql" // 空のインポートを追加
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
)

type UserRepository interface {
	CreateUser(user model.User) error
	GetUserByEmail(email string) (model.User, error)
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (u userRepository) CreateUser(user model.User) error {
	// emailが登録されているかチェック
	if _, err := u.GetUserByEmail(user.Email); err == nil {
		return xerrors.Errorf(": %w", err)
	}
	// usersテーブルにinsertさせる
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	ins, err := u.DB.Prepare("INSERT INTO users (name,email,password) VALUES(? ,?, ?)")
	if err != nil {
		return err
	}
	defer ins.Close()

	_, err = ins.Exec(user.Name, user.Email, hashPassword)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

func (u userRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	//  emailよりユーザ情報を取得
	if err := u.DB.QueryRow("SELECT id,name,email,password FROM users WHERE email = ?", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, xerrors.Errorf(": %w", err)
		}
		return model.User{}, xerrors.Errorf(": %w", err)

	}

	return user, nil
}
