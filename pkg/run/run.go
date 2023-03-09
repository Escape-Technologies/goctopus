package run

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/Escape-Technologies/goctopus/internal/config"
	"github.com/Escape-Technologies/goctopus/internal/http"
	"github.com/Escape-Technologies/goctopus/internal/utils"
	"github.com/Escape-Technologies/goctopus/internal/workers"
	out "github.com/Escape-Technologies/goctopus/pkg/output"

	log "github.com/sirupsen/logrus"
)

// @todo refactor this to decouple from the filesystem and cli
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
	outputFile, err := os.OpenFile(config.Conf.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// Limit the number of workers to the number of domains if the number of domains is less than the max workers
	maxWorkers := utils.MinInt(count, config.Conf.MaxWorkers)
	log.Infof("Starting %d workers\n", maxWorkers)
	output := make(chan *out.FingerprintOutput, maxWorkers)
	go workers.Orchestrator(inputBuffer, maxWorkers, output, count)

	foundCount := 0
	// -- OUTPUT --
	var wg sync.WaitGroup
	for out := range output {
		foundCount++
		jsonOutput, err := json.Marshal(out)
		log.Infof("Found: %+v\n", string(jsonOutput))
		if err != nil {
			log.Error(err)
		}
		wg.Add(1)
		go http.SendToWebhook(jsonOutput, &wg)
		outputFile.Write(jsonOutput)
		outputFile.Write([]byte("\n"))
	}
	wg.Wait()
	log.Infof("Done. Found %d graphql endpoints", foundCount)
}

//@todo run from list of domains
