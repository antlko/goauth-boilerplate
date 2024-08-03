package requests

type SignUpRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type VerifyAndRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
