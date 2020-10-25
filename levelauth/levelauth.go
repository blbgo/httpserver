package levelauth

import (
	"encoding/binary"
	"net/http"
	"strconv"

	"github.com/blbgo/httpserver/authrouter"
	"github.com/blbgo/httpserver/sessions"
)

// LevelAuth is an interface to support setting the auth level in a session
type LevelAuth interface {
	SetAuthLevel(w http.ResponseWriter, req *http.Request, level uint32) error
}

type levelAuth struct {
	sessions.Sessions
}

// NewLevelAuth returns a LevelAuth and a authrouter.CheckAuth
func NewLevelAuth(ses sessions.Sessions) (LevelAuth, authrouter.CheckAuth) {
	r := &levelAuth{Sessions: ses}
	return r, r
}

func (r *levelAuth) SetAuthLevel(w http.ResponseWriter, req *http.Request, level uint32) error {
	ses, err := r.Sessions.GetOrStart(w, req)
	if err != nil {
		return err
	}
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, level)
	ses["LevelAuth"] = buf
	return nil
}

func (r *levelAuth) HasAuth(req *http.Request, required string) bool {
	requiredLevel, err := strconv.ParseUint(required, 10, 32)
	if err != nil {
		return false
	}
	ses, err := r.Sessions.Get(req)
	if err != nil {
		return false
	}
	buf, ok := ses["LevelAuth"]
	if !ok || len(buf) < 4 || binary.BigEndian.Uint32(buf) < uint32(requiredLevel) {
		return false
	}
	return true
}
