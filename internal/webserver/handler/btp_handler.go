package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"sigs.k8s.io/yaml"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/go-logr/logr"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/durgeshmeena/envoy-gateway-controller/internal/pkg/validation"
)

type ClientBTP struct {
	// ClientBTP is the struct for the client BTP, which is used to
	// store the BTP data received from the client.
	RateLimitHttpRoute string                 `json:"rateLimitHttpRoute"`
	RateLimitType      egv1a1.RateLimitType   `json:"rateLimitType"`
	RateLimitRules     []egv1a1.RateLimitRule `json:"rateLimitRules"`
}

type CreateClientBTPHandler struct {
	// logger
	Logger *logr.Logger
}

// validate RateLimitType, enum with values: Global or Local
func validateRateLimitType(rateLimitType string) error {
	if rateLimitType != "Global" && rateLimitType != "Local" {
		return errors.New("rateLimitType is not valid, rateLimitType must be either Global or Local")
	}
	return nil
}

func validateClientBTP(clientBTPData *ClientBTP, webLogger logr.Logger) error {
	var errs []error
	// validate if RateLimitHttpRoute is not empty
	if clientBTPData.RateLimitHttpRoute == "" {
		validationError := errors.New("rateLimitHttpRoute is empty, rateLimitHttpRoute is required")
		errs = append(errs, validationError)
		webLogger.Error(validationError, "RateLimitHttpRoute")
	}
	// validate RateLimitType
	if err := validateRateLimitType(string(clientBTPData.RateLimitType)); err != nil {
		errs = append(errs, err)
		webLogger.Error(err, "RateLimitType")
	}
	// validate RateLimitRules is not empty
	if len(clientBTPData.RateLimitRules) == 0 {
		validationError := errors.New("rateLimitRules is empty, rateLimitRules is required")
		errs = append(errs, validationError)
		webLogger.Error(validationError, "RateLimitRules")
	}
	// validate RateLimitRules
	for i, rule := range clientBTPData.RateLimitRules {
		if err := validation.ValidateRateLimitRule(&rule); err != nil {
			errs = append(errs, err)
			rateLimitRuleNo := i + 1
			rateLimitRuleName := "RateLimitRule " + strconv.Itoa(rateLimitRuleNo)
			webLogger.Error(err, rateLimitRuleName)
		}
	}
	return utilerrors.NewAggregate(errs)
}

// create ServeHTTP method for CreateClientBTPHandler
func (h *CreateClientBTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// dereference the logger
	webLog := *h.Logger

	webLog.Info("Received POST request", "time", time.Now())

	// inforce maximum request size of 1MB
	// A request body larger than that will now result in
	// Decode() returning a "http: request body too large" error.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// disallow unknown fields
	// this will cause Decode() to return an error if the JSON
	// contains any keys which do not match the ClientBTP struct.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// create a new ClientBTP struct and decode the request body into it
	var clientBTPData ClientBTP
	err := dec.Decode(&clientBTPData)
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
			http.Error(w, msg, http.StatusBadRequest)

		// check if there is any io.ErrUnexpectedEOF error while decoding json
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			webLog.Error(err, msg)
			http.Error(w, msg, http.StatusBadRequest)

		// Catch any type errors,
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			webLog.Error(err, msg)
			http.Error(w, msg, http.StatusBadRequest)

		// Catch the error caused by extra unexpected fields in the request body
		// case strings.HasPrefix(err.Error(), "json: unknown field "):
		// 	fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		// 	msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		// 	webLog.Error(err, msg)
		// 	http.Error(w, msg, http.StatusBadRequest)

		// if request body is empty
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			webLog.Error(err, msg)
			http.Error(w, msg, http.StatusBadRequest)

		// return request body too large error, if request body is larger than 1MB
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			webLog.Error(err, msg)
			http.Error(w, msg, http.StatusRequestEntityTooLarge)

		// default return Internal Server Error (500)
		default:
			webLog.Error(err, "Failed to decode request body")
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}
		return
	}

	webLog.Info("valid request body", "data", clientBTPData)
	// validate the ClientBTP data
	if err := validateClientBTP(&clientBTPData, webLog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		webLog.Error(err, "Failed to validate ClientBTP data")
		return
	}

	// convert the ClientBTP struct to YAML
	clientBTPYamlData, err := yaml.Marshal(clientBTPData)
	if err != nil {
		webLog.Error(err, "Failed to marshal ClientBTP data to YAML")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log the ClientBTP data
	webLog.Info("ClientBTP data", "data", string(clientBTPYamlData))

	// send yaml response
	w.Header().Set("Content-Type", "application/x-yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(clientBTPYamlData)

}
