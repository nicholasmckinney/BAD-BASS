package edge

import (
	"Webphish/internal"
	"fmt"
	"strings"
	"unsafe"
)

type WebResponse struct {
	Content      []byte
	StatusCode   int
	ReasonPhrase string
	Headers      string
}

type WebResourceRequestedCallback func(uri string) WebResponse
type WebMessageCallback func(msg string, ready chan bool) internal.ErrorCode

// This is copied from webview/webview.
// The documentation is included for convenience.

// Hint is used to configure window sizing and resizing behavior.
type Hint int

const (
	// HintNone specifies that width and height are default size
	HintNone Hint = iota

	// HintFixed specifies that window size can not be changed by a user
	HintFixed

	// HintMin specifies that width and height are minimum bounds
	HintMin

	// HintMax specifies that width and height are maximum bounds
	HintMax
)

type ContentType string

const (
	ContentTypeHTML       = "text/html"
	ContentTypeJavascript = "text/javascript"
	ContentTypePNG        = "image/png"
	ContentTypeJPEG       = "image/jpeg"
)

var (
	ContentTypes = map[string]string{
		"html": ContentTypeHTML,
		"js":   ContentTypeJavascript,
		"png":  ContentTypePNG,
		"jpg":  ContentTypeJPEG,
		"jpeg": ContentTypeJPEG,
	}
)

func GetContentType(ext string) ContentType {
	var contentType string
	var ok bool

	lowered := strings.ToLower(ext)
	if contentType, ok = ContentTypes[lowered]; !ok {
		return ContentTypeHTML
	}
	return ContentType(contentType)
}

// WebView is the interface for the webview.
type WebView interface {
	Close()

	// Run runs the main loop until it's terminated. After this function exits -
	// you must destroy the webview.
	Run()

	// Terminate stops the main loop. It is safe to call this function from
	// a background thread.
	Terminate()

	// Dispatch posts a function to be executed on the main thread. You normally
	// do not need to call this function, unless you want to tweak the native
	// window.
	Dispatch(f func())

	// Destroy destroys a webview and closes the native window.
	Destroy()

	// Window returns a native window handle pointer. When using GTK backend the
	// pointer is GtkWindow pointer, when using Cocoa backend the pointer is
	// NSWindow pointer, when using Win32 backend the pointer is HWND pointer.
	Window() unsafe.Pointer

	// SetTitle updates the title of the native window. Must be called from the UI
	// thread.
	SetTitle(title string)

	// SetSize updates native window size. See Hint constants.
	SetSize(w int, h int, hint Hint)

	// Navigate navigates webview to the given URL. URL may be a data URI, i.e.
	// "data:text/text,<html>...</html>". It is often ok not to url-encode it
	// properly, webview will re-encode it for you.
	Navigate(url string)

	AddWebResourceRequestedFilter(uri string)

	// Init injects JavaScript code at the initialization of the new page. Every
	// time the webview will open a the new page - this initialization code will
	// be executed. It is guaranteed that code is executed before window.onload.
	Init(js string)

	// Eval evaluates arbitrary JavaScript code. Evaluation happens asynchronously,
	// also the result of the expression is ignored. Use RPC bindings if you want
	// to receive notifications about the results of the evaluation.
	Eval(js string)

	// Bind binds a callback function so that it will appear under the given name
	// as a global JavaScript function. Internally it uses webview_init().
	// Callback receives a request string and a user-provided argument pointer.
	// Request string is a JSON array of all the arguments passed to the
	// JavaScript function.
	//
	// f must be a function
	// f must return either value and error or just error
	Bind(name string, f interface{}) error
}

func OK(content []byte, contentType ContentType) WebResponse {
	return WebResponse{
		Content:      content,
		StatusCode:   200,
		ReasonPhrase: "OK",
		Headers:      fmt.Sprintf("Content-Type: %s", contentType),
	}
}

func NotFound() WebResponse {
	return WebResponse{
		Content:      nil,
		StatusCode:   404,
		ReasonPhrase: "Not Found",
		Headers:      "",
	}
}
