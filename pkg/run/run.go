package run

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Escape-Technologies/goctopus/internal/config"
	"github.com/Escape-Technologies/goctopus/internal/utils"
	"github.com/Escape-Technologies/goctopus/internal/workers"
	"github.com/Escape-Technologies/goctopus/pkg/fingerprint"

	log "github.com/sirupsen/logrus"
)

// @todo refactor this to decouple from the filesystem
// make it a run function that takes a list of domains
// maybe a second function that takes a channel of domains
// the file io should be in internal/io
func RunFromFile(input *os.File) {
	count, err := utils.CountLines(input)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	inputBuffer := bufio.NewScanner(input)
	// Scan the file line by line
	inputBuffer.Split(bufio.ScanLines)
	// removes the file if it exists
	os.Remove(config.Conf.OutputFile)
	out, err := os.OpenFile(config.Conf.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer out.Close()

	output := make(chan *fingerprint.FingerprintOutput, config.Conf.MaxWorkers)
	go workers.Orchestrator(inputBuffer, config.Conf.MaxWorkers, output, count)

	// -- OUTPUT --
	for output := range output {
		log.Infof("Found: %+v\n", output)
		jsonOutput, err := json.Marshal(output)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		out.Write(jsonOutput)
		out.Write([]byte("\n"))
	}
}

//@todo run from list of domains
