package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
)


type MalformedRequest struct {
	Status int
	msg string
}

func (mr *MalformedRequest) Error() string {
	return mr.msg

}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}, webLog logr.Logger) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.Split(ct, ";")[0])
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &MalformedRequest{Status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		// webLog.Error(err, "Failed to decode request body")
		// http.Error(w, err.Error(), http.StatusBadRequest)

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		// switch on the type of the error to extract the
		// underlying cause. We can then return a more specific error message
		switch {
		// Catch any syntax errors in the JSON and return a client error response
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			webLog.Error(err, msg)
			return &MalformedRequest{ Status: http.StatusBadRequest, msg: msg }

		// check if there is any io.ErrUnexpectedEOF error while decoding json
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			webLog.Error(err, msg)
			return &MalformedRequest{ Status: http.StatusBadRequest, msg: msg }

		// Catch any type errors,
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			webLog.Error(err, msg)
			return &MalformedRequest{ Status: http.StatusBadRequest, msg: msg }

		// Catch the error caused by extra unexpected fields in the request body
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			webLog.Error(err, msg)
			return &MalformedRequest{ Status: http.StatusBadRequest, msg: msg }

		// if request body is empty
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			webLog.Error(err, msg)
			return &MalformedRequest{ Status: http.StatusBadRequest, msg: msg }

		// return request body too large error, if request body is larger than 1MB
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			webLog.Error(err, msg)
			http.Error(w, msg, http.StatusRequestEntityTooLarge)

		// default return Internal Server Error (500)
		default:
			webLog.Error(err, "Failed to decode request body")
			return &MalformedRequest{ Status: http.StatusInternalServerError, msg: "Failed to decode request body" }

		}
	
	}

	// no error
	return nil
}