package config

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

var (
	Conf *Config
)

type Config struct {
	InputFile       string
	OutputFile      string
	MaxWorkers      int
	Verbose         bool
	Silent          bool
	Timeout         int
	Introspection   bool
	FieldSuggestion bool
	WebhookUrl      string
}

func ParseFlags() {
	config := Config{}
	flag.StringVar(&config.InputFile, "i", "input.txt", "Input file")
	flag.StringVar(&config.OutputFile, "o", "output.jsonl", "Output file (json-lines format)")
	flag.StringVar(&config.WebhookUrl, "webhook", "", "Webhook URL")
	flag.IntVar(&config.MaxWorkers, "w", 10, "Max workers")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose")
	flag.BoolVar(&config.Silent, "s", false, "Silent")
	flag.IntVar(&config.Timeout, "t", 2, "Request timeout (seconds)")
	flag.BoolVar(&config.Introspection, "introspect", false, "Enable introspection fingerprinting")
	flag.BoolVar(&config.FieldSuggestion, "suggest", false, "Enable fields suggestion fingerprinting.\nNeeds \"introspection\" to be enabled.")

	flag.Parse()

	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	if config.Silent {
		log.SetLevel(log.ErrorLevel)
	}

	Conf = &config
}
