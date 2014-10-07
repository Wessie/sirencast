package main

import (
	"log"

	"github.com/Wessie/sirencast"
	"github.com/Wessie/sirencast/icecast"
	_ "github.com/Wessie/sirencast/web"
)

func main() {
	environment, err := sirencast.Setup()
	if err != nil {
		log.Fatal(err)
	}

	ice := icecast.NewServer()
	sirencast.RegisterDetector(ice.Detect)

	if err = sirencast.Run(environment); err != nil {
		log.Fatal(err)
	}
}
