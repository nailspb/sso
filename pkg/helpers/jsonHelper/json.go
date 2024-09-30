package jsonHelper

import (
	"encoding/json"
	"io"
	"net/http"
	"sso/pkg/helpers/errorHelper"
)

const (
	ErrorJsonDecode    = "json decode error"
	ErrorJsonEncode    = "json encode error"
	ErrorWriteResponse = "error write response"
)

func Decode[T any](s *T, body io.ReadCloser) (*T, error) {
	const op = "pkg.helper.jsonHelper.decode()"
	decoder := json.NewDecoder(body)
	err := decoder.Decode(s)
	if err != nil {
		return nil, errorHelper.WrapError(op, ErrorJsonDecode, err)
	}
	return s, nil
}

func Encode(payload any) ([]byte, error) {
	const op = "pkg.helper.jsonHelper.encode()"
	marshal, err := json.Marshal(payload)
	if err != nil {
		return nil, errorHelper.WrapError(op, ErrorJsonEncode, err)
	}
	return marshal, nil
}

func WriteResponse(resp any, w http.ResponseWriter) error {
	const op = "pkg.helper.jsonHelper.WriteJsonResponse()"
	if encode, err := Encode(resp); err != nil {
		return errorHelper.WrapError(op, ErrorJsonEncode, err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(encode); err != nil {
			return errorHelper.WrapError(op, ErrorWriteResponse, err)
		}
	}
	return nil
}
