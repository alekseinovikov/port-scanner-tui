package main

import "sync"

var (
	foundPorts []uint = make([]uint, 0)
	finished   bool   = false
	scanned    int    = 0
	rwLock     sync.RWMutex
)

func StartScanning(address string) {
	results := Scan(address)

	go func() {
		for port := range results {
			rwLock.Lock()
			scanned++
			if port == 0 {
				rwLock.Unlock()
				continue
			}

			foundPorts = append(foundPorts, port)
			rwLock.Unlock()
		}

		rwLock.Lock()
		finished = true
		rwLock.Unlock()
	}()
}

func GetFoundPorts() (bool, []uint, int) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	return finished, foundPorts[:], scanned
}
