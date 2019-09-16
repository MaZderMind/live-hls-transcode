package main

import (
	"github.com/thoas/go-funk"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type DirectoryIndex struct {
	template            *template.Template
	streamingExtensions []string
}

type TemplateFileDto struct {
	Name     string
	IsDir    bool
	Size     int64
	IsStream bool
}

func NewDirectoryIndex(streamingExtensions []string) DirectoryIndex {
	templateFile, err := template.New("directory-index.gohtml").ParseFiles("templates/directory-index.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	return DirectoryIndex{
		templateFile,
		streamingExtensions,
	}
}

func (directoryIndex DirectoryIndex) Handle(writer http.ResponseWriter, request *http.Request, mappingResult PathMappingResult) {
	files, err := ioutil.ReadDir(mappingResult.CalculatedPath)
	if err != nil {
		log.Printf("Error reading dir: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	if err = directoryIndex.template.Execute(writer, struct {
		IsRoot bool
		Files  []TemplateFileDto
	}{
		path.Clean(request.URL.Path) == "/",
		directoryIndex.buildTemplateFileDtos(files),
	}); err != nil {
		log.Printf("Template-Formatting failed: %s", err)
	}
}

func (directoryIndex DirectoryIndex) buildTemplateFileDtos(fileInfos []os.FileInfo) []TemplateFileDto {
	templateFiles := make([]TemplateFileDto, len(fileInfos))
	for i, fileInfo := range fileInfos {
		extension := strings.ToLower(strings.TrimLeft(filepath.Ext(fileInfo.Name()), "."))

		templateFiles[i] = TemplateFileDto{
			fileInfo.Name(),
			fileInfo.IsDir(),
			fileInfo.Size(),
			funk.ContainsString(directoryIndex.streamingExtensions, extension),
		}
	}

	return templateFiles
}
