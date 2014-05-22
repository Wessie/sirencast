package sirencast

import (
	"sync"
)

// DefaultDetectors are the global default detectors
var DefaultDetectors = NewDetectors()

// ConnHandler is a function that can handle a SirenConn, these are returned
// by detectors.
type ConnHandler func(*SirenConn)

// Detector is a function that should read from the Peeker and determine if it
// knows a ConnHandler that can handle the stream of data. Returning an appropriate
// ConnHandler if one is known, or nil if not.
//
// A detector should be reasonably fast in determining if it has a handler available
// or not due to being the entry-point of all connections.
type Detector func(Peeker) ConnHandler

type Detectors struct {
	mu        *sync.RWMutex
	Detectors []Detector
	Default   ConnHandler
}

// NewDetectors returns a new *Detectors
func NewDetectors() *Detectors {
	return &Detectors{
		mu:        new(sync.RWMutex),
		Detectors: make([]Detector, 0),
	}
}

// Register registers a new detector in the Detectors. The Detector
// is called when Detect is called.
func (ds *Detectors) Register(d Detector) {
	ds.mu.Lock()
	ds.Detectors = append(ds.Detectors, d)
	ds.mu.Unlock()
}

// Detect tries to detect what kind of stream we're receiving
// by letting all registered detectors peek at the front of
// the stream.
//
// Detect returns on the first non-nil return value from a Detector
func (ds *Detectors) Detect(input Peeker) (handler ConnHandler) {
	ds.mu.RLock()
	for _, d := range ds.Detectors {
		handler = d(input)

		// Reset for next detector or for return
		input.Reset()

		if handler != nil {
			break
		}
	}

	if handler == nil {
		return ds.Default
	}
	ds.mu.RUnlock()

	return handler
}

// RegisterDetector calls DefaultDetectors.Register
func RegisterDetector(d Detector) {
	DefaultDetectors.Register(d)
}

// Detect calls DefaultDetectors.Detect
func Detect(input Peeker) ConnHandler {
	return DefaultDetectors.Detect(input)
}
