package requests

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Check struct {
	Token string `json:"token"`
}
