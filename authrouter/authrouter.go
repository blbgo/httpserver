package authrouter

import (
	//"errors"
	"net/http"

	"github.com/blbgo/httpserver"
)

type CheckAuth interface {
	HasAuth(req *http.Request, required string) bool
}

// AuthRouter is an interface to support defining routes that require auth
type AuthRouter interface {
	AuthHandler(method, path, required string, handler http.Handler)
}

type authRouter struct {
	httpserver.Router
	CheckAuth
}

// NewAuthRouter returns a AuthRouter
func NewAuthRouter(router httpserver.Router, checkAuth CheckAuth) AuthRouter {
	return &authRouter{
		Router:    router,
		CheckAuth: checkAuth,
	}
}

func (r *authRouter) AuthHandler(method, path, required string, handler http.Handler) {
	r.Handler(
		method,
		path,
		http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				if r.HasAuth(req, required) {
					handler.ServeHTTP(rw, req)
				} else {
					http.NotFound(rw, req)
				}
			},
		),
	)
}
