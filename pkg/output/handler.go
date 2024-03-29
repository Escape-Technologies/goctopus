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
	if config.OutputFile == "" {
		return nil, nil
	}
	os.Remove(config.OutputFile)
	outputFile, err := os.OpenFile(config.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return outputFile, err
}

func handleSingleOutput(output *FingerprintOutput, outputFile *os.File, wg *sync.WaitGroup, config *config.Config) {
	hasOutputFile := config.OutputFile != ""
	hasWebhook := config.WebhookUrl != ""

	jsonOutput, err := json.Marshal(output)
	log.Infof("Found: %+v\n", string(jsonOutput))
	if err != nil {
		log.Error(err)
	}

	if hasOutputFile {
		content := append(jsonOutput, []byte("\n")...)
		if _, err := outputFile.Write(content); err != nil {
			log.Error(err)
		}
	}

	if hasWebhook {
		wg.Add(1)
		go func() {
			if err := http.SendToWebhook(config.WebhookUrl, jsonOutput, wg); err != nil {
				log.Error(err)
			}
		}()
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
