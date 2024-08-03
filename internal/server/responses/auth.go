package responses

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"` // Optional (from Auth0 doc refreshing access/refresh tokens)
}
