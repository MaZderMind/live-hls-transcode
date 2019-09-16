package main

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type PathMapper struct {
	rootDir string
}

type PathMappingResult struct {
	CalculatedPath string
	FileExtension  string
	FileInfo       os.FileInfo
	StatError      error
}

func (result PathMappingResult) HandleError(writer http.ResponseWriter) {
	if os.IsNotExist(result.StatError) {
		writer.WriteHeader(http.StatusNotFound)
	} else if os.IsPermission(result.StatError) {
		writer.WriteHeader(http.StatusForbidden)
	}
}

func NewPathMapper(rootDir string) PathMapper {
	return PathMapper{
		rootDir,
	}
}

func (pathMapper PathMapper) MapUrlPathToFilesystem(urlPath string) PathMappingResult {
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}

	cleanPath := path.Clean(urlPath)
	calculatedPath := path.Join(pathMapper.rootDir, cleanPath)
	fileInfo, statError := os.Stat(calculatedPath)

	return PathMappingResult{
		calculatedPath,
		strings.ToLower(strings.TrimLeft(filepath.Ext(calculatedPath), ".")),
		fileInfo,
		statError,
	}
}
