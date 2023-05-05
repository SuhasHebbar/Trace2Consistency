package multitrace

import (
	flag "github.com/spf13/pflag"
)

type Config struct {
	ValidTraces  []string
	FaultyTraces []string
}

func GetConfig() *Config {
	// command line args
	faultyTraces := flag.StringArrayP("--faulty", "-f", []string{}, "List of all faulty traces!")

	// parse command line args
	flag.Parse()
	validTraces := flag.Args()

	return &Config{
		ValidTraces:  validTraces,
		FaultyTraces: *faultyTraces,
	}
}
