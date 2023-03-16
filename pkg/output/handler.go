package output

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/Escape-Technologies/goctopus/pkg/config"
	"github.com/Escape-Technologies/goctopus/pkg/http"
	log "github.com/sirupsen/logrus"
)

func openOutputFile(config *config.Config) (*os.File, error) {
	// removes the file if it exists
	os.Remove(config.OutputFile)
	outputFile, err := os.OpenFile(config.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return outputFile, err
}

func handleSingleOutput(output *FingerprintOutput, outputFile *os.File, wg *sync.WaitGroup, config *config.Config) {
	isOutputFile := config.OutputFile != ""
	isWebhook := config.WebhookUrl != ""

	jsonOutput, err := json.Marshal(output)
	log.Infof("Found: %+v\n", string(jsonOutput))
	if err != nil {
		log.Error(err)
	}

	if isOutputFile {
		if err != nil {
			log.Error(err)
		}
		content := append(jsonOutput, []byte("\n")...)
		outputFile.Write(content)
	}

	if isWebhook {
		if err != nil {
			log.Error(err)
		}
		wg.Add(1)
		go http.SendToWebhook(config.WebhookUrl, jsonOutput, wg)
	}

}

func HandleOutput(output chan *FingerprintOutput) int {
	c := config.Get()
	outputFile, err := openOutputFile(c)
	if err != nil {
		log.Panic(err)
	}

	foundCount := 0
	var wg sync.WaitGroup
	for out := range output {
		foundCount++
		handleSingleOutput(out, outputFile, &wg, c)
	}

	wg.Wait()
	return foundCount
}
