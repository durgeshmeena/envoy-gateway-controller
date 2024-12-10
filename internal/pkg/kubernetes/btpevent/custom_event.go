package btpevent

import (
    "sigs.k8s.io/controller-runtime/pkg/event"
)

var BTPUpdateChannel = make(chan event.GenericEvent)