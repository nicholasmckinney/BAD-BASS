package edge

import jchv "github.com/jchv/go-webview2/pkg/edge"

type _ICoreWebView2NavigationCompletedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke jchv.ComProc
}

type ICoreWebView2NavigationCompletedEventHandler struct {
	vtbl *_ICoreWebView2NavigationCompletedEventHandlerVtbl
	impl _ICoreWebView2NavigationCompletedEventHandlerImpl
}

func _ICoreWebView2NavigationCompletedEventHandlerIUnknownQueryInterface(this *ICoreWebView2NavigationCompletedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2NavigationCompletedEventHandlerIUnknownAddRef(this *ICoreWebView2NavigationCompletedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2NavigationCompletedEventHandlerIUnknownRelease(this *ICoreWebView2NavigationCompletedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2NavigationCompletedEventHandlerInvoke(this *ICoreWebView2NavigationCompletedEventHandler, sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr {
	return this.impl.NavigationCompleted(sender, args)
}

type _ICoreWebView2NavigationCompletedEventHandlerImpl interface {
	_IUnknownImpl
	NavigationCompleted(sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr
}

var _ICoreWebView2NavigationCompletedEventHandlerFn = _ICoreWebView2NavigationCompletedEventHandlerVtbl{
	_IUnknownVtbl{
		jchv.NewComProc(_ICoreWebView2NavigationCompletedEventHandlerIUnknownQueryInterface),
		jchv.NewComProc(_ICoreWebView2NavigationCompletedEventHandlerIUnknownAddRef),
		jchv.NewComProc(_ICoreWebView2NavigationCompletedEventHandlerIUnknownRelease),
	},
	jchv.NewComProc(_ICoreWebView2NavigationCompletedEventHandlerInvoke),
}

func newICoreWebView2NavigationCompletedEventHandler(impl _ICoreWebView2NavigationCompletedEventHandlerImpl) *ICoreWebView2NavigationCompletedEventHandler {
	return &ICoreWebView2NavigationCompletedEventHandler{
		vtbl: &_ICoreWebView2NavigationCompletedEventHandlerFn,
		impl: impl,
	}
}
