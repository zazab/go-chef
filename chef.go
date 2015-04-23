package chef

import (
	"github.com/op/go-logging"
	"github.com/zazab/chef-golang"
)

// Needed for defining functions on chef.Chef struct
type Chef struct {
	chef.Chef
	log *logging.Logger
}

func Connect(filename ...string) (*Chef, error) {
	c, err := chef.Connect(filename...)
	if err != nil {
		return nil, err
	}
	return &Chef{*c, logging.MustGetLogger("chef.client")}, nil
}
