package request

type RefreshRequestDto struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}