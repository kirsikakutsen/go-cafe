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

const RefreshTokenTTL = 7 * 24 * time.Hour

func generateAuthTokens(userID uint) (string, string, error) {
	if err := database.DB.
		Where("user_id = ?", userID).
		Delete(&models.RefreshToken{}).Error; err != nil {
		return "", "", err
	}

	accessToken, err := auth.SignPayLoad(userID)
	if err != nil {
		return "", "", errors.New("Failed to generate token")
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", errors.New("Failed to generate refresh token")
	}

	hashed := auth.HashToken(refreshToken)

	err = database.DB.Create(&models.RefreshToken{
		UserID:    userID,
		TokenHash: hashed,
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
	}).Error

	if err != nil {
		return "", "", errors.New("Failed to create refresh token")
	}

	return accessToken, refreshToken, nil
}

func Login(c *echo.Context) error {
	var req request.LoginRequestDto

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewBasicErrorDto(errors.New("Invalid request body")))
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

	accessToken, refreshToken, err := generateAuthTokens(acc.ID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.NewBasicErrorDto(errors.New("Failed to generate tokens")),
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
		return c.JSON(http.StatusBadRequest, response.NewBasicErrorDto(errors.New("Invalid request body")))
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


	accessToken, refreshToken, err := generateAuthTokens(acc.ID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.NewBasicErrorDto(
				errors.New("Failed to generate tokens"),
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

func Refresh(c *echo.Context) error {
	var req request.RefreshRequestDto

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewBasicErrorDto(err))
	}

	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewBasicErrorDto(err))
	}

	hashed := auth.HashToken(req.RefreshToken)

	var token models.RefreshToken

	if err := database.DB.Where("token_hash = ?", hashed).
		First(&token).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, response.NewBasicErrorDto(errors.New("Invalid refresh token")))
	}

	if token.ExpiresAt.Before(time.Now()) {
		return c.JSON(http.StatusUnauthorized, response.NewBasicErrorDto(errors.New("Refresh token expired")))
	}

	accessToken, refreshToken, err := generateAuthTokens(token.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewBasicErrorDto(errors.New("Failed to generate tokens")),)
	}

	return c.JSON(
		http.StatusOK,
		response.NewBasicSuccessDto(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}),
	)
}

func Logout(c *echo.Context) error {
	var req request.RefreshRequestDto

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewBasicErrorDto(err))
	}

	hashed := auth.HashToken(req.RefreshToken)

	if err := database.DB.Delete(
		&models.RefreshToken{},
		"token_hash = ?",
		hashed,
	).Error; err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.NewBasicErrorDto(errors.New("Failed to logout")),
		)
	}

	return c.NoContent(http.StatusOK)
}
