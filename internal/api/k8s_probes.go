package api

import "net/http"

func (api *API) readinessProbeHandle(_ http.ResponseWriter, _ *http.Request) {
	// Server ready as soon as it has started.
}

func (api *API) livenessProbeHandle(_ http.ResponseWriter, _ *http.Request) {
	// Define your liveness logic here.
}
