package middleware

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"sso/pkg/helpers/slogHelper"
	"sso/pkg/http/routing"
)

func RequestId(log *slog.Logger) routing.MiddlewareFunc {
	log = slogHelper.ConfigureForMiddleware(log, "RequestId")
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			id := r.Header.Get("X-Request-Id")
			if id == "" {
				newId := uuid.New()
				id = fmt.Sprintf("%s", newId.String())
			}
			w.Header().Set("X-Request-ID", id)
			ctx = context.WithValue(ctx, "X-Request-Id", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
