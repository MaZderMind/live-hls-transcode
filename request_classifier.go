package main

import (
	"github.com/thoas/go-funk"
	"net/http"
)

type RequestType int

const (
	DirectoryIndexRequest RequestType = iota
	FileRequest
	StreamStatusRequest
	StreamPlaylistRequest
	StreamSegmentRequest
)

type RequestClassifier struct {
	streamableExtensions []string
}

func NewRequestClassifier(streamableExtensions []string) RequestClassifier {
	return RequestClassifier{
		streamableExtensions,
	}
}

func (requestClassifier *RequestClassifier) ClassifyRequest(request *http.Request, mappingResult PathMappingResult) RequestType {
	if mappingResult.FileInfo.IsDir() {
		return DirectoryIndexRequest
	} else {
		isStreamableExtension := funk.Contains(requestClassifier.streamableExtensions, mappingResult.FileExtension)
		isStreamRequest := request.URL.Query()["stream"] != nil

		if isStreamableExtension && isStreamRequest {
			isPlaylistRequest := request.URL.Query()["playlist"] != nil
			isSegmentRequest := request.URL.Query()["segment"] != nil

			if isPlaylistRequest {
				return StreamPlaylistRequest
			} else if isSegmentRequest {
				return StreamSegmentRequest
			}

			return StreamStatusRequest
		}

		return FileRequest
	}
}

func (requestClassifier *RequestClassifier) GetSegmentFilename(request *http.Request) string {
	return request.URL.Query().Get("segment")
}
