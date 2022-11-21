package api

import (
	"context"
)

func AuthenticateUser(ctx context.Context, loginData LoginData) (bool, error) {
	// choice, err := repository.Default.Choice.Get(ctx, expression.IDEquals(int(choiceID)))
	// if err != nil {
	// 	return false, err
	// }

	return true, nil
}
