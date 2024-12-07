package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"sigs.k8s.io/yaml"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/go-logr/logr"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/durgeshmeena/envoy-gateway-controller/internal/pkg/utils/datastore"
	"github.com/durgeshmeena/envoy-gateway-controller/internal/pkg/validation"
	"github.com/durgeshmeena/envoy-gateway-controller/internal/webserver/model"
)

type CreateClientBTPHandler struct {
	// logger
	Logger 		*logr.Logger
	// datastore
	FileStore 	*datastore.FileStore
}

func NewCreateClientBTPHandler(logger *logr.Logger, fileStore *datastore.FileStore) *CreateClientBTPHandler {
	return &CreateClientBTPHandler{
		Logger: logger,
		FileStore: fileStore,
	}
}

// // validate RateLimitType, enum with values: Global or Local
// func validateRateLimitType(rateLimitType string) error {
// 	if rateLimitType != "Global" && rateLimitType != "Local" {
// 		return errors.New("rateLimitType is not valid, rateLimitType must be either Global or Local")
// 	}
// 	return nil
// }

func validateClientBTP(clientBTPData *model.ClientBTP, webLogger logr.Logger) error {
	var errs []error
	// validate if Namespace is not empty
	if clientBTPData.Namespace == "" {
		validationError := errors.New("namespace is empty, namespace is required")
		errs = append(errs, validationError)
		webLogger.Error(validationError, "Namespace")
	}

	// validate if RateLimitHttpRoute is not empty
	if clientBTPData.RateLimitHttpRoute == "" {
		validationError := errors.New("rateLimitHttpRoute is empty, rateLimitHttpRoute is required")
		errs = append(errs, validationError)
		webLogger.Error(validationError, "RateLimitHttpRoute")
	}
	// validate RateLimitType
	// if err := validateRateLimitType(string(clientBTPData.RateLimitType)); err != nil {
	// 	errs = append(errs, err)
	// 	webLogger.Error(err, "RateLimitType")
	// }
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

	// create a new ClientBTP struct and decode the request body into it
	var clientBTPData model.ClientBTP

	// validate the request body
	err := validation.DecodeJSONBody(w, r, &clientBTPData, webLog)
	if err != nil {

		var mr *validation.MalformedRequest
		if errors.As(err, &mr) {
			webLog.Error(err, "Failed to decode request body")
			http.Error(w, mr.Error(), mr.Status)
		} else {
			webLog.Error(err, "Internal Server Error")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
	// create a new BTP struct
	// distinct name, combination of httproute, time
	name := "btp-" + clientBTPData.RateLimitHttpRoute + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	// create metadata
	metaData := model.Metadata{
		Name:               name,
		Namespace:          clientBTPData.Namespace,
		RateLimitHttpRoute: clientBTPData.RateLimitHttpRoute,
	}

	// create targetRefs for BTP using the httpRoute
	var targetRefs []gwapiv1a2.LocalPolicyTargetReferenceWithSectionName
	targetRefs = append(targetRefs, gwapiv1a2.LocalPolicyTargetReferenceWithSectionName{
		LocalPolicyTargetReference: gwapiv1a2.LocalPolicyTargetReference{
			Group: "gateway.networking.k8s.io",
			Kind:  "HTTPRoute",
			Name:  gwapiv1a2.ObjectName(clientBTPData.RateLimitHttpRoute),
		},
	})

	// ratelimit spec with global ratelimit
	rateLimitSpec := egv1a1.RateLimitSpec{
		Type: "Global",
		Global: &egv1a1.GlobalRateLimit{
			Rules: clientBTPData.RateLimitRules,
		},
	}

	// create BTP spec
	btpSpec := egv1a1.BackendTrafficPolicySpec{
		// reference to  httproute
		PolicyTargetReferences: egv1a1.PolicyTargetReferences{
			TargetRefs: targetRefs,
		},
		// ratelimit spec
		RateLimit: &rateLimitSpec,
	}

	// blank policy status
	status := gwapiv1a2.PolicyStatus{}

	// final BTP
	btp := model.BTP{
		Metadata: metaData,
		BTPSpec:  btpSpec,
		Status:   status,
	}

	btpYaml, err := yaml.Marshal(btp)
	if err != nil {
		webLog.Error(err, "Failed to marshal BTP data to YAML")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// save BTP to file
	if err := h.FileStore.SaveBTPToFile(btp); err != nil {
		webLog.Error(err, "Failed to save BTP to file")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// send the respopnse back
	w.WriteHeader(http.StatusOK)
	w.Write(btpYaml)
}
