package entities

import (
	"net/mail"
	"time"
)

const INN = "7702060003"

type Info struct {
	AppConfig AppConfig   `toml:"app_config"`
	Session   SessionInfo `toml:"session_info"`
}

type AppConfig struct {
	User   UserInfo   `toml:"user_info"`
	Driver DriverInfo `toml:"driver_info"`
}

type UserInfo struct {
	Login    string `toml:"login"`
	Password string `toml:"password"`
}

type DriverInfo struct {
	Path          string        `toml:"path"`
	Com           string        `toml:"com"`
	Time          string        `toml:"time"`
	Connection    string        `toml:"baseurl"`
	PollingPeriod time.Duration `toml:"polling_period"`
	TimeoutPeriod time.Duration `toml:"timeout_duration"`
	UpdatePath    string        `toml:"update_path"`
}

type SessionInfo struct {
	AccessToken string    `json:"accessToken" toml:"access_token"`
	TokenType   string    `json:"token_type" toml:"token_type"`
	CreatedAt   time.Time `toml:"created_at"`

	UserData struct {
		ID       int    `json:"id" toml:"id"`
		FullName string `json:"fullName" toml:"full_name"`
		Username string `json:"username" toml:"user_name"`
		Avatar   string `json:"avatar" toml:"avatar"`
		Email    string `json:"email" toml:"email"`
		Role     string `json:"role" toml:"role"`

		//Ability  []struct {
		//	Subject string `json:"subject" toml:"subject"`
		//	Action  string `json:"action" toml:"action"`
		//} `json:"ability" toml:"ability"`
	} `json:"userData" toml:"user_data"`
}

func (u *UserInfo) ValidateUser() bool {
	_, err := mail.ParseAddress(u.Login)
	return err == nil && u.Password != ""
}

func (s *SessionInfo) IsDead() bool {
	return time.Since(s.CreatedAt).Hours() > 24
}
