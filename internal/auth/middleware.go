package auth

import (
	"errors"
	"go-cafe/internal/dto/response"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func VerifyAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		authHeader := (*c).Request().Header.Get("Authorization")

		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, response.NewBasicErrorDto(errors.New("Missing authorization header")))
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, response.NewBasicErrorDto(errors.New("Invalid authorization header")))
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := VerifyPayload(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, response.NewBasicErrorDto(errors.New("Invalid token")))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.JSON(http.StatusUnauthorized, response.NewBasicErrorDto(errors.New("Invalid token claims")))
		}

		accID, ok := claims["accID"].(float64)
		if !ok {
			return c.JSON(http.StatusUnauthorized, response.NewBasicErrorDto(errors.New("INvalid account ID")))
		}

		(c).Set("accID", uint(accID))

		return next(c)
	}
}