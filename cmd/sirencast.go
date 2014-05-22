package main

import (
	"log"

	"github.com/Wessie/sirencast"
	// "github.com/Wessie/sirencast/sirencast/icecast"
)

func main() {
	environment, err := sirencast.Setup()

	if err != nil {
		log.Fatal(err)
	}

	if err = sirencast.Run(environment); err != nil {
		log.Fatal(err)
	}
}
