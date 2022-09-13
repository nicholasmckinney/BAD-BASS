//go:build windows
// +build windows

package browser

import (
	"Webphish/internal/edge"
	"Webphish/internal/w32"
	"encoding/json"
	"errors"
	jchv "github.com/jchv/go-webview2/pkg/edge"
	"golang.org/x/sys/windows"
	"log"
	"reflect"
	"sync"
	"unsafe"
)

type LONG int32

var (
	windowContext     = map[uintptr]interface{}{}
	windowContextSync sync.RWMutex
)

func getWindowContext(wnd uintptr) interface{} {
	windowContextSync.RLock()
	defer windowContextSync.RUnlock()
	return windowContext[wnd]
}

func setWindowContext(wnd uintptr, data interface{}) {
	windowContextSync.Lock()
	defer windowContextSync.Unlock()
	windowContext[wnd] = data
}

type WindowRect struct {
	TopLeftX     LONG
	TopLeftY     LONG
	BottomRightX LONG
	BottomRightY LONG
}

type CenterWndParams struct {
	TopLeftX LONG
	TopLeftY LONG
	Width    LONG
	Height   LONG
}

type browser interface {
	Close()
	Environment() *edge.ICoreWebView2Environment
	Embed(hwnd uintptr) bool
	Resize()
	Navigate(url string)
	AddWebResourceRequestedFilter(filter string, ctx jchv.COREWEBVIEW2_WEB_RESOURCE_CONTEXT)
	Init(script string)
	Eval(script string)
	NotifyParentWindowPositionChanged() error
	Focus()
}

type webview struct {
	hwnd        uintptr
	mainthread  uintptr
	browser     browser
	parent      windows.HWND
	autofocus   bool
	maxsz       w32.Point
	minsz       w32.Point
	m           sync.Mutex
	bindings    map[string]interface{}
	dispatchq   []func()
	webCallback edge.WebResourceRequestedCallback
	msgCallback edge.WebMessageCallback
}

type WindowOptions struct {
	Title  string
	Width  uint
	Height uint
	IconId uint
	Center bool
	Parent windows.HWND
}

type WebViewOptions struct {
	Window unsafe.Pointer
	Debug  bool
	Parent windows.HWND

	// DataPath specifies the datapath for the WebView2 runtime to use for the
	// browser instance.
	DataPath string

	// AutoFocus will try to keep the WebView2 widget focused when the window
	// is focused.
	AutoFocus bool

	// WindowOptions customizes the window that is created to embed the
	// WebView2 widget.
	WindowOptions WindowOptions

	WebRequestCallback edge.WebResourceRequestedCallback
	WebMessageCallback edge.WebMessageCallback
}

// New creates a new webview in a new window.
func New(debug bool) edge.WebView { return NewWithOptions(WebViewOptions{Debug: debug}) }

// NewWindow creates a new webview using an existing window.
//
// Deprecated: Use NewWithOptions.
func NewWindow(debug bool, window unsafe.Pointer) edge.WebView {
	return NewWithOptions(WebViewOptions{Debug: debug, Window: window})
}

// NewWithOptions creates a new webview using the provided options.
func NewWithOptions(options WebViewOptions) edge.WebView {
	w := &webview{}
	w.bindings = map[string]interface{}{}
	w.autofocus = options.AutoFocus
	// callback for out of browser package. used by lib user
	w.webCallback = options.WebRequestCallback
	w.msgCallback = options.WebMessageCallback

	chromium := edge.NewChromium()
	chromium.MessageCallback = w.msgcb
	// sets callback from underlying webkit to this object. Users of this object can optionally
	// set their own callback by using SetWebResourceRequestedCallback
	chromium.WebResourceRequestedCallback = w.webResourceRequested

	chromium.DataPath = options.DataPath
	chromium.SetPermission(edge.CoreWebView2PermissionKindClipboardRead, edge.CoreWebView2PermissionStateAllow)

	w.browser = chromium
	w.mainthread, _, _ = w32.Kernel32GetCurrentThreadID.Call()
	if !w.CreateWithOptions(options.WindowOptions) {
		return nil
	}

	settings, err := chromium.GetSettings()
	if err != nil {
		log.Fatal(err)
	}

	settings.PutIsScriptEnabled(true)

	// disable context menu
	err = settings.PutAreDefaultContextMenusEnabled(options.Debug)
	if err != nil {
		log.Fatal(err)
	}
	// disable developer tools
	err = settings.PutAreDevToolsEnabled(options.Debug)
	if err != nil {
		log.Fatal(err)
	}

	return w
}

