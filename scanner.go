package main

import (
	"fmt"
	"net"
	"sync"
)

var workers = 100
var Ports uint = 65535

func Scan(address string) chan uint {
	results := make(chan uint, workers)
	workload := make(chan uint)
	resultWaitGroup := sync.WaitGroup{}

	for range workers {
		go worker(&address, workload, results, &resultWaitGroup)
	}

	resultWaitGroup.Add(int(Ports))
	go func() {
		for port := uint(1); port <= Ports; port++ {
			workload <- port
		}
	}()

	go func() {
		resultWaitGroup.Wait()
		close(workload)
		close(results)
	}()

	return results
}

func worker(address *string, ports chan uint, results chan uint, resultWaitGroup *sync.WaitGroup) {
	for port := range ports {
		fullAddress := fmt.Sprintf("%s:%d", *address, port)
		conn, err := net.Dial("tcp", fullAddress)
		if err != nil {
			results <- 0
			resultWaitGroup.Done()
			continue
		}
		_ = conn.Close()
		results <- port
		resultWaitGroup.Done()
	}
}
