package request

type SignupRequestDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	ColorScheme string `json:"color_scheme"`
}