package edge

import (
	jchv "github.com/jchv/go-webview2/pkg/edge"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2WebResourceRequestVtbl struct {
	_IUnknownVtbl
	GetUri     jchv.ComProc
	PutUri     jchv.ComProc
	GetMethod  jchv.ComProc
	PutMethod  jchv.ComProc
	GetContent jchv.ComProc
	PutContent jchv.ComProc
	GetHeaders jchv.ComProc
}

type ICoreWebView2WebResourceRequest struct {
	vtbl *_ICoreWebView2WebResourceRequestVtbl
}

func (i *ICoreWebView2WebResourceRequest) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call()
	return r
}

func (i *ICoreWebView2WebResourceRequest) GetUri() (string, error) {
	var err error
	// Create *uint16 to hold result
	var _uri *uint16
	_, _, err = i.vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_uri)),
	)
	if err != windows.ERROR_SUCCESS {
		return "", err
	} // Get result and cleanup
	uri := windows.UTF16PtrToString(_uri)
	windows.CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, nil
}
