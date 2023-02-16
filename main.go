package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"goctopus/fingerprint"
	"goctopus/utils"

	log "github.com/sirupsen/logrus"
)

func worker(domains chan string, endpoints chan string, workerId int, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debugf("Worker %d instantiated\n", workerId)
	for domain := range domains {
		log.Debugf("Worker %d started on: %v\n", workerId, domain)
		endpoint, err := fingerprint.FingerprintDomain(domain)
		if err == nil {
			log.Debugf("Worker %d found endpoint: %v\n", workerId, endpoint)
			endpoints <- endpoint
		}	
	}
	log.Debugf("Worker %d finished\n", workerId)
}

func orchestrator(inputBuffer *bufio.Scanner, maxWorkers int, endpoints chan string) {

	domains := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		go worker(domains, endpoints, i, &wg)
	}

	for inputBuffer.Scan() {
		domain := inputBuffer.Text()
		log.Debugf("Adding %v to the queue\n", domain)
		domains <- domain
	}

	close(domains)
	log.Debugf("Orchestrator finished, waiting for workers to finish...")
	wg.Wait()
	close(endpoints)
	log.Debugf("All workers finished")
}

func main() {
	// -- PARAMS --
	// @todo pass this in args/flags
	// @todo make a shared config package: https://stackoverflow.com/questions/36528091/golang-sharing-configurations-between-packages
	inputFile := "input.txt"
	outputFile := "endpoints.txt"
	maxWorkers := 200

	log.SetLevel(log.DebugLevel)

	input, err := os.Open(inputFile)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer input.Close()

	count, err := utils.CountLines(input)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Infof("Starting goctopus on %d endpoints...\n", count)
	inputBuffer := bufio.NewScanner(input)
	// Scan the file line by line
	inputBuffer.Split(bufio.ScanLines)
	// removes the file if it exists
	os.Remove(outputFile)
	out, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer out.Close()

	endpoints := make(chan string, maxWorkers)
	go orchestrator(inputBuffer, maxWorkers, endpoints)

	// -- OUTPUT --
	for endpoint := range endpoints {
		log.Infof("Found endpoint: %v\n", endpoint)
		fmt.Fprintln(out, endpoint)
	}
}