type rpcMessage struct {
	ID     int               `json:"id"`
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

func jsString(v interface{}) string { b, _ := json.Marshal(v); return string(b) }

func (w *webview) Close() {
	w.browser.Close()
}

func (w *webview) webResourceRequested(request *edge.ICoreWebView2WebResourceRequest, args *edge.ICoreWebView2WebResourceRequestedEventArgs) {
	var uri string
	var err error
	var env *edge.ICoreWebView2Environment
	var resp *jchv.ICoreWebView2WebResourceResponse

	if uri, err = request.GetUri(); err != nil {
	}
	env = w.browser.Environment()
	if w.webCallback != nil {
		callbackResponse := w.webCallback(uri)
		resp, err = env.CreateWebResourceResponse(callbackResponse.Content, callbackResponse.StatusCode, callbackResponse.ReasonPhrase, callbackResponse.Headers)
		if err != nil {
			return
		}
		if err = args.PutResponse(resp); err != nil {
			//console.MessageBoxPlain("Failure", "Failed to put response")
			return
		}
	}
}

func (w *webview) msgcb(msg string) {
	//console.MessageBoxPlain("Message Received", msg)
	if w.msgCallback != nil {
		ready := make(chan bool, 1) // TODO: dirty. i dirtied the ApplicationOutputWriter interface to add a signal to allow for setup of named pipe server before allowing client connections. bad alternative to a worse alternative: using a sleep interval that could lead to weird race conditions
		w.msgCallback(msg, ready)
	}
}

func (w *webview) callbinding(d rpcMessage) (interface{}, error) {
	w.m.Lock()
	f, ok := w.bindings[d.Method]
	w.m.Unlock()
	if !ok {
		return nil, nil
	}

	v := reflect.ValueOf(f)
	isVariadic := v.Type().IsVariadic()
	numIn := v.Type().NumIn()
	if (isVariadic && len(d.Params) < numIn-1) || (!isVariadic && len(d.Params) != numIn) {
		return nil, errors.New("function arguments mismatch")
	}
	args := []reflect.Value{}
	for i := range d.Params {
		var arg reflect.Value
		if isVariadic && i >= numIn-1 {
			arg = reflect.New(v.Type().In(numIn - 1).Elem())
		} else {
			arg = reflect.New(v.Type().In(i))
		}
		if err := json.Unmarshal(d.Params[i], arg.Interface()); err != nil {
			return nil, err
		}
		args = append(args, arg.Elem())
	}

	errorType := reflect.TypeOf((*error)(nil)).Elem()
	res := v.Call(args)
	switch len(res) {
	case 0:
		// No results from the function, just return nil
		return nil, nil

	case 1:
		// One result may be a value, or an error
		if res[0].Type().Implements(errorType) {
			if res[0].Interface() != nil {
				return nil, res[0].Interface().(error)
			}
			return nil, nil
		}
		return res[0].Interface(), nil

	case 2:
		// Two results: first one is value, second is error
		if !res[1].Type().Implements(errorType) {
			return nil, errors.New("second return value must be an error")
		}
		if res[1].Interface() == nil {
			return res[0].Interface(), nil
		}
		return res[0].Interface(), res[1].Interface().(error)

	default:
		return nil, errors.New("unexpected number of return values")
	}
}

func wndproc(hwnd, msg, wp, lp uintptr) uintptr {
	if w, ok := getWindowContext(hwnd).(*webview); ok {
		switch msg {
		case w32.WMMove, w32.WMMoving:
			_ = w.browser.NotifyParentWindowPositionChanged()
		case w32.WMNCLButtonDown:
			_, _, _ = w32.User32SetFocus.Call(w.hwnd)
			r, _, _ := w32.User32DefWindowProcW.Call(hwnd, msg, wp, lp)
			return r
		case w32.WMSize:
			w.browser.Resize()
		case w32.WMActivate:
			if wp == w32.WAInactive {
				break
			}
			if w.autofocus {
				w.browser.Focus()
			}
		case w32.WMClose:
			_, _, _ = w32.User32DestroyWindow.Call(hwnd)
		case w32.WMDestroy:
			w.Terminate()
		case w32.WMGetMinMaxInfo:
			lpmmi := (*w32.MinMaxInfo)(unsafe.Pointer(lp))
			if w.maxsz.X > 0 && w.maxsz.Y > 0 {
				lpmmi.PtMaxSize = w.maxsz
				lpmmi.PtMaxTrackSize = w.maxsz
			}
			if w.minsz.X > 0 && w.minsz.Y > 0 {
				lpmmi.PtMinTrackSize = w.minsz
			}
		default:
			r, _, _ := w32.User32DefWindowProcW.Call(hwnd, msg, wp, lp)
			return r
		}
		return 0
	}
	r, _, _ := w32.User32DefWindowProcW.Call(hwnd, msg, wp, lp)
	return r
}

func (w *webview) Create(debug bool, window unsafe.Pointer) bool {
	// This function signature stopped making sense a long time ago.
	// It is but legacy cruft at this point.
	return w.CreateWithOptions(WindowOptions{})
}

func GetWindowCenterParameters(hwnd windows.HWND) (*CenterWndParams, bool) {
	parentRect := WindowRect{}
	success, _, _ := w32.User32GetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&parentRect)))
	if int32(success) <= 0 {
		//console.MessageBoxPlain("Error", "Could not get parent rect")
		return nil, false
	}

	parentWidth := parentRect.BottomRightX - parentRect.TopLeftX
	parentHeight := parentRect.BottomRightY - parentRect.TopLeftY

	//windowWidth := parentWidth / 2
	//windowHeight := parentHeight / 2

	//posX := uint(parentRect.TopLeftX) + ((uint(parentWidth - windowWidth)) / 2)
	//posY := uint(parentRect.TopLeftY) + (uint(parentHeight-windowHeight))/2

	windowWidth := parentWidth
	windowHeight := parentHeight
	posX := parentRect.TopLeftX
	posY := parentRect.TopLeftY

	return &CenterWndParams{
		TopLeftX: LONG(posX),
		TopLeftY: LONG(posY),
		Width:    windowWidth,
		Height:   windowHeight,
	}, true
}

