package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Cleanup struct {
	ticker          *time.Ticker
	stop            chan bool
	statusManager   *StreamStatusManager
	lifetimeMinutes uint32
}

func NewCleanup(statusManager *StreamStatusManager, lifetimeMinutes uint32) Cleanup {
	return Cleanup{
		time.NewTicker(time.Hour),
		make(chan bool),
		statusManager,
		lifetimeMinutes,
	}
}

func (cleanup Cleanup) Start() {
	go func() {
		for {
			select {
			case <-cleanup.stop:
				return
			case <-cleanup.ticker.C:
				cleanup.doRun()
			}
		}
	}()
}

func (cleanup Cleanup) Stop() {
	cleanup.stop <- true
}

func (cleanup Cleanup) doRun() {
	log.Print("Cleanup: Running")

	streams := cleanup.statusManager.StreamInfos()
	for calculatedPath, info := range streams {
		expirationDate := info.LastAccess.Add(time.Minute * time.Duration(cleanup.lifetimeMinutes))
		if time.Now().After(expirationDate) && ! info.IsRunning() {
			log.Printf("Cleanup: Deleting %s (%s)", calculatedPath, info.TempDir)

			err := os.RemoveAll(info.TempDir)
			if err != nil {
				fmt.Printf("Cleanup: Unable to Delete %s", info.TempDir)
				continue
			}

			cleanup.statusManager.DeleteStreamInfo(calculatedPath)
		}
	}
}
