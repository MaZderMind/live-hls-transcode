package main

import (
	"github.com/gosimple/slug"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type StreamStatusManager struct {
	tempDir    string
	streamInfo map[string]StreamInfo
	transcoder Transcoder
}

func NewStreamStatusManager(tempDir string) StreamStatusManager {
	err := os.MkdirAll(tempDir, 0770)
	if err != nil {
		log.Fatalf("Cannot create Temp-Dir %s", tempDir)
	}

	return StreamStatusManager{
		tempDir,
		make(map[string]StreamInfo),
		NewTranscoder(),
	}
}

// Stored Information about a Stream
type StreamInfo struct {
	CalculatedPath string
	UrlPath        string
	TempDir        string
	handle         *TranscoderHandle
	LastAccess     time.Time
}

func (info StreamInfo) DominantStatusCode() StreamStatus {
	if info.handle != nil {
		if info.handle.IsFinished() {
			return TranscodingFinished
		} else if info.handle.IsReady() {
			return StreamReady
		} else if info.handle.IsRunning() {
			return StreamInPreparation
		} else {
			return StreamTranscodingFailed
		}
	} else {
		return NoStream
	}
}

func (info StreamInfo) IsReady() bool {
	return info.handle != nil && info.handle.IsReady()
}

func (info StreamInfo) IsRunning() bool {
	return info.handle != nil && info.handle.IsRunning()
}

func (info StreamInfo) IsFinished() bool {
	return info.handle != nil && info.handle.IsFinished()
}

func (info StreamInfo) TotalDuration() time.Duration {
	if info.handle != nil {
		return info.handle.TotalDuration()
	}

	return time.Duration(0)
}

func (info StreamInfo) ProcessedDuration() time.Duration {
	if info.handle != nil {
		return info.handle.ProcessedDuration()
	}

	return time.Duration(0)
}

func (info StreamInfo) ProcessedPercent() float64 {
	if info.handle != nil {
		return info.handle.ProcessedPercent()
	}

	return 0.
}

type StreamStatus int

const (
	NoStream = iota
	StreamTranscodingFailed
	StreamInPreparation
	StreamReady
	TranscodingFinished
)

func (manager StreamStatusManager) StreamInfo(calculatedPath string) StreamInfo {
	return manager.streamInfo[calculatedPath]
}

func (manager StreamStatusManager) StartStream(calculatedPath string, urlPath string) {
	_, hasInfo := manager.streamInfo[calculatedPath]
	if hasInfo {
		log.Printf("%s: Stream already active", calculatedPath)
		return
	}

	tempDir := manager.createTempDir(calculatedPath)

	log.Printf("%s: Starting Stream-Transcoder into %s", calculatedPath, tempDir)
	handle := manager.transcoder.StartTranscoder(calculatedPath, tempDir)

	manager.streamInfo[calculatedPath] = StreamInfo{
		calculatedPath,
		urlPath,
		tempDir,
		handle,
		time.Now(),
	}
}

func (manager StreamStatusManager) createTempDir(calculatedPath string) string {
	tempDir, err := ioutil.TempDir(manager.tempDir, slug.Make(calculatedPath))

	if err != nil {
		log.Fatal(err)
	}

	return tempDir
}

func (manager StreamStatusManager) StopStream(calculatedPath string) {
	info, hasInfo := manager.streamInfo[calculatedPath]
	if hasInfo {
		if info.handle.IsFinished() {
			log.Printf("%s: Stream-Transcoder already finished, ignoring Stop-Command", calculatedPath)
		} else {
			log.Printf("%s: Stopping unfinished Stream-Transcoder", calculatedPath)
			info.handle.Stop()

			delete(manager.streamInfo, calculatedPath)
		}
	}
}

func (manager StreamStatusManager) OtherRunningTranscoders(excludingThisCalculatedPath string) []StreamInfo {
	otherRunningTranscoders := make([]StreamInfo, 0)

	for _, streamInfo := range manager.streamInfo {
		if streamInfo.CalculatedPath == excludingThisCalculatedPath || !streamInfo.IsRunning() {
			continue
		}

		otherRunningTranscoders = append(otherRunningTranscoders, streamInfo)
	}

	return otherRunningTranscoders
}

func (manager StreamStatusManager) UpdateLastAccess(calculatedPath string) {
	info, hasInfo := manager.streamInfo[calculatedPath]

	if hasInfo {
		info.LastAccess = time.Now()
		manager.streamInfo[calculatedPath] = info
	}
}

func (manager StreamStatusManager) StreamInfos() map[string]StreamInfo {
	return manager.streamInfo
}

func (manager StreamStatusManager) DeleteStreamInfo(calculatedPath string) {
	delete(manager.streamInfo, calculatedPath)
}
