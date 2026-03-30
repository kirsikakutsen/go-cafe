package response

import (
	"go-cafe/internal/models"
	"time"
)

type AccountResponseDto struct {
	ID uint `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	ColorScheme string `json:"color_scheme"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccountDto(acc models.Account) AccountResponseDto {
	return AccountResponseDto{
		ID: acc.ID,
		Username: acc.Username,
		Email: acc.Email,
		ColorScheme: acc.ColorScheme,
		CreatedAt: acc.CreatedAt,
	}
}