package main

import (
	"bufio"
	"flag"
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

func orchestrator(inputBuffer *bufio.Scanner, maxWorkers int, endpoints chan string, count int) {

	domains := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		go worker(domains, endpoints, i, &wg)
	}

	i := 0
	for inputBuffer.Scan() {
		domain := inputBuffer.Text()
		log.Infof("(%d/%d) Adding %v to the queue\n", i,count, domain)
		domains <- domain
		i++
	}

	close(domains)
	log.Debugf("Orchestrator finished, waiting for workers to finish...")
	wg.Wait()
	close(endpoints)
	log.Debugf("All workers finished")
}

type Config struct {
	InputFile string
	OutputFile string
	MaxWorkers int
	Verbose bool
}

func parseFlags() Config {
	config := Config{}
	flag.StringVar(&config.InputFile, "i", "input.txt", "Input file")
	flag.StringVar(&config.OutputFile, "o", "endpoints.txt", "Output file")
	flag.IntVar(&config.MaxWorkers, "w", 10, "Max workers")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose")
	flag.Parse()
	return config
}

func main() {
	// -- PARAMS --
	// @todo pass this in args/flags
	// @todo make a shared config package: https://stackoverflow.com/questions/36528091/golang-sharing-configurations-between-packages
	config := parseFlags()

	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	input, err := os.Open(config.InputFile)
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
	os.Remove(config.OutputFile)
	out, err := os.OpenFile(config.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer out.Close()

	endpoints := make(chan string, config.MaxWorkers)
	go orchestrator(inputBuffer, config.MaxWorkers, endpoints, count)

	// -- OUTPUT --
	for endpoint := range endpoints {
		log.Infof("Found endpoint: %v\n", endpoint)
		fmt.Fprintln(out, endpoint)
	}
}