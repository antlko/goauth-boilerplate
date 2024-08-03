package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id       int64  `db:"id"`
	Login    string `db:"login"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return UserRepo{db: db}
}

func (u UserRepo) GetById(ctx context.Context, id int64) (User, error) {
	var user User
	if err := u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id); err != nil {
		return User{}, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (u UserRepo) GetByLogin(ctx context.Context, login string) (User, error) {
	var user User
	if err := u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE login = $1", login); err != nil {
		return User{}, fmt.Errorf("get user by login: %w", err)
	}
	return user, nil
}

func (u UserRepo) GetByLoginOrEmail(ctx context.Context, login, email string) (User, error) {
	var user User
	if err := u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE login = $1 OR email = $2", login, email); err != nil {
		return User{}, fmt.Errorf("get user by login: %w", err)
	}
	return user, nil
}

func (u UserRepo) Insert(ctx context.Context, user User) error {
	_, err := u.db.NamedExecContext(ctx, "INSERT INTO users (login, email, password) VALUES (:login, :email, :password);", user)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}
