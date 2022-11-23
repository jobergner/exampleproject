package api

import (
	"context"
	"errors"
	"exampleproject/db"
	"exampleproject/entity"
	"exampleproject/log"
	"exampleproject/repository"
	"exampleproject/repository/expression"

	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(ctx context.Context, loginData LoginData) (bool, error) {
	user, err := repository.Default.User.Get(ctx, expression.UserNameEquals(loginData.Name), expression.IsNotArchived())
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(loginData.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		log.Log(log.CompareHash, log.Err(err))
	}

	return true, nil
}

func CreateUser(ctx context.Context, loginData LoginData) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(loginData.Password), bcrypt.MinCost)
	if err != nil {
		log.Log(log.CreateHash, log.Err(err), log.Password(loginData.Password))
		return err
	}

	user := &entity.User{
		Name:         loginData.Name,
		PasswordHash: passwordHash,
	}

	if _, err := repository.Default.User.Create(ctx, db.NoTSX, user); err != nil {
		return err
	}

	return nil
}
