package controller

import (
	"github.com/durgeshmeena/envoy-gateway-controller/internal/webserver/model"
	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

)


func (r *BackendTrafficPolicyReconciler) expectedBackendTrafficPolicy(btp *model.BTP) *egv1a1.BackendTrafficPolicy {
	bp := &egv1a1.BackendTrafficPolicy{
		Spec: btp.BTPSpec,
		Status: gwapiv1a2.PolicyStatus{},
	}
	bp.ObjectMeta.Name = btp.Metadata.Name
	bp.ObjectMeta.Namespace = btp.Metadata.Namespace

	return bp
}


// compare the expected BTPSpec with the actual BTPSpec
func (r *BackendTrafficPolicyReconciler) compareBackendTrafficPolicy(expectedBTPSpec *egv1a1.BackendTrafficPolicySpec, actualBTPSpec *egv1a1.BackendTrafficPolicySpec) bool {
    // compare targetRefs
	if len(expectedBTPSpec.PolicyTargetReferences.TargetRefs) != len(actualBTPSpec.PolicyTargetReferences.TargetRefs) {
		return false
	}
	for i, expectedRef := range expectedBTPSpec.PolicyTargetReferences.TargetRefs {
		actualRef := actualBTPSpec.PolicyTargetReferences.TargetRefs[i]
		if expectedRef != actualRef {
			return false
		}
	}

	// compare ratelimitspec
	if expectedBTPSpec.RateLimit != nil && actualBTPSpec.RateLimit != nil {
		if expectedBTPSpec.RateLimit.Global != actualBTPSpec.RateLimit.Global {
			return false
		}
		if expectedBTPSpec.RateLimit.Local != actualBTPSpec.RateLimit.Local {
			return false
		}
	}

	return true
	
}