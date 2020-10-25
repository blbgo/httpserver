package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/blbgo/general"
)

type httpServer struct {
	*http.Server
	Config
	general.Shutdowner
	closing bool
}

// New creates and starts an http server returning a general.DelayCloser that will allow clean
// shutdown
func New(config Config, handler http.Handler, shutdowner general.Shutdowner) general.DelayCloser {
	r := &httpServer{
		Server: &http.Server{
			Addr:         config.Addr(), //":1333",
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      handler,
		},
		Config:     config,
		Shutdowner: shutdowner,
	}

	go r.httpServerThread()

	return r
}

func (r *httpServer) httpServerThread() {
	var err error
	if r.HTTPS() {
		err = r.ListenAndServeTLS(r.CertFile(), r.KeyFile())
	} else {
		err = r.ListenAndServe()
	}
	if !r.closing {
		r.closing = true
		fmt.Println("httpServer.ListenAndServe returned before Close called")
		fmt.Println("must be startup error, shutting down")
		r.Shutdowner.Shutdown(err)
		return
	}
	if err == http.ErrServerClosed {
		fmt.Println("ListenAndServe returned http.ErrServerClosed")
		return
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func (r *httpServer) Close(doneChan chan<- error) {
	go r.closeServer(doneChan)
}

func (r *httpServer) closeServer(doneChan chan<- error) {
	// already closing must have been startup error
	if r.closing {
		doneChan <- nil
		return
	}
	r.closing = true
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := r.Server.Shutdown(ctx)
	if err != nil {
		fmt.Println("Error on http server shutdown: ", err)
		doneChan <- err
		return
	}
	fmt.Println("Clean http server shutdown.")
	doneChan <- nil
}
