package handler

import (
	"encoding/json"
	"errors"
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
		webLog.Error(err, "Failed to decode request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
