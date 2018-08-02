package main

// inputFlagsStruct contains all possible verbosity output levels available as CLI args
type inputFlagsStruct struct {
	Debug   bool
	Error   bool
	Quiet   bool
	Color   bool
	AuthKey string
	Sleep   int
}

// ParseVerbosity selects proper verbosity level based on verbosity flag values
func (v *inputFlagsStruct) ParseVerbosity() (verbosityLevel int8) {
	if inputFlags.Debug {
		verbosityLevel = VerbosityDebug
	} else if inputFlags.Quiet {
		verbosityLevel = VerbosityQuiet
	} else if inputFlags.Error {
		verbosityLevel = VerbosityError
	} else {
		verbosityLevel = VerbosityNormal
	}

	return verbosityLevel
}

// StringVerbosity converts verbosity level to text name representation
func (v *inputFlagsStruct) StringVerbosity() string {
	if inputFlags.Debug {
		return `Verbose`
	} else if inputFlags.Quiet {
		return `Quiet`
	} else if inputFlags.Error {
		return `Errors Only`
	} else {
		return `Normal`
	}
}
