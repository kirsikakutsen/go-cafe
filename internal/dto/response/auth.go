package response

type AuthResponseDto struct {
	Token string `json:"token"`
	Account AccountResponseDto `json:"account"`
}