package auth

import (
	"errors"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
)

type registerRequest struct {
	Email string `json:"email"`
	First string `json:"first"`
	Last  string `json:"last"`
}

func (req *registerRequest) isValid() error {
	req.Email = strings.TrimSpace(req.Email)
	req.First = strings.TrimSpace(req.First)
	req.Last = strings.TrimSpace(req.Last)

	if req.Email == "" {
		return errors.New("missing email")
	}

	if req.First == "" {
		return errors.New("missing first name")
	}

	if req.Last == "" {
		return errors.New("missing last name")
	}

	return nil
}

type loginRequest struct {
	Email    string `json:"email"`
	PassCode string `json:"passCode"`
}

func (req *loginRequest) isValid() error {
	req.Email = strings.TrimSpace(req.Email)
	req.PassCode = strings.TrimSpace(req.PassCode)

	if req.Email == "" {
		return errors.New("missing email")
	}

	if req.PassCode == "" {
		return errors.New("missing passcode")
	}

	if len(req.PassCode) != 6 {
		return errors.New("invalid passcode")
	}

	return nil
}

type authClaims struct {
	jwt.RegisteredClaims
}
type user struct {
	id     string
	secret string
}
