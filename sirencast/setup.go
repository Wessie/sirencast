package sirencast

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

	return e, nil
}
