package config

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/Escape-Technologies/goctopus/internal/utils"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	InputFile            string
	Addresses            []string
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

func parseArgs() []string {
	if flag.Arg(0) == "" {
		log.Error("Invalid addresses argument")
		utils.PrintUsage()
		os.Exit(1)
	}
	input := flag.Arg(0)
	return strings.Split(input, ",")
}

func LoadFromArgs() {
	flag.Usage = utils.PrintUsage
	config := Config{}
	// -- INPUT --
	flag.StringVar(&config.InputFile, "f", "", "Input file")

	// -- CONFIG --
	flag.StringVar(&config.OutputFile, "o", "", "Output file (json-lines format)")
	flag.StringVar(&config.WebhookUrl, "webhook", "", "Webhook URL")
	flag.IntVar(&config.MaxWorkers, "w", 40, "Max workers")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose")
	flag.BoolVar(&config.Silent, "s", false, "Silent")
	flag.IntVar(&config.Timeout, "t", 30, "Request timeout (seconds)")
	flag.BoolVar(&config.Introspection, "introspect", false, "Enable introspection fingerprinting")
	flag.BoolVar(&config.FieldSuggestion, "suggest", false, "Enable fields suggestion fingerprinting.\nNeeds \"introspection\" to be enabled.")
	flag.BoolVar(&config.SubdomainEnumeration, "subdomain", false, "Enable subdomain enumeration")

	flag.Parse()

	if config.InputFile == "" {
		config.Addresses = parseArgs()
	}

	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	if config.Silent {
		log.SetLevel(log.ErrorLevel)
	}

	if err := validateConfig(&config, true); err != nil {
		log.Error(err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	c = &config
}

// Validates the config.
// `cli` is used to determine if the config is loaded from the CLI or from a file.
// If cli is false, then the addresses check is skipped.
func validateConfig(conf *Config, isCli bool) error {
	if conf.MaxWorkers < 1 {
		return errors.New("[Invalid config] Max workers must be greater than 0")
	}

	if conf.Timeout < 1 {
		return errors.New("[Invalid config] Timeout must be greater than 0")
	}

	if !conf.Introspection && conf.FieldSuggestion {
		return errors.New("[Invalid config] Introspection has to be enabled to use field suggestion fingerprinting")
	}

	if isCli && conf.InputFile == "" && len(conf.Addresses) == 0 {
		return errors.New("[Invalid config] Please specify an input file or a list of addresses")
	}

	return nil
}

func Load(config *Config) {
	if err := validateConfig(config, false); err != nil {
		log.Error(err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	c = config

	if c.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	if c.Silent {
		log.SetLevel(log.ErrorLevel)
	}
}
