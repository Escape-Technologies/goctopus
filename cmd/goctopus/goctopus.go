package main

import (
	"os"

	"github.com/Escape-Technologies/goctopus/internal/utils"
	"github.com/Escape-Technologies/goctopus/pkg/config"
	"github.com/Escape-Technologies/goctopus/pkg/goctopus"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.LoadFromArgs()
	if !config.Get().Silent {
		utils.PrintASCII()
	}

	if config.Get().InputFile != "" {
		input, err := os.Open(config.Get().InputFile)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		defer input.Close()
		goctopus.FingerprintFromFile(input)
	} else {
		goctopus.FingerprintFromSlice(config.Get().Addresses)
	}
}
