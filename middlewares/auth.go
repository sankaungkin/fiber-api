package middlewares

import (
	jwtware "github.com/gofiber/contrib/jwt"
)

const jwtSecret = "superSecretKey"

func NewAuthMiddleware() {
	jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(jwtSecret)},
	})

}
