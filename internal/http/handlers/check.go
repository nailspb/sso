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

func Check(logger *slog.Logger, storage storage.Storage) http.HandlerFunc {
	//setup logger
	logger = slogHelper.AddOperation(logger, "http.handlers.auth.check()")
	//return function
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &responses.Auth{
			Response: responses.Response{
				Status: responses.StatusOk,
			},
		}
		log := slogHelper.AddRequestId(logger, r.Context())
		params, err := jsonHelper.Decode(&requests.Check{}, r.Body)
		if err != nil {
			resp.Response.Error = responses.ErrorBadRequest
			log.Error(resp.Response.Error, slogHelper.GetErrAttr(err))
		} else {
			if err := services.Check(params.Token, storage); err != nil {
				log.Error("Error on validate token", slogHelper.GetErrAttr(err))
				resp.Status = responses.StatusError
				resp.Error = responses.ErrorTokenNotValid
			}
			err = jsonHelper.WriteResponse(resp, w)
			if err != nil {
				log.Error(ErrorWriteResponse, slogHelper.GetErrAttr(err))
			}
		}

	}
}
