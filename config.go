package httpserver

import (
	"github.com/blbgo/general"
)

// Config must be implemented and provided to New
type Config interface {
	// Addr must return the address and port the server should listen on
	Addr() string
	HTTPS() bool
	CertFile() string
	KeyFile() string
}

type config struct {
	AddrValue     string
	HTTPSValue    bool
	CertFileValue string
	KeyFileValue  string
}

// NewConfig provides a Config based on general.Config
func NewConfig(c general.Config) (Config, error) {
	r := &config{}
	var err error

	r.AddrValue, err = c.Value("HTTPServer", "Addr")
	if err != nil {
		return nil, err
	}
	httpsString, err := c.Value("HTTPServer", "HTTPS")
	if err != nil {
		return nil, err
	}
	r.HTTPSValue = httpsString == "true"
	r.CertFileValue, err = c.Value("HTTPServer", "CertFile")
	if err != nil {
		return nil, err
	}
	r.KeyFileValue, err = c.Value("HTTPServer", "KeyFile")
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Addr method of httpserver.Config, returns the address and port the http server should run
// under
func (r *config) Addr() string {
	return r.AddrValue
}

// HTTPS indicates if the server should be run in https mode
func (r *config) HTTPS() bool {
	return r.HTTPSValue
}

// CertFile returns the cert tile name, only used if Https returns true
func (r *config) CertFile() string {
	return r.CertFileValue
}

// KeyFile returns the key tile name, only used if Https returns true
func (r *config) KeyFile() string {
	return r.KeyFileValue
}
