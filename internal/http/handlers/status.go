package handlers

import (
	"log/slog"
	"net/http"
	"sso/internal/http/responses"
	jsonHelper "sso/pkg/helpers/jsonHelper"
	slogHelper "sso/pkg/helpers/slogHelper"
)

func Status(logger *slog.Logger) http.HandlerFunc {
	//setup logger
	logger = slogHelper.AddOperation(logger, "http.handlers.test.status()")
	//return function
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &responses.Response{
			Status: responses.StatusOk,
		}
		log := slogHelper.AddRequestId(logger, r.Context())
		err := jsonHelper.WriteResponse(resp, w)
		if err != nil {
			log.Error(ErrorWriteResponse, slogHelper.GetErrAttr(err))
		}
	}
}
