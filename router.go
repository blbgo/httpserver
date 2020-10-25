package httpserver

import (
	//"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Router is an interface to support defining routes
type Router interface {
	Handler(method, path string, handler http.Handler)
}

// NewHandlerAndRouter returns a http.Handler and a Router
func NewHandlerAndRouter() (http.Handler, Router) {
	r := &httprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		NotFound:               http.HandlerFunc(http.NotFound),
	}
	return r, r
}

// RouteParams is an interface to abstract out getting route paramiters from httprouter
type RouteParams interface {
	Get(req *http.Request) []string
}

type routeParams struct{}

// NewRouteParams provides a RouteParams interface
func NewRouteParams() RouteParams {
	return routeParams{}
}

func (r routeParams) Get(req *http.Request) []string {
	p := httprouter.ParamsFromContext(req.Context())
	result := make([]string, len(p))
	for i, v := range p {
		result[i] = v.Value
	}
	return result
}
