package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateJWT(tokenString, secretKey string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("неверный токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("неверные claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("отсутствует user_id в токене")
	}

	return userID, nil
}

func GenerateJWT(userID, secretKey string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
