package middlewares

import (
	"errors"

	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)


func RequireJwtAuth(c *fiber.Ctx) error {
	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == ""{
		return errors.New("Not found jwt Secret")
	}

	reqtoken := c.GetReqHeaders()["Token"]
	if reqtoken == nil {
		return errors.New("Not found Token Header")
	}

	token, err := jwt.Parse(reqtoken[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return []byte(jwt_secret), nil
	})

	if err != nil {
		return err
	}
	if token == nil {
		return errors.New("invalid token")
	}

	c.Locals("user", token.Claims)

	return nil
}
