package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"goctopus/fingerprint"
)

func worker(domain string, endpoints chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	endpoint, err := fingerprint.FingerprintDomain(domain)
	if err == nil {
		endpoints <- endpoint
	}
}

func orchestrator(inputBuffer *bufio.Scanner, maxWorkers int, endpoints chan string) {
	var wg sync.WaitGroup

	// this is nasty, but it works
	// it will read the file line by line and spawn maxWorkers workers
	for {
		for i := 0; i < maxWorkers; i++ {
			if !inputBuffer.Scan() {
				break
			}
			wg.Add(1)
			go worker(inputBuffer.Text(), endpoints, &wg)
		}
		wg.Wait()
		if !inputBuffer.Scan() {
			break
		}
	}
	close(endpoints)
}

func main() {
	// -- PARAMS --
	// @todo pass this in args/flags
	inputFile := "input.txt"
	outputFile := "endpoints.txt"
	maxWorkers := 100

	input, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer input.Close()
	inputBuffer := bufio.NewScanner(input)
	// Scan the file line by line
	inputBuffer.Split(bufio.ScanLines)
	// removes the file if it exists
	os.Remove(outputFile)
	out, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer out.Close()

	endpoints := make(chan string)
	go orchestrator(inputBuffer, maxWorkers, endpoints)

	// -- OUTPUT --
	for endpoint := range endpoints {
		fmt.Fprintln(out, endpoint)
	}
}