package sessions

import (
	"errors"
	"net/http"
	"sync"

	"github.com/segmentio/ksuid"
)

// Sessions is an interface to manage sessions for http requests
type Sessions interface {
	Start(w http.ResponseWriter) (map[string][]byte, error)
	Get(req *http.Request) (map[string][]byte, error)
	End(w http.ResponseWriter)
	GetOrStart(w http.ResponseWriter, req *http.Request) (map[string][]byte, error)
}

// NoSession is an error indicating no session was found
var NoSession = errors.New("no session")

type sessions struct {
	name     string
	maxAge   int
	lock     sync.Mutex
	sessions map[string]map[string][]byte
}

// NewSessions provides a Sessions implementation
func NewSessions(c Config) Sessions {
	return &sessions{
		name:     c.Name(),
		maxAge:   c.MaxAge(),
		sessions: make(map[string]map[string][]byte),
	}
}

func (r *sessions) Start(w http.ResponseWriter) (map[string][]byte, error) {
	id, err := ksuid.NewRandom()
	if err != nil {
		return nil, err
	}
	sid := id.String()
	ses := make(map[string][]byte)
	r.lock.Lock()
	r.sessions[sid] = ses
	r.lock.Unlock()
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     r.name,
			Value:    sid,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   r.maxAge,
		},
	)
	return ses, nil
}

func (r *sessions) Get(req *http.Request) (map[string][]byte, error) {
	cookie, err := req.Cookie(r.name)
	if err != nil {
		return nil, err
	}
	if cookie.Value == "" {
		return nil, NoSession
	}
	r.lock.Lock()
	ses, ok := r.sessions[cookie.Value]
	r.lock.Unlock()
	if !ok {
		return nil, NoSession
	}
	return ses, nil
}

func (r *sessions) End(w http.ResponseWriter) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     r.name,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		},
	)
}

func (r *sessions) GetOrStart(
	w http.ResponseWriter,
	req *http.Request,
) (map[string][]byte, error) {
	ses, err := r.Get(req)
	if err == nil {
		return ses, nil
	}
	return r.Start(w)
}
