package main

import (
	"os"

	"goctopus/internal/config"
	"goctopus/internal/utils"
	"goctopus/pkg/run"

	log "github.com/sirupsen/logrus"
)



func main() {
	// -- PARAMS --
	// @todo make config a singleton package: https://stackoverflow.com/questions/36528091/golang-sharing-configurations-between-packages
	utils.PrintASCII()
	config := config.ParseFlags()

	input, err := os.Open(config.InputFile)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer input.Close()

	run.Run(input, config)
}