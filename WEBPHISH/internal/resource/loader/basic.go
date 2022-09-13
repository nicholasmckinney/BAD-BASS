package loader

import (
	"Webphish/internal"
	"Webphish/internal/edge"
	"Webphish/internal/resource"
	"bytes"
	"encoding/xml"
	"path/filepath"
	"regexp"
	"strings"
)

type BasicLoader struct {
	webInjectConfiguration resource.Configuration
	resourceMap            map[string]string
}

const URLPrefix = "http://contoso.com/"

func NewBasicLoader() *BasicLoader {
	return &BasicLoader{
		resourceMap: make(map[string]string),
	}
}

func (loader *BasicLoader) Get(fileName string) edge.WebResponse {
	fileName = strings.TrimPrefix(fileName, URLPrefix)
	fileExt := filepath.Ext(fileName)
	contentType := edge.GetContentType(fileExt)
	fileContent := loader.resourceMap[fileName]
	buf := bytes.NewBufferString(fileContent)
	return edge.OK(buf.Bytes(), contentType)
}

func (loader *BasicLoader) HasMatchingResources(title string) bool {
	_, found := loader.MatchingResource(title)
	return found == internal.GenericSuccess
}

func (loader *BasicLoader) MatchingResource(title string) (string, internal.ErrorCode) {
	for _, matcher := range loader.webInjectConfiguration.ResourceMatchers {
		matched, _ := regexp.MatchString(matcher.Name, title)
		if matched {
			return matcher.MatchingResourcePath, internal.GenericSuccess
		}
	}
	return "", internal.GenericError
}

func (loader *BasicLoader) Put(file *resource.File) internal.ErrorCode {
	if file.Path == "conf.xml" {
		err := xml.Unmarshal(file.Content, &loader.webInjectConfiguration)
		if err != nil {
			return internal.ERR_LOADER_LOAD_CONFIGURATION
		}
		return internal.GenericSuccess
	}
	if file.IsDirectory {
		return internal.GenericSuccess
	}
	file.Path = strings.ReplaceAll(file.Path, "\\", "/")
	buf := bytes.NewBuffer(file.Content)
	loader.resourceMap[file.Path] = buf.String()
	return internal.GenericSuccess
}
