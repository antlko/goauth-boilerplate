package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Config struct {
	JwtSecretKey         string `env:"JWT_SECRET_KEY"`
	JwtAccessTokenHours  int64  `env:"JWT_ACCESS_TOKEN_HOURS"  envDefault:"24"`
	JwtRefreshTokenHours int64  `env:"JWT_REFRESH_TOKEN_HOURS" envDefault:"168"`
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Authorizer struct {
	secretKey       []byte
	accessDuration  int64
	refreshDuration int64
}

func NewAuthorizer(config Config) Authorizer {
	return Authorizer{
		secretKey:       []byte(config.JwtSecretKey),
		accessDuration:  config.JwtAccessTokenHours,
		refreshDuration: config.JwtRefreshTokenHours,
	}
}

func (a Authorizer) CreateTokens(username string) (Tokens, error) {
	accessToken, err := a.createToken(username, a.accessDuration)
	if err != nil {
		return Tokens{}, fmt.Errorf("create access token: %w", err)
	}
	refreshToken, err := a.createToken(username, a.accessDuration)
	if err != nil {
		return Tokens{}, fmt.Errorf("create refresh token: %w", err)
	}
	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a Authorizer) ValidateAndUpdate(token string) (Tokens, error) {
	ok, identity, err := a.verifyToken(token)
	if err != nil {
		return Tokens{}, fmt.Errorf("token verification: %w")
	}
	if !ok {
		return Tokens{}, fmt.Errorf("token not verified")
	}
	tokens, err := a.CreateTokens(identity)
	if err != nil {
		return Tokens{}, fmt.Errorf("tokens not updated: %w", err)
	}
	return tokens, nil
}

func (a Authorizer) Validate(token string) (bool, string, error) {
	ok, identity, err := a.verifyToken(token)
	if err != nil {
		return false, "", fmt.Errorf("token verification: %w")
	}
	return ok, identity, nil
}

func (a Authorizer) createToken(username string, hours int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * time.Duration(hours)).Unix(),
		})

	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		return "", fmt.Errorf("token sign: %w", err)
	}
	return tokenString, nil
}

func (a Authorizer) verifyToken(tokenString string) (bool, string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return a.secretKey, nil
	})

	identityData, ok := claims["username"]
	if !ok {
		return false, "", fmt.Errorf("claim not found")
	}
	identity, ok := identityData.(string)
	if !ok {
		return false, "", fmt.Errorf("claim invalid")

	}
	if err != nil {
		return false, identity, fmt.Errorf("parse token: %w", err)
	}
	return token.Valid, identity, nil
}
