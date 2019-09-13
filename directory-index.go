package main

import (
	"fmt"
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
	rootDir             string
	template            *template.Template
	streamingExtensions []string
}

type TemplateFile struct {
	Name     string
	IsDir    bool
	Size     int64
	IsStream bool
}

func NewDirectoryIndex(rootDir string, streamingExtensions []string) DirectoryIndex {
	templateFile, err := template.New("directory-index.gohtml").ParseFiles("templates/directory-index.gohtml")

	if err != nil {
		log.Fatal(err)
	}
	return DirectoryIndex{
		rootDir,
		templateFile,
		streamingExtensions,
	}
}

func (directoryIndex DirectoryIndex) Handle(writer http.ResponseWriter, request *http.Request, fileHandler func(calculatedPath string)) {
	calculatedPath := directoryIndex.calculatePath(request.URL.Path)

	info, err := os.Stat(calculatedPath)
	if os.IsNotExist(err) {
		writer.WriteHeader(404)
		return
	} else if os.IsPermission(err) {
		writer.WriteHeader(401)
		return
	} else if !info.IsDir() {
		fileHandler(calculatedPath)
		return;
	}

	files, err := ioutil.ReadDir(calculatedPath)
	if err != nil {
		log.Printf("Error reading dir: %s", err)
		writer.WriteHeader(500)
		return
	}

	writer.Header().Add("Content-Type", "text/html")
	if err = directoryIndex.template.Execute(writer, struct {
		IsRoot bool
		Files  []TemplateFile
	}{
		path.Clean(request.URL.Path) == "/",
		directoryIndex.buildTemplateFile(files),
	}); err != nil {
		fmt.Printf("Template-Formatting failed: %s", err)
	}
}

func (directoryIndex DirectoryIndex) buildTemplateFile(fileInfos []os.FileInfo) []TemplateFile {
	templateFiles := make([]TemplateFile, len(fileInfos))
	for i, fileInfo := range fileInfos {
		extension := strings.ToLower(strings.TrimLeft(filepath.Ext(fileInfo.Name()), "."))

		templateFiles[i] = TemplateFile{
			fileInfo.Name(),
			fileInfo.IsDir(),
			fileInfo.Size(),
			funk.ContainsString(directoryIndex.streamingExtensions, extension),
		}
	}

	return templateFiles
}

func (directoryIndex DirectoryIndex) calculatePath(urlPath string) string {
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}

	cleanPath := path.Clean(urlPath)
	return path.Join(directoryIndex.rootDir, cleanPath)
}
