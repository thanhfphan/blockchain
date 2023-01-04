package logging

type Logger interface {
	// The program should exit soon after this called
	Fatalf(format string, args ...interface{})
	// The program should be able to recover from this error
	Errorf(format string, args ...interface{})
	// Log that event has occurred that may indicate a future error
	Warnf(format string, args ...interface{})
	// Log an event that maybe useful for user
	Infof(format string, args ...interface{})
	// Log an event that maybe useful for a programmer when debugging
	Debugf(format string, args ...interface{})
	// Lo extremely detailed events that can be useful for inspecting the program
	Verbof(format string, args ...interface{})

	Stop()
}
