package auth

import (
	"context"
	"fmt"
	"os"
	"ws-trial/db"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func JwtAuth(ctx context.Context, token *string, queries *db.Queries) (bool, string, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error", err.Error())
	}
	t, err := jwt.Parse(*token, func(token *jwt.Token) (interface{}, error) {
		secret := []byte(os.Getenv("NEXTAUTH_SECRET"))
		return secret, nil
	})
	if err != nil {
		return false, "", err
	}
	if !t.Valid {
		return false, "", nil
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return false, "", nil
	}

	email, emailExists := claims["email"].(string)
	if !emailExists {
		return false, "", nil
	}

	user, err := queries.GetUserFromEmail(ctx, email)
	if err != nil {
		return false, "", nil
	}

	if user.Email == email {
		return true, email, nil
	} else {
		return false, "", nil
	}
}
