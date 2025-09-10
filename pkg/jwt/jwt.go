package jwt

import (
	"time"

	"github.com/TimNikolaev/drag-sso/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user *models.User, app *models.App, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":    user.ID,
		"email":  user.Email,
		"exp":    time.Now().Add(duration).Unix(),
		"app_id": app.ID,
	})

	tokenString, err := token.SignedString([]byte(app.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
