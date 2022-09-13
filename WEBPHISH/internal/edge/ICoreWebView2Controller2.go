package edge

import (
	jchv "github.com/jchv/go-webview2/pkg/edge"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2Controller2Vtbl struct {
	_IUnknownVtbl
	GetIsVisible                      jchv.ComProc
	PutIsVisible                      jchv.ComProc
	GetBounds                         jchv.ComProc
	PutBounds                         jchv.ComProc
	GetZoomFactor                     jchv.ComProc
	PutZoomFactor                     jchv.ComProc
	AddZoomFactorChanged              jchv.ComProc
	RemoveZoomFactorChanged           jchv.ComProc
	SetBoundsAndZoomFactor            jchv.ComProc
	MoveFocus                         jchv.ComProc
	AddMoveFocusRequested             jchv.ComProc
	RemoveMoveFocusRequested          jchv.ComProc
	AddGotFocus                       jchv.ComProc
	RemoveGotFocus                    jchv.ComProc
	AddLostFocus                      jchv.ComProc
	RemoveLostFocus                   jchv.ComProc
	AddAcceleratorKeyPressed          jchv.ComProc
	RemoveAcceleratorKeyPressed       jchv.ComProc
	GetParentWindow                   jchv.ComProc
	PutParentWindow                   jchv.ComProc
	NotifyParentWindowPositionChanged jchv.ComProc
	Close                             jchv.ComProc
	GetCoreWebView2                   jchv.ComProc
	GetDefaultBackgroundColor         jchv.ComProc
	PutDefaultBackgroundColor         jchv.ComProc
}

type ICoreWebView2Controller2 struct {
	vtbl *_ICoreWebView2Controller2Vtbl
}

func (i *ICoreWebView2Controller2) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call()
	return r
}

func (i *ICoreWebView2Controller2) GetDefaultBackgroundColor() (*jchv.COREWEBVIEW2_COLOR, error) {
	var err error
	var backgroundColor *jchv.COREWEBVIEW2_COLOR
	_, _, err = i.vtbl.GetDefaultBackgroundColor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&backgroundColor)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return backgroundColor, nil
}

func (i *ICoreWebView2Controller2) PutDefaultBackgroundColor(backgroundColor jchv.COREWEBVIEW2_COLOR) error {
	var err error

	// Cast to a uint32 as that's what the call is expecting
	col := *(*uint32)(unsafe.Pointer(&backgroundColor))

	_, _, err = i.vtbl.PutDefaultBackgroundColor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(col),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
