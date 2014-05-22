package sirencast

import (
	"log"
	"net/http"
)

type Environment struct {
	Server ServerEnvironment
}

type ServerEnvironment struct {
	BindAddress string
}

func Setup() (*Environment, error) {
	e := &Environment{
		Server: ServerEnvironment{
			BindAddress: ":9050",
		},
	}

	// Setup a listener for HTTP requests
	httpListener := NewHTTPListener()
	DefaultDetectors.Default = httpListener.Handler
	go func() {
		if err := http.Serve(httpListener, nil); err != nil {
			log.Println("HTTP server exited:", err)
		}
	}()

	return e, nil
}
