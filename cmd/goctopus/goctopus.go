package main

import (
	"os"

	"github.com/Escape-Technologies/goctopus/internal/config"
	"github.com/Escape-Technologies/goctopus/internal/utils"
	"github.com/Escape-Technologies/goctopus/pkg/run"

	log "github.com/sirupsen/logrus"
)

func main() {
	// -- PARAMS --
	// @todo make config a singleton package: https://stackoverflow.com/questions/36528091/golang-sharing-configurations-between-packages
	config.ParseFlags()
	if !config.Conf.Silent {
		utils.PrintASCII()
	}

	input, err := os.Open(config.Conf.InputFile)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer input.Close()

	run.RunFromFile(input)
}
