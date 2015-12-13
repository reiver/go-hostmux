package hostmux


import (
	"net/http"
	"strings"
)


type HostHandler interface {
	http.Handler

	// Host registers a sub-handler to be handed off to, when
	// a request comes in with one of the hosts.
	Host(subhandler http.Handler, hosts ...string) HostHandler

	// Else registers a sub-handler to be haned off to, when
	// where is no host handler to hand off to.
	Else(subhandler http.Handler)
}


type internalHostHandler struct {
	hostToHandler map[string] http.Handler
	elseHandler   http.Handler
}


// New creates a new HostHandler.
func New() HostHandler {
	hostToHandler := make(map[string] http.Handler)

	handler := internalHostHandler{
		hostToHandler:hostToHandler,
	}

	return &handler
}


func (handler *internalHostHandler) Host(subhandler http.Handler, hosts ...string) HostHandler {
	for _, host := range hosts {
		handler.hostToHandler[strings.ToLower(host)] = subhandler
	}

	return handler
}


func (handler *internalHostHandler) Else(subhandler http.Handler) {

	handler.elseHandler = subhandler
}


func (handler *internalHostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Get the host from the HTTP request.
	//
	// This will be used next to figure out which sub-handler
	// pass the HTTP request off to.
	//
	// These sub-handlers would have been previously registered
	// with the Register method.
	if nil == r {
		http.NotFound(w, r)
		return
	}

	host := r.Host
	host = strings.ToLower(host)


	// Pass of the HTTP request to the sub-handler that was
	// previously registered for this host.
	subhandler, ok := handler.hostToHandler[host]
	if !ok {
		if nil != handler.elseHandler {
			handler.elseHandler.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		} 
		return
	}

	subhandler.ServeHTTP(w, r)
}
