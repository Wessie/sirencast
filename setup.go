package sirencast

import (
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/Wessie/sirencast/config"
)

func Setup() (*config.Config, error) {
	if !config.Loaded {
		return nil, errors.New("unable to load configuration")
	}

	// TODO: Move all of this into reusable functions
	// Setup a listener for HTTP requests
	if config.Active.HTTPServer.Disabled {
		return config.Active, nil
	}

	var l net.Listener
	// Setup a protocol detector default and a fake listener
	// for HTTP if the configuration tells us to not run the
	// HTTP server on a different address.
	if config.Active.HTTPServer.Addr == "" {
		httpListener := NewHTTPListener(config.Active.Addr)
		DefaultDetectors.Default = httpListener.Handler
		l = httpListener
	} else {
		var err error
		l, err = net.Listen("tcp", config.Active.HTTPServer.Addr)
		if err != nil {
			log.Printf("HTTP unable to listen on '%s': %s\n", config.Active.HTTPServer.Addr, err)
		}
	}
	log.Printf("HTTP server listening on '%s'\n", l.Addr())

	go func() {
		if err := http.Serve(l, nil); err != nil {
			log.Println("HTTP server exited:", err)
		}
	}()

	return config.Active, nil
}
