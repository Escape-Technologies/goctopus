package config

import (
	"errors"
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	InputFile            string
	OutputFile           string
	MaxWorkers           int
	Verbose              bool
	Silent               bool
	Timeout              int
	Introspection        bool
	FieldSuggestion      bool
	WebhookUrl           string
	SubdomainEnumeration bool
}

var (
	c *Config
)

func Get() *Config {
	if c == nil {
		log.Panic("Config not initialized")
	}
	return c
}

func LoadFromArgs() {
	config := Config{}
	// -- INPUT --
	flag.StringVar(&config.InputFile, "i", "", "Input file")
	// @todo accept stdin args as input (addresses comma separated)

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
	flag.BoolVar(&config.SubdomainEnumeration, "subdomain", false, "Enable subdomain enumeration")

	flag.Parse()

	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	if config.Silent {
		log.SetLevel(log.ErrorLevel)
	}

	if err := ValidateConfig(&config); err != nil {
		log.Error(err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	c = &config
}

func ValidateConfig(conf *Config) error {
	if conf.MaxWorkers < 1 {
		return errors.New("[Invalid config] Max workers must be greater than 0")
	}

	if conf.Timeout < 1 {
		return errors.New("[Invalid config] Timeout must be greater than 0")
	}

	if !conf.Introspection && conf.FieldSuggestion {
		return errors.New("[Invalid config] Introspection has to be enabled to use field suggestion fingerprinting")
	}

	if conf.InputFile == "" {
		return errors.New("[Invalid config] Please specify an input file")
	}

	return nil
}

func Load(config *Config) {
	if err := ValidateConfig(config); err != nil {
		log.Error(err)
		flag.PrintDefaults()
		os.Exit(1)
	}
	c = config
}
