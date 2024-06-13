package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"sso/internal/domain/models"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID
	tokenString, err := token.SignedString([]byte(app.Secret)) /// !!! в будущем секрет не надо передавать

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

///!!!Написать на эту функцию тесты
