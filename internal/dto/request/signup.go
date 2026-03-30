package request

type SignupRequestDto struct {
	Username string `json:"username" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	ColorScheme string `json:"color_scheme"`
}