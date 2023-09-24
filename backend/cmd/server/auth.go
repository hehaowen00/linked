package main

import "database/sql"

type TokenInfo struct {
	Audience      string `json:"aud"`
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
}

type UserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail string `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type StateInfo struct {
	CSRFToken   string `csv:"csrf"`
	RedirectUrl string `csv:"redirect_url"`
}

func insertUserProfile(db *sql.DB, info *UserInfo) {
}
