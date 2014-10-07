package config

var Default = Config{
	Addr: "localhost:9050",
	HTTP: HTTPServer{
		Disabled: false,
		Addr:     "",
	},
}
