package main

import (
	"github.com/gosimple/slug"
	"io/ioutil"
	"log"
	"os"
)

// Stored Information about a Stream
type StreamInfo struct {
	tempDir string
	handle  *TranscoderHandle
}

type StreamStatusManager struct {
	tempDir    string
	streamInfo map[string]StreamInfo
	transcoder Transcoder
}

func NewStreamStatusManager(tempDir string) StreamStatusManager {
	err := os.MkdirAll(tempDir, 0770)
	if err != nil {
		log.Fatal("Cannot create Temp-Dir %s", tempDir)
	}

	return StreamStatusManager{
		tempDir,
		make(map[string]StreamInfo),
		NewTranscoder(),
	}
}

type StreamStatus int

const (
	NoStream = iota
	StreamTranscodingFailed
	StreamInPreparation
	StreamReady
	TranscodingFinished
)

func (manager StreamStatusManager) GetStreamStatus(calculatedPath string) StreamStatus {
	info, hasInfo := manager.streamInfo[calculatedPath]

	if hasInfo {
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

func (manager StreamStatusManager) GetStreamTempdir(calculatedPath string) string {
	info, hasInfo := manager.streamInfo[calculatedPath]

	if hasInfo {
		return info.tempDir
	} else {
		return ""
	}
}

func (manager StreamStatusManager) StartStream(calculatedPath string) {
	_, hasInfo := manager.streamInfo[calculatedPath]
	if hasInfo {
		log.Printf("%s: Stream already active", calculatedPath)
		return
	}

	tempDir := manager.createTempDir(calculatedPath)

	log.Printf("%s: Starting Stream-Transcoder into %s", calculatedPath, tempDir)
	handle := manager.transcoder.StartTranscoder(calculatedPath, tempDir)

	manager.streamInfo[calculatedPath] = StreamInfo{
		tempDir,
		handle,
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
