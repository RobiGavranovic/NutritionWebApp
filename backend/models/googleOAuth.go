package models

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Authuser    string `json:"authuser"`
	Expires_in  int    `json:"expires_in"`
	Prompt      string `json:"prompt"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
