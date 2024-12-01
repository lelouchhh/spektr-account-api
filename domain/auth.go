package domain

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"passwd"`
	Token    string `json:"token"`
}
