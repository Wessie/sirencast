package config

var Default = Config{
	Addr: "localhost:9050",
	HTTPServer: HTTPServer{
		Disabled: false,
		Addr:     "",
	},
}
