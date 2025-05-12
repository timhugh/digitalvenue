package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/timhugh/digitalvenue/common"
	"github.com/timhugh/digitalvenue/db/queries"
)

func New(connection string) (*sql.DB, error) {
	return sql.Open("postgres", connection)
}

type createUserDB interface {
	CreateUser(context.Context, queries.CreateUserParams) (queries.User, error)
}

func CreateUser(ctx context.Context, db createUserDB, params queries.CreateUserParams, password string) (queries.User, error) {
	passwordHash, salt, err := common.HashPassword(password, "")
	if err != nil {
		return queries.User{}, fmt.Errorf("failed to hash password: %w", err)
	}
	params.PasswordHash = sql.NullString{String: passwordHash, Valid: true}
	params.PasswordSalt = sql.NullString{String: salt, Valid: true}

	return db.CreateUser(ctx, params)
}

type authUserDB interface {
	GetUserByUsername(context.Context, string) (queries.User, error)
}

func AuthUser(ctx context.Context, db authUserDB, username string, password string) (queries.User, error) {
	user, err := db.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return queries.User{}, err
	}

	var expectedHash string
	var salt string
	if err == nil {
		expectedHash = user.PasswordHash.String
		salt = user.PasswordSalt.String
	}

	hash, _, err := common.HashPassword(password, salt)
	if hash != expectedHash {
		return queries.User{}, fmt.Errorf("failed to authenticate user")
	}

	return user, nil
}
