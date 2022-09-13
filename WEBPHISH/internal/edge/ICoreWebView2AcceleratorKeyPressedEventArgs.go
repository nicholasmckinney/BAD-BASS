package edge

import (
	jchv "github.com/jchv/go-webview2/pkg/edge"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2AcceleratorKeyPressedEventArgsVtbl struct {
	_IUnknownVtbl
	GetKeyEventKind      jchv.ComProc
	GetVirtualKey        jchv.ComProc
	GetKeyEventLParam    jchv.ComProc
	GetPhysicalKeyStatus jchv.ComProc
	GetHandled           jchv.ComProc
	PutHandled           jchv.ComProc
}

type ICoreWebView2AcceleratorKeyPressedEventArgs struct {
	vtbl *_ICoreWebView2AcceleratorKeyPressedEventArgsVtbl
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call()
	return r
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs) GetKeyEventKind() (jchv.COREWEBVIEW2_KEY_EVENT_KIND, error) {
	var err error
	var keyEventKind jchv.COREWEBVIEW2_KEY_EVENT_KIND
	_, _, err = i.vtbl.GetKeyEventKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&keyEventKind)),
	)
	if err != windows.ERROR_SUCCESS {
		return 0, err
	}
	return keyEventKind, nil
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs) GetVirtualKey() (uint, error) {
	var err error
	var virtualKey uint
	_, _, err = i.vtbl.GetVirtualKey.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&virtualKey)),
	)
	if err != windows.ERROR_SUCCESS {
		return 0, err
	}
	return virtualKey, nil
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs) GetPhysicalKeyStatus() (jchv.COREWEBVIEW2_PHYSICAL_KEY_STATUS, error) {
	var err error
	var physicalKeyStatus jchv.COREWEBVIEW2_PHYSICAL_KEY_STATUS
	_, _, err = i.vtbl.GetPhysicalKeyStatus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&physicalKeyStatus)),
	)
	if err != windows.ERROR_SUCCESS {
		return jchv.COREWEBVIEW2_PHYSICAL_KEY_STATUS{}, err
	}
	return physicalKeyStatus, nil
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs) PutHandled(handled bool) error {
	var err error

	_, _, err = i.vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(handled)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
