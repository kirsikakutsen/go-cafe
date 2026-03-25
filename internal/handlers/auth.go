package handlers

import (
	"errors"
	"go-cafe/internal/auth"
	"go-cafe/internal/database"
	"go-cafe/internal/dto/request"
	"go-cafe/internal/dto/response"
	"go-cafe/internal/models"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func Login(c *echo.Context) error {
	var req request.LoginRequestDto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewBasicErrorDto(err))
	}

	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewBasicErrorDto(err))
	}

	var acc models.Account
	if err := database.DB.Where("email = ?", req.Email).First(&acc).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, response.NewBasicErrorDto(errors.New("Invalid credentials")))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, response.NewBasicErrorDto(errors.New("Invalid credentials")))
	}

	accessToken, err := auth.SignPayLoad(acc.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewBasicErrorDto(errors.New("Failed to generate token")))
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.NewBasicErrorDto(
				errors.New("Failed to generate refresh token"),
			),
		)
	}

	hashed := auth.HashToken(refreshToken)

	if err := database.DB.Create(&models.RefreshToken{
		UserID:    acc.ID,
		TokenHash: hashed,
		ExpiresAt: time.Now().
			Add(7 * 24 * time.Hour),
	}).Error; err != nil {

		return c.JSON(
			http.StatusInternalServerError,
			response.NewBasicErrorDto(
				errors.New("Failed to create refresh token"),
			),
		)
	}

	resp := response.AuthResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Account:      response.NewAccountDto(acc),
	}

	return c.JSON(http.StatusOK, response.NewBasicSuccessDto(resp))
}

func Signup(c *echo.Context) error {
	var req request.SignupRequestDto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewBasicErrorDto(err))
	}

	if len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, response.NewBasicErrorDto(errors.New("Password must be at least 8 characters long")))
	}

	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewBasicErrorDto(err))
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewBasicErrorDto(errors.New("Failed to create account")))
	}

	acc := models.Account{
		Username:    req.Username,
		Email:       req.Email,
		Password:    string(hashedPass),
		ColorScheme: req.ColorScheme,
	}

	if err := database.DB.Create(&acc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewBasicErrorDto(errors.New("Failed to create account")))
	}

	accessToken, err := auth.SignPayLoad(acc.ID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.NewBasicErrorDto(
				errors.New("Failed to generate token"),
			),
		)
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.NewBasicErrorDto(
				errors.New("Failed to generate refresh token"),
			),
		)
	}

	hashed := auth.HashToken(refreshToken)

	if err := database.DB.Create(&models.RefreshToken{
		UserID:    acc.ID,
		TokenHash: hashed,
		ExpiresAt: time.Now().
			Add(7 * 24 * time.Hour),
	}).Error; err != nil {

		return c.JSON(
			http.StatusInternalServerError,
			response.NewBasicErrorDto(
				errors.New("Failed to create refresh token"),
			),
		)
	}

	resp := response.AuthResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Account:      response.NewAccountDto(acc),
	}

	return c.JSON(http.StatusCreated, response.NewBasicSuccessDto(resp))
}
