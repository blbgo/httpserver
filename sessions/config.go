package sessions

import (
	"strconv"

	"github.com/blbgo/general"
)

// Config must be implemented and provided to New
type Config interface {
	Name() string
	MaxAge() int
}

type config struct {
	NameValue   string
	MaxAgeValue int
}

// NewConfig provides a Config based on general.Config
func NewConfig(c general.Config) (Config, error) {
	r := &config{}
	var err error

	r.NameValue, err = c.Value("HTTPSessions", "Name")
	if err != nil {
		return nil, err
	}
	MaxAgeString, err := c.Value("HTTPSessions", "MaxAge")
	if err != nil {
		return nil, err
	}
	r.MaxAgeValue, err = strconv.Atoi(MaxAgeString)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *config) Name() string {
	return r.NameValue
}

func (r *config) MaxAge() int {
	return r.MaxAgeValue
}
