package responses

type User struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}
