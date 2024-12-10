package btpevent

import (
    "sigs.k8s.io/controller-runtime/pkg/event"
	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
)

// Channel to send BTP update events
var BTPUpdateChannel = make(chan event.TypedGenericEvent[egv1a1.BackendTrafficPolicy])      