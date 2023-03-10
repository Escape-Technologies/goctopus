package config

import (
	"flag"
	"os"

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
	NoSubdomain     bool
}

func ParseFlags() {
	config := Config{}
	// -- INPUT --
	flag.StringVar(&config.InputFile, "i", "", "Input file")
	// @TODO
	// flag.StringVar(&config.InputFile, "d", "", "Input domains (comma separated)")
	// flag.StringVar(&config.InputFile, "u", "", "Input urls (comma separated)")

	// -- CONFIG --
	// @todo make output file optional ?
	flag.StringVar(&config.OutputFile, "o", "output.jsonl", "Output file (json-lines format)")
	flag.StringVar(&config.WebhookUrl, "webhook", "", "Webhook URL")
	flag.IntVar(&config.MaxWorkers, "w", 100, "Max workers")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose")
	flag.BoolVar(&config.Silent, "s", false, "Silent")
	flag.IntVar(&config.Timeout, "t", 30, "Request timeout (seconds)")
	flag.BoolVar(&config.Introspection, "introspect", false, "Enable introspection fingerprinting")
	flag.BoolVar(&config.FieldSuggestion, "suggest", false, "Enable fields suggestion fingerprinting.\nNeeds \"introspection\" to be enabled.")
	flag.BoolVar(&config.NoSubdomain, "no-subdomain", false, "Disable subdomain enumeration")

	flag.Parse()

	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	if config.Silent {
		log.SetLevel(log.ErrorLevel)
	}

	ValidateConfig(&config)
	Conf = &config
}

func ValidateConfig(conf *Config) {
	if conf.MaxWorkers < 1 {
		log.Error("[Invalid args] Max workers must be greater than 0")
		configError()
	}

	if conf.Timeout < 1 {
		log.Error("[Invalid args] Timeout must be greater than 0")
		configError()
	}

	if !conf.Introspection && conf.FieldSuggestion {
		log.Error("[Invalid args] Introspection has to be enabled to use field suggestion fingerprinting")
		configError()
	}

	if conf.InputFile == "" {
		log.Error("[Invalid args] Please specify an input file")
		configError()
	}

}

func configError() {
	flag.PrintDefaults()
	os.Exit(1)
}
