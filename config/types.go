package config

// Config is the configuration root and deals with configuration for
// the streaming server.
type Config struct {
	// Addr is the address to bind the server to, this is used for
	// the streaming component and optionally for the HTTP server
	// if no alternative address is used and the HTTP server isn't
	// disabled.
	Addr       string `json:"address"`
	HTTPServer `json:"http_server"`
}

// HTTPServer is an optional configuration for the HTTP server included
// with sirencast. This allows you to bind the HTTP server on a different
// address and/or port and disable it completely if so wished.
type HTTPServer struct {
	// Disabled indicates if we have the HTTP server disabled or not
	Disabled bool `json:"disabled"`
	// Addr is an alternative address to bind the HTTP server on,
	// instead of using `Server.Addr`.
	Addr string `json:"address,omitempty"`
}
