package main

import (
	"net/http"

	"github.com/thoas/go-funk"
)

type RequestType int

const (
	DirectoryIndexRequest RequestType = iota
	FileRequest
	PlayerRequest
	StreamStatusRequest
	StreamPlaylistRequest
	StreamSegmentRequest
)

type RequestClassifier struct {
	transcodeExtensions []string
	playerExtensions    []string
}

func NewRequestClassifier(transcodeExtensions []string, playerExtensions []string) RequestClassifier {
	return RequestClassifier{
		transcodeExtensions,
		playerExtensions,
	}
}

func (requestClassifier *RequestClassifier) ClassifyRequest(request *http.Request, mappingResult PathMappingResult) RequestType {
	if mappingResult.FileInfo.IsDir() {
		return DirectoryIndexRequest
	} else {
		isTranscodeExtension := funk.Contains(requestClassifier.transcodeExtensions, mappingResult.FileExtension)
		isStreamRequest := request.URL.Query()["stream"] != nil

		isPlayerExtension := funk.Contains(requestClassifier.playerExtensions, mappingResult.FileExtension)
		isPlayerRequest := request.URL.Query()["play"] != nil

		if isTranscodeExtension && isStreamRequest {
			isPlaylistRequest := request.URL.Query()["playlist"] != nil
			isSegmentRequest := request.URL.Query()["segment"] != nil

			if isPlaylistRequest {
				return StreamPlaylistRequest
			} else if isSegmentRequest {
				return StreamSegmentRequest
			}

			return StreamStatusRequest
		} else if isPlayerExtension && isPlayerRequest {
			return PlayerRequest
		} else {
			return FileRequest
		}
	}
}

func (requestClassifier *RequestClassifier) GetSegmentFilename(request *http.Request) string {
	return request.URL.Query().Get("segment")
}
