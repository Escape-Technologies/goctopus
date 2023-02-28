package config

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

var (
	Conf *Config
)

type Config struct {
	InputFile     string
	OutputFile    string
	MaxWorkers    int
	Verbose       bool
	Silent        bool
	Timeout       int
	Introspection bool
}

func ParseFlags() {
	config := Config{}
	flag.StringVar(&config.InputFile, "i", "input.txt", "Input file")
	flag.StringVar(&config.OutputFile, "o", "output.jsonl", "Output file (json-lines format)")
	flag.IntVar(&config.MaxWorkers, "w", 10, "Max workers")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose")
	flag.BoolVar(&config.Silent, "s", false, "Silent")
	flag.IntVar(&config.Timeout, "t", 2, "Request timeout (seconds)")
	flag.BoolVar(&config.Introspection, "introspection", false, "Enable introspection fingerprinting (default: true)")

	flag.Parse()

	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	if config.Silent {
		log.SetLevel(log.ErrorLevel)
	}

	Conf = &config
}
