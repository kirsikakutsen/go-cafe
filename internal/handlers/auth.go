package handlers

import (
	"errors"
	"go-cafe/internal/auth"
	"go-cafe/internal/database"
	"go-cafe/internal/dto/request"
	"go-cafe/internal/dto/response"
	"go-cafe/internal/models"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func Login(c *echo.Context) error {
	var req request.AuthRequestDto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewBasicErrorDto(err))
	}

	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewBasicErrorDto(err))
	}

	var acc models.Account
	if err := database.DB.Where("username = ?", req.Username).First(&acc).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, response.NewBasicErrorDto(errors.New("No such account")))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, response.NewBasicErrorDto(errors.New("Wrong credentials")))
	}

	token, err := auth.SignPayLoad(acc.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewBasicErrorDto(errors.New("Failed to generate token")))
	}

	resp := response.AuthResponseDto{
		Token: token,
		Account: response.NewAccountDto(acc),
	}

	return c.JSON(http.StatusOK, response.NewBasicSuccessDto(resp))
}
