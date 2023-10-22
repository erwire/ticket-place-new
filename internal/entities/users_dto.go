package entities

type Users struct {
	Login    string `json:"login,omitempty" db:"login"`
	Password string `json:"password,omitempty" db:"password"`
}
