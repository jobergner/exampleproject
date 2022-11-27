package api

import (
	"context"
	"errors"
	"exampleproject/entity"
	"exampleproject/entity/selector/user"
	"exampleproject/log"
	"exampleproject/repository"

	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(ctx context.Context, loginData LoginData) (bool, error) {
	var u entity.User
	if err := repository.Default.User.Get(ctx, user.NameEquals(loginData.Name), user.IsNotArchived()); err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(loginData.Password)); err != nil {
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

	u := &entity.User{
		Name:         loginData.Name,
		PasswordHash: passwordHash,
	}

	if _, err := repository.Default.User.Create(ctx, u); err != nil {
		return err
	}

	return nil
}
