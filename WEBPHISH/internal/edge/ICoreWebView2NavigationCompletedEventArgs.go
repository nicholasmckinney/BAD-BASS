package edge

import jchv "github.com/jchv/go-webview2/pkg/edge"

type _ICoreWebView2NavigationCompletedEventArgsVtbl struct {
	_IUnknownVtbl
	GetIsSuccess      jchv.ComProc
	GetWebErrorStatus jchv.ComProc
	GetNavigationId   jchv.ComProc
}

type ICoreWebView2NavigationCompletedEventArgs struct {
	vtbl *_ICoreWebView2NavigationCompletedEventArgsVtbl
}

func (i *ICoreWebView2NavigationCompletedEventArgs) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call()
	return r
}
