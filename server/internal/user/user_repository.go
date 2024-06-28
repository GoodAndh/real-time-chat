package user

import (
	"context"
	"database/sql"
	"fmt"
	"realTime/server/utils"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUsers(ctx context.Context, user *User) (int, error) {

	if _, err := r.GetUserByUsername(ctx, user.Username); err == nil {
		return 0, fmt.Errorf("username telah digunakan")
	}

	if _, err := r.GetUserByEmail(ctx, user.Email); err == nil {
		return 0, fmt.Errorf("email telah digunakan")
	}

	query := "insert into users(username,email,password,created_at,last_updated) values(?,?,?,?,?)"
	result, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.CreatedAt, user.LastUpdated)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	err := r.db.QueryRowContext(ctx, "select id,username,email,password from users where username = ?", username).Scan(&user.ID, &user.Username, &user.Email,&user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := r.db.QueryRowContext(ctx, "select id,username,email,password from users where email = ?", email).Scan(&user.ID, &user.Username, &user.Email,&user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}
