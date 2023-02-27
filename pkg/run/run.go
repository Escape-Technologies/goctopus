package run

import (
	"bufio"
	"fmt"
	"goctopus/internal/config"
	"goctopus/internal/utils"
	"goctopus/internal/workers"
	"os"

	log "github.com/sirupsen/logrus"
)

func Run(input *os.File, config config.Config) {
	count, err := utils.CountLines(input)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	inputBuffer := bufio.NewScanner(input)
	// Scan the file line by line
	inputBuffer.Split(bufio.ScanLines)
	// removes the file if it exists
	os.Remove(config.OutputFile)
	out, err := os.OpenFile(config.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer out.Close()

	endpoints := make(chan string, config.MaxWorkers)
	go workers.Orchestrator(inputBuffer, config.MaxWorkers, endpoints, count)

	// -- OUTPUT --
	for endpoint := range endpoints {
		log.Infof("Found endpoint: %v\n", endpoint)
		fmt.Fprintln(out, endpoint)
	}
}