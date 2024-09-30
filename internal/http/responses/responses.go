package responses

const (
	StatusOk    = "ok"
	StatusError = "error"
)

const (
	ErrorBadRequest         = "bad request"
	ErrorEmptyLoginPassword = "empty login or password"
	ErrorUserNotFound       = "user not found"
	ErrorTokenNotValid      = "token signature is invalid"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type Auth struct {
	Response
	Token string `json:"token,omitempty"`
}
