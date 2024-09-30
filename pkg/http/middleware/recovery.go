package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"sso/pkg/helpers/slogHelper"
	"sso/pkg/http/routing"
)

func Recovery(log *slog.Logger, debugLevel string) routing.MiddlewareFunc {
	log = slogHelper.ConfigureForMiddleware(log, "Recovery")
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(http.StatusText(http.StatusInternalServerError), slogHelper.GetErrAttr(err.(error)), slog.String("stack", string(debug.Stack())))
					if debugLevel == "prod" {
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					} else {
						http.Error(w, http.StatusText(http.StatusInternalServerError)+"\r\n\r\n"+string(debug.Stack()), http.StatusInternalServerError)
					}

				}
			}()
			next.ServeHTTP(w, req)
		}
	}
}
