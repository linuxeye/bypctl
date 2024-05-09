package cmder

import "github.com/scottames/cmder/pkg/log"

var logger log.Logger

// getLogger returns the cmder implementation of the Logger interface
//
// If SetLogger called prior to NewLogger the Logger passed to SetLogger will be returned
// otherwise a new Logger will be returned
func getLogger() log.Logger {
	if logger != nil {
		return logger
	}

	logger = log.New()

	return logger
}

// SetLogger allows for setting the cmder from an external implementation
func SetLogger(l log.Logger) {
	logger = l
}
