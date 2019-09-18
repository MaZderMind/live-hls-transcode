package main

import (
	"github.com/gosimple/slug"
	"io/ioutil"
	"log"
	"os"
)

// Stored Information about a Stream
type StreamInfo struct {
	CalculatedPath string
	UrlPath        string
	TempDir        string
	Handle         *TranscoderHandle
}

func (info StreamInfo) DominantStatusCode() StreamStatus {
	if info.Handle != nil {
		if info.Handle.IsFinished() {
			return TranscodingFinished
		} else if info.Handle.IsReady() {
			return StreamReady
		} else if info.Handle.IsRunning() {
			return StreamInPreparation
		} else {
			return StreamTranscodingFailed
		}
	} else {
		return NoStream
	}
}

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
		if info.Handle.IsFinished() {
			log.Printf("%s: Stream-Transcoder already finished, ignoring Stop-Command", calculatedPath)
		} else {
			log.Printf("%s: Stopping unfinished Stream-Transcoder", calculatedPath)
			info.Handle.Stop()

			delete(manager.streamInfo, calculatedPath)
		}
	}
}

func (manager StreamStatusManager) OtherStreamInfos(thisCalculatedPath string) []StreamInfo {
	otherStreamInfos := make([]StreamInfo, 0)

	for _, streamInfo := range manager.streamInfo {
		if streamInfo.CalculatedPath == thisCalculatedPath {
			continue
		}

		otherStreamInfos = append(otherStreamInfos, streamInfo)
	}

	return otherStreamInfos
}
