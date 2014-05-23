package config

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
)

// DefaultFilename is the filename used as configuration if none is supplied
const DefaultFilename = "sirencast.json"

// Filename is the filename of the file we read the configuration from.
var Filename string

// Active is the active configuration, and should be referenced in code using the
// configuration. The configuration will be updated atomically if required.
var Active *Config = &Default

// Loaded indicates if we were able to load a configuration successfully or not. It
// is suggested to not continue executing if this variable is false when main is executed.
var Loaded = false

func init() {
	flag.StringVar(&Filename, "conf", DefaultFilename, "configuration file to load")
	allowDefault := flag.Bool("default", false, "allow default configuration to be used if no file could be loaded")
	writeDefault := flag.Bool("write-default", false, "write the default configuration to the specified `conf` filename if unable to load it")
	flag.Parse()

	f, err := os.Open(Filename)
	if err != nil {
		log.Printf("unable to load configuration file '%s'", Filename)
		goto defaults
	}
	defer f.Close()

	if err := ReadConfig(Active, f); err != nil {
		log.Printf("unable to parse configuration file: %s", err)
		goto defaults
	}

	Loaded = true
	return

defaults:
	Loaded = *allowDefault

	// Nothing left to do if this is false
	if !*writeDefault {
		return
	}

	// else write the default config to file
	if err := CreateDefault(Filename); err != nil {
		log.Printf("unable to write default configuration to file '%s': %s", Filename, err)
	}
}

// ReadConfig reads a configuration from `r` and stores it
// into `conf`. The current format is JSON.
func ReadConfig(conf *Config, r io.Reader) error {
	return json.NewDecoder(r).Decode(conf)
}

// WriteConfig writes the configuration `conf` to writer
// `w`. The current format is JSON.
func WriteConfig(w io.Writer, conf Config) error {
	return json.NewEncoder(w).Encode(conf)
}

// CreateDefault creates file `filename` and writes the
// default configuration to it with `WriteConfig`.
func CreateDefault(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return WriteConfig(f, Default)
}
