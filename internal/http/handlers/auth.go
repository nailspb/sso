package handlers

import (
	"log/slog"
	"net/http"
	"sso/internal/http/requests"
	"sso/internal/http/responses"
	"sso/internal/services"
	"sso/internal/storage"
	"sso/pkg/helpers/jsonHelper"
	"sso/pkg/helpers/slogHelper"
)

const (
	ErrorAuth          = "Auth error"
	ErrorWriteResponse = "Error on write json response"
	MsgIssuedToken     = "The user has been issued a token"
)

func Auth(logger *slog.Logger, storage storage.Storage) http.HandlerFunc {
	logger = slogHelper.AddOperation(logger, "http.handlers.auth()")
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &responses.Auth{
			Response: responses.Response{
				Status: responses.StatusError,
			},
		}
		log := slogHelper.AddRequestId(logger, r.Context())
		params, err := jsonHelper.Decode(&requests.Auth{}, r.Body)
		if err != nil {
			resp.Response.Error = responses.ErrorBadRequest
			log.Error(resp.Response.Error, slogHelper.GetErrAttr(err))
		} else if params.Login == "" || params.Password == "" {
			resp.Response.Error = responses.ErrorEmptyLoginPassword
			log.Warn(resp.Response.Error)
		} else {
			token, err := services.Auth(params.Login, params.Password, storage)
			if err != nil {
				resp.Response.Error = responses.ErrorUserNotFound
				log.Error(ErrorAuth, slogHelper.GetErrAttr(err))
			} else {
				resp.Status = responses.StatusOk
				resp.Token = token
				log.Info(MsgIssuedToken, slog.String("user_login", params.Login))
			}
		}
		err = jsonHelper.WriteResponse(resp, w)
		if err != nil {
			log.Error(ErrorWriteResponse, slogHelper.GetErrAttr(err))
		}

	}
}
