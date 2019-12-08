package main

import (
	"facette.io/natsort"
	"github.com/dustin/go-humanize"
	"github.com/thoas/go-funk"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type DirectoryIndex struct {
	template            *template.Template
	streamingExtensions []string
}

type TemplateFileDto struct {
	Name     string
	IsDir    bool
	Size     string
	IsStream bool
}

func NewDirectoryIndex(streamingExtensions []string) DirectoryIndex {
	return DirectoryIndex{
		readTemplate("directory-index.gohtml"),
		streamingExtensions,
	}
}

func (directoryIndex *DirectoryIndex) Handle(writer http.ResponseWriter, request *http.Request, mappingResult PathMappingResult) {
	directoryIndex.redirectPathsWithoutSlash(writer, request)

	files, err := ioutil.ReadDir(mappingResult.CalculatedPath)
	if err != nil {
		log.Printf("Error reading dir: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	if err = directoryIndex.template.Execute(writer, struct {
		IsRoot  bool
		UrlPath string
		Files   []TemplateFileDto
	}{
		path.Clean(request.URL.Path) == "/",
		mappingResult.UrlPath,
		directoryIndex.buildTemplateFileDtos(files),
	}); err != nil {
		log.Printf("Template-Formatting failed: %s", err)
	}
}

func (directoryIndex *DirectoryIndex) buildTemplateFileDtos(fileInfos []os.FileInfo) []TemplateFileDto {
	sortByNameDirectoriesFirst(fileInfos)

	templateFiles := make([]TemplateFileDto, 0)
	for _, fileInfo := range fileInfos {
		if fileInfo.Name()[0] == '.' {
			continue
		}

		extension := strings.ToLower(strings.TrimLeft(filepath.Ext(fileInfo.Name()), "."))

		templateFiles = append(templateFiles, TemplateFileDto{
			fileInfo.Name(),
			fileInfo.IsDir(),
			humanize.Bytes(uint64(fileInfo.Size())),
			funk.ContainsString(directoryIndex.streamingExtensions, extension),
		})
	}

	return templateFiles
}

func sortByNameDirectoriesFirst(fileInfos []os.FileInfo) {
	sort.Slice(fileInfos, func(aIndex, bIndex int) bool {
		a := fileInfos[aIndex]
		b := fileInfos[bIndex]

		if a.IsDir() && !b.IsDir() {
			return true
		} else if !a.IsDir() && b.IsDir() {
			return false
		} else {
			return natsort.Compare(
				strings.ToLower(a.Name()),
				strings.ToLower(b.Name()))
		}
	})
}

func (directoryIndex *DirectoryIndex) redirectPathsWithoutSlash(writer http.ResponseWriter, request *http.Request) {
	requestPath := request.URL.Path
	if !strings.HasSuffix(requestPath, "/") {
		http.Redirect(writer, request, requestPath+"/", http.StatusSeeOther)
	}
}
