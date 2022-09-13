package edge

import (
	"unsafe"

	"Webphish/internal/w32"
	jchv "github.com/jchv/go-webview2/pkg/edge"
	"golang.org/x/sys/windows"
)

type _ICoreWebView2ControllerVtbl struct {
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
}

type ICoreWebView2Controller struct {
	vtbl *_ICoreWebView2ControllerVtbl
}

func (i *ICoreWebView2Controller) Close() error {
	_, _, err := i.vtbl.Close.Call()
	return err
}

func (i *ICoreWebView2Controller) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call()
	return r
}

func (i *ICoreWebView2Controller) GetBounds() (*w32.Rect, error) {
	var err error
	var bounds w32.Rect
	_, _, err = i.vtbl.GetBounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return &bounds, nil
}

func (i *ICoreWebView2Controller) PutBounds(bounds w32.Rect) error {
	var err error

	_, _, err = i.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Controller) AddAcceleratorKeyPressed(eventHandler *ICoreWebView2AcceleratorKeyPressedEventHandler, token *_EventRegistrationToken) error {
	var err error
	_, _, err = i.vtbl.AddAcceleratorKeyPressed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Controller) PutIsVisible(isVisible bool) error {
	var err error

	_, _, err = i.vtbl.PutIsVisible.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(isVisible)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Controller) GetICoreWebView2Controller2() *ICoreWebView2Controller2 {

	var result *ICoreWebView2Controller2

	iidICoreWebView2Controller2 := jchv.NewGUID("{c979903e-d4ca-4228-92eb-47ee3fa96eab}")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Controller2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Controller) NotifyParentWindowPositionChanged() error {
	var err error
	_, _, err = i.vtbl.NotifyParentWindowPositionChanged.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Controller) MoveFocus(reason uintptr) error {
	var err error
	_, _, err = i.vtbl.MoveFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(reason),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