func (w *webview) CreateWithOptions(opts WindowOptions) bool {
	var hinstance windows.Handle
	_ = windows.GetModuleHandleEx(0, nil, &hinstance)

	var icon uintptr
	if opts.IconId == 0 {
		// load default icon
		icow, _, _ := w32.User32GetSystemMetrics.Call(w32.SystemMetricsCxIcon)
		icoh, _, _ := w32.User32GetSystemMetrics.Call(w32.SystemMetricsCyIcon)
		icon, _, _ = w32.User32LoadImageW.Call(uintptr(hinstance), 32512, icow, icoh, 0)
	} else {
		// load icon from resource
		icon, _, _ = w32.User32LoadImageW.Call(uintptr(hinstance), uintptr(opts.IconId), 1, 0, 0, w32.LR_DEFAULTSIZE|w32.LR_SHARED)
	}

	className, _ := windows.UTF16PtrFromString("webview")
	wc := w32.WndClassExW{
		CbSize:        uint32(unsafe.Sizeof(w32.WndClassExW{})),
		HInstance:     hinstance,
		LpszClassName: className,
		HIcon:         windows.Handle(icon),
		HIconSm:       windows.Handle(icon),
		LpfnWndProc:   windows.NewCallback(wndproc),
	}
	_, _, _ = w32.User32RegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	windowName, _ := windows.UTF16PtrFromString(opts.Title)

	centerParams, success := GetWindowCenterParameters(opts.Parent)
	if !success {
		//console.MessageBoxPlain("Error", "Failed to get size/coord parameters from parent")
		return false
	}

	//console.MessageBoxPlain("Creating window", "Creating window...")
	w.hwnd, _, _ = w32.User32CreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		w32.WSPopup,
		uintptr(centerParams.TopLeftX),
		uintptr(centerParams.TopLeftY),
		uintptr(centerParams.Width),
		uintptr(centerParams.Height),
		uintptr(opts.Parent),
		0,
		uintptr(hinstance),
		0,
	)
	setWindowContext(w.hwnd, w)
	//console.MessageBoxPlain("Created window", "Created window!!")
	_, _, _ = w32.User32ShowWindow.Call(w.hwnd, w32.SWShow)
	_, _, _ = w32.User32UpdateWindow.Call(w.hwnd)
	_, _, _ = w32.User32SetFocus.Call(w.hwnd)

	if !w.browser.Embed(w.hwnd) {
		return false
	}
	w.browser.Resize()
	w.parent = opts.Parent
	return true
}

