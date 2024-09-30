package handlers

import (
	"crypto/x509"
	"encoding/pem"
	"log/slog"
	"net/http"
	"sso/internal/storage"
	slogHelper "sso/pkg/helpers/slogHelper"
)

func Key(logger *slog.Logger, storage storage.Storage) http.HandlerFunc {
	//setup logger
	logger = slogHelper.AddOperation(logger, "http.handlers.auth.key()")
	//return function
	return func(w http.ResponseWriter, r *http.Request) {
		log := slogHelper.AddRequestId(logger, r.Context())
		k, _ := storage.GetRsaKey()
		keyPem := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(&k.PublicKey),
			},
		)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="public.pem"`)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(keyPem); err != nil {
			log.Error("Error writing key in response body", slogHelper.GetErrAttr(err))
		}
	}
}
