package main

type Profile struct {
	UserId       string `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	DisplayName  string `json:"display_name"`
	EmailAddress string `json:"email_address"`
}
