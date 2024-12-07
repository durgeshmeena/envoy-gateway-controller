package model

import (
	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// ClientBTP struct, to validate the BTP data received from the client
type ClientBTP struct {
	// ClientBTP is the struct for the client BTP, which is used to
	// store the BTP data received from the client.
	Namespace          string `json:"namespace"`
	RateLimitHttpRoute string `json:"rateLimitHttpRoute"`
	// RateLimitType      egv1a1.RateLimitType   `json:"rateLimitType"`
	RateLimitRules []egv1a1.RateLimitRule `json:"rateLimitRules"`
}

// metadata for the BTP
type Metadata struct {
	Name               string `json:"name"`
	Namespace          string `json:"namespace"`
	RateLimitHttpRoute string `json:"rateLimitHttpRoute"`
}

// BTP struct, which will be stored in the json file
type BTP struct {
	// metadata
	Metadata Metadata `json:"metadata"`
	// btp spec
	BTPSpec egv1a1.BackendTrafficPolicySpec `json:"btpSpec"`
	// status
	Status gwapiv1a2.PolicyStatus `json:"status"`
}