func (w *webview) Destroy() {
}

func (w *webview) Run() {
	var msg w32.Msg
	for {
		_, _, _ = w32.User32GetMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0,
			0,
			0,
		)
		if msg.Message == w32.WMApp {
			w.m.Lock()
			q := append([]func(){}, w.dispatchq...)
			w.dispatchq = []func(){}
			w.m.Unlock()
			for _, v := range q {
				v()
			}
		} else if msg.Message == w32.WMQuit {
			return
		}
		r, _, _ := w32.User32GetAncestor.Call(uintptr(msg.Hwnd), w32.GARoot)
		r, _, _ = w32.User32IsDialogMessage.Call(r, uintptr(unsafe.Pointer(&msg)))
		if r != 0 {
			continue
		}
		_, _, _ = w32.User32TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		_, _, _ = w32.User32DispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))

	}
}

func (w *webview) Terminate() {
	_, _, _ = w32.User32PostQuitMessage.Call(0)
}

func (w *webview) Window() unsafe.Pointer {
	return unsafe.Pointer(w.hwnd)
}

func (w *webview) Navigate(url string) {
	w.browser.Navigate(url)
}

func (w *webview) AddWebResourceRequestedFilter(uri string) {
	w.browser.AddWebResourceRequestedFilter(uri, jchv.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)
}

func (w *webview) SetTitle(title string) {
	_title, err := windows.UTF16FromString(title)
	if err != nil {
		_title, _ = windows.UTF16FromString("")
	}
	_, _, _ = w32.User32SetWindowTextW.Call(w.hwnd, uintptr(unsafe.Pointer(&_title[0])))
}

func (view *webview) SetSize(w int, h int, hint edge.Hint) {
	index := w32.GWLStyle
	style, _, _ := w32.User32GetWindowLongPtrW.Call(view.hwnd, uintptr(index))
	_, _, _ = w32.User32SetWindowLongPtrW.Call(view.hwnd, uintptr(index), style)

	r := w32.Rect{}
	r.Left = 0
	r.Top = 0
	r.Right = int32(w)
	r.Bottom = int32(h)
	_, _, _ = w32.User32AdjustWindowRect.Call(uintptr(unsafe.Pointer(&r)), w32.WSOverlappedWindow, 0)
	_, _, _ = w32.User32SetWindowPos.Call(
		view.hwnd, 0, uintptr(r.Left), uintptr(r.Top), uintptr(r.Right-r.Left), uintptr(r.Bottom-r.Top),
		w32.SWPNoZOrder|w32.SWPNoActivate|w32.SWPNoMove|w32.SWPFrameChanged)
	view.browser.Resize()

}

func (w *webview) Init(js string) {
	w.browser.Init(js)
}

func (w *webview) Eval(js string) {
	w.browser.Eval(js)
}

func (w *webview) Dispatch(f func()) {
	w.m.Lock()
	w.dispatchq = append(w.dispatchq, f)
	w.m.Unlock()
	_, _, _ = w32.User32PostThreadMessageW.Call(w.mainthread, w32.WMApp, 0, 0)
}

func (w *webview) Bind(name string, f interface{}) error {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return errors.New("only functions can be bound")
	}
	if n := v.Type().NumOut(); n > 2 {
		return errors.New("function may only return a value or a value+error")
	}
	w.m.Lock()
	w.bindings[name] = f
	w.m.Unlock()

	w.Init("(function() { var name = " + jsString(name) + ";" + `
		var RPC = window._rpc = (window._rpc || {nextSeq: 1});
		window[name] = function() {
		  var seq = RPC.nextSeq++;
		  var promise = new Promise(function(resolve, reject) {
			RPC[seq] = {
			  resolve: resolve,
			  reject: reject,
			};
		  });
		  window.external.invoke(JSON.stringify({
			id: seq,
			method: name,
			params: Array.prototype.slice.call(arguments),
		  }));
		  return promise;
		}
	})()`)

	return nil
}
