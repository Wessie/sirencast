package sirencast

import (
	"errors"
	"log"
	"net/http"

	"github.com/Wessie/sirencast/config"
)

func Setup() (*config.Config, error) {
	if !config.Loaded {
		return nil, errors.New("unable to load configuration")
	}

	// Setup a listener for HTTP requests
	httpListener := NewHTTPListener()
	DefaultDetectors.Default = httpListener.Handler
	go func() {
		if err := http.Serve(httpListener, nil); err != nil {
			log.Println("HTTP server exited:", err)
		}
	}()

	return config.Active, nil
}
