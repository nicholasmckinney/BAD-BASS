package edge

import jchv "github.com/jchv/go-webview2/pkg/edge"

type _ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke jchv.ComProc
}

type ICoreWebView2AcceleratorKeyPressedEventHandler struct {
	vtbl *_ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl
	impl _ICoreWebView2AcceleratorKeyPressedEventHandlerImpl
}

func (i *ICoreWebView2AcceleratorKeyPressedEventHandler) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call()
	return r
}
func _ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownQueryInterface(this *ICoreWebView2AcceleratorKeyPressedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownAddRef(this *ICoreWebView2AcceleratorKeyPressedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownRelease(this *ICoreWebView2AcceleratorKeyPressedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2AcceleratorKeyPressedEventHandlerInvoke(this *ICoreWebView2AcceleratorKeyPressedEventHandler, sender *ICoreWebView2Controller, args *ICoreWebView2AcceleratorKeyPressedEventArgs) uintptr {
	return this.impl.AcceleratorKeyPressed(sender, args)
}

type _ICoreWebView2AcceleratorKeyPressedEventHandlerImpl interface {
	_IUnknownImpl
	AcceleratorKeyPressed(sender *ICoreWebView2Controller, args *ICoreWebView2AcceleratorKeyPressedEventArgs) uintptr
}

var _ICoreWebView2AcceleratorKeyPressedEventHandlerFn = _ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl{
	_IUnknownVtbl{
		jchv.NewComProc(_ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownQueryInterface),
		jchv.NewComProc(_ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownAddRef),
		jchv.NewComProc(_ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownRelease),
	},
	jchv.NewComProc(_ICoreWebView2AcceleratorKeyPressedEventHandlerInvoke),
}

func newICoreWebView2AcceleratorKeyPressedEventHandler(impl _ICoreWebView2AcceleratorKeyPressedEventHandlerImpl) *ICoreWebView2AcceleratorKeyPressedEventHandler {
	return &ICoreWebView2AcceleratorKeyPressedEventHandler{
		vtbl: &_ICoreWebView2AcceleratorKeyPressedEventHandlerFn,
		impl: impl,
	}
}
