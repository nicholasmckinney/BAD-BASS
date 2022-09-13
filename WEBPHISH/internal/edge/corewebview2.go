//go:build windows
// +build windows

package edge

import (
	"log"
	"runtime"
	"unsafe"

	"Webphish/internal/w32"

	jchv "github.com/jchv/go-webview2/pkg/edge"
	"github.com/jchv/go-webview2/webviewloader"
	"golang.org/x/sys/windows"
)

func init() {
	runtime.LockOSThread()

	r, _, _ := w32.Ole32CoInitializeEx.Call(0, 2)
	if int(r) < 0 {
		log.Printf("Warning: CoInitializeEx call failed: E=%08x", r)
	}
}

type _EventRegistrationToken struct {
	Value int64
}

type CoreWebView2PermissionKind uint32

const (
	CoreWebView2PermissionKindUnknownPermission CoreWebView2PermissionKind = iota
	CoreWebView2PermissionKindMicrophone
	CoreWebView2PermissionKindCamera
	CoreWebView2PermissionKindGeolocation
	CoreWebView2PermissionKindNotifications
	CoreWebView2PermissionKindOtherSensors
	CoreWebView2PermissionKindClipboardRead
)

type CoreWebView2PermissionState uint32

const (
	CoreWebView2PermissionStateDefault CoreWebView2PermissionState = iota
	CoreWebView2PermissionStateAllow
	CoreWebView2PermissionStateDeny
)

func createCoreWebView2EnvironmentWithOptions(browserExecutableFolder, userDataFolder *uint16, environmentOptions uintptr, environmentCompletedHandle *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) (uintptr, error) {
	return webviewloader.CreateCoreWebView2EnvironmentWithOptions(
		browserExecutableFolder,
		userDataFolder,
		environmentOptions,
		uintptr(unsafe.Pointer(environmentCompletedHandle)),
	)
}

// IUnknown

type _IUnknownVtbl struct {
	QueryInterface jchv.ComProc
	AddRef         jchv.ComProc
	Release        jchv.ComProc
}

type _IUnknownImpl interface {
	QueryInterface(refiid, object uintptr) uintptr
	AddRef() uintptr
	Release() uintptr
}

// ICoreWebView2

type iCoreWebView2Vtbl struct {
	_IUnknownVtbl
	GetSettings                            jchv.ComProc
	GetSource                              jchv.ComProc
	Navigate                               jchv.ComProc
	NavigateToString                       jchv.ComProc
	AddNavigationStarting                  jchv.ComProc
	RemoveNavigationStarting               jchv.ComProc
	AddContentLoading                      jchv.ComProc
	RemoveContentLoading                   jchv.ComProc
	AddSourceChanged                       jchv.ComProc
	RemoveSourceChanged                    jchv.ComProc
	AddHistoryChanged                      jchv.ComProc
	RemoveHistoryChanged                   jchv.ComProc
	AddNavigationCompleted                 jchv.ComProc
	RemoveNavigationCompleted              jchv.ComProc
	AddFrameNavigationStarting             jchv.ComProc
	RemoveFrameNavigationStarting          jchv.ComProc
	AddFrameNavigationCompleted            jchv.ComProc
	RemoveFrameNavigationCompleted         jchv.ComProc
	AddScriptDialogOpening                 jchv.ComProc
	RemoveScriptDialogOpening              jchv.ComProc
	AddPermissionRequested                 jchv.ComProc
	RemovePermissionRequested              jchv.ComProc
	AddProcessFailed                       jchv.ComProc
	RemoveProcessFailed                    jchv.ComProc
	AddScriptToExecuteOnDocumentCreated    jchv.ComProc
	RemoveScriptToExecuteOnDocumentCreated jchv.ComProc
	ExecuteScript                          jchv.ComProc
	CapturePreview                         jchv.ComProc
	Reload                                 jchv.ComProc
	PostWebMessageAsJSON                   jchv.ComProc
	PostWebMessageAsString                 jchv.ComProc
	AddWebMessageReceived                  jchv.ComProc
	RemoveWebMessageReceived               jchv.ComProc
	CallDevToolsProtocolMethod             jchv.ComProc
	GetBrowserProcessID                    jchv.ComProc
	GetCanGoBack                           jchv.ComProc
	GetCanGoForward                        jchv.ComProc
	GoBack                                 jchv.ComProc
	GoForward                              jchv.ComProc
	GetDevToolsProtocolEventReceiver       jchv.ComProc
	Stop                                   jchv.ComProc
	AddNewWindowRequested                  jchv.ComProc
	RemoveNewWindowRequested               jchv.ComProc
	AddDocumentTitleChanged                jchv.ComProc
	RemoveDocumentTitleChanged             jchv.ComProc
	GetDocumentTitle                       jchv.ComProc
	AddHostObjectToScript                  jchv.ComProc
	RemoveHostObjectFromScript             jchv.ComProc
	OpenDevToolsWindow                     jchv.ComProc
	AddContainsFullScreenElementChanged    jchv.ComProc
	RemoveContainsFullScreenElementChanged jchv.ComProc
	GetContainsFullScreenElement           jchv.ComProc
	AddWebResourceRequested                jchv.ComProc
	RemoveWebResourceRequested             jchv.ComProc
	AddWebResourceRequestedFilter          jchv.ComProc
	RemoveWebResourceRequestedFilter       jchv.ComProc
	AddWindowCloseRequested                jchv.ComProc
	RemoveWindowCloseRequested             jchv.ComProc
}

type ICoreWebView2 struct {
	vtbl *iCoreWebView2Vtbl
}

func (i *ICoreWebView2) GetSettings() (*jchv.ICoreWebViewSettings, error) {
	var err error
	var settings *jchv.ICoreWebViewSettings
	_, _, err = i.vtbl.GetSettings.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&settings)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return settings, nil
}

// ICoreWebView2Environment

type iCoreWebView2EnvironmentVtbl struct {
	_IUnknownVtbl
	CreateCoreWebView2Controller     jchv.ComProc
	CreateWebResourceResponse        jchv.ComProc
	GetBrowserVersionString          jchv.ComProc
	AddNewBrowserVersionAvailable    jchv.ComProc
	RemoveNewBrowserVersionAvailable jchv.ComProc
}

type ICoreWebView2Environment struct {
	vtbl *iCoreWebView2EnvironmentVtbl
}

func (e *ICoreWebView2Environment) CreateWebResourceResponse(content []byte, statusCode int, reasonPhrase string, headers string) (*jchv.ICoreWebView2WebResourceResponse, error) {
	var err error
	var stream uintptr

	if len(content) > 0 {
		// Create stream for response
		stream, err = w32.SHCreateMemStream(content)
		if err != nil {
			return nil, err
		}
	}

	// Convert string 'uri' to *uint16
	_reason, err := windows.UTF16PtrFromString(reasonPhrase)
	if err != nil {
		return nil, err
	}
	// Convert string 'uri' to *uint16
	_headers, err := windows.UTF16PtrFromString(headers)
	if err != nil {
		return nil, err
	}
	var response *jchv.ICoreWebView2WebResourceResponse
	_, _, err = e.vtbl.CreateWebResourceResponse.Call(
		uintptr(unsafe.Pointer(e)),
		stream,
		uintptr(statusCode),
		uintptr(unsafe.Pointer(_reason)),
		uintptr(unsafe.Pointer(_headers)),
		uintptr(unsafe.Pointer(&response)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return response, nil

}

// ICoreWebView2WebMessageReceivedEventArgs

type iCoreWebView2WebMessageReceivedEventArgsVtbl struct {
	_IUnknownVtbl
	GetSource                jchv.ComProc
	GetWebMessageAsJSON      jchv.ComProc
	TryGetWebMessageAsString jchv.ComProc
}

type iCoreWebView2WebMessageReceivedEventArgs struct {
	vtbl *iCoreWebView2WebMessageReceivedEventArgsVtbl
}

// ICoreWebView2PermissionRequestedEventArgs

type iCoreWebView2PermissionRequestedEventArgsVtbl struct {
	_IUnknownVtbl
	GetURI             jchv.ComProc
	GetPermissionKind  jchv.ComProc
	GetIsUserInitiated jchv.ComProc
	GetState           jchv.ComProc
	PutState           jchv.ComProc
	GetDeferral        jchv.ComProc
}

type iCoreWebView2PermissionRequestedEventArgs struct {
	vtbl *iCoreWebView2PermissionRequestedEventArgsVtbl
}

// ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler

type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl interface {
	_IUnknownImpl
	EnvironmentCompleted(res uintptr, env *ICoreWebView2Environment) uintptr
}

type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl struct {
	_IUnknownVtbl
	Invoke jchv.ComProc
}

type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler struct {
	vtbl *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl
	impl iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl
}

func _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownQueryInterface(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownAddRef(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownRelease(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, res uintptr, env *ICoreWebView2Environment) uintptr {
	return this.impl.EnvironmentCompleted(res, env)
}

var iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerFn = iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl{
	_IUnknownVtbl{
		jchv.NewComProc(_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownQueryInterface),
		jchv.NewComProc(_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownAddRef),
		jchv.NewComProc(_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownRelease),
	},
	jchv.NewComProc(_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke),
}

func newICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler(impl iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl) *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler {
	return &iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler{
		vtbl: &iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerFn,
		impl: impl,
	}
}

// ICoreWebView2WebMessageReceivedEventHandler

type iCoreWebView2WebMessageReceivedEventHandlerImpl interface {
	_IUnknownImpl
	MessageReceived(sender *ICoreWebView2, args *iCoreWebView2WebMessageReceivedEventArgs) uintptr
}

type iCoreWebView2WebMessageReceivedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke jchv.ComProc
}

type iCoreWebView2WebMessageReceivedEventHandler struct {
	vtbl *iCoreWebView2WebMessageReceivedEventHandlerVtbl
	impl iCoreWebView2WebMessageReceivedEventHandlerImpl
}

func _ICoreWebView2WebMessageReceivedEventHandlerIUnknownQueryInterface(this *iCoreWebView2WebMessageReceivedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2WebMessageReceivedEventHandlerIUnknownAddRef(this *iCoreWebView2WebMessageReceivedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2WebMessageReceivedEventHandlerIUnknownRelease(this *iCoreWebView2WebMessageReceivedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2WebMessageReceivedEventHandlerInvoke(this *iCoreWebView2WebMessageReceivedEventHandler, sender *ICoreWebView2, args *iCoreWebView2WebMessageReceivedEventArgs) uintptr {
	return this.impl.MessageReceived(sender, args)
}

var iCoreWebView2WebMessageReceivedEventHandlerFn = iCoreWebView2WebMessageReceivedEventHandlerVtbl{
	_IUnknownVtbl{
		jchv.NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerIUnknownQueryInterface),
		jchv.NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerIUnknownAddRef),
		jchv.NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerIUnknownRelease),
	},
	jchv.NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerInvoke),
}

func newICoreWebView2WebMessageReceivedEventHandler(impl iCoreWebView2WebMessageReceivedEventHandlerImpl) *iCoreWebView2WebMessageReceivedEventHandler {
	return &iCoreWebView2WebMessageReceivedEventHandler{
		vtbl: &iCoreWebView2WebMessageReceivedEventHandlerFn,
		impl: impl,
	}
}

// ICoreWebView2PermissionRequestedEventHandler

type iCoreWebView2PermissionRequestedEventHandlerImpl interface {
	_IUnknownImpl
	PermissionRequested(sender *ICoreWebView2, args *iCoreWebView2PermissionRequestedEventArgs) uintptr
}

type iCoreWebView2PermissionRequestedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke jchv.ComProc
}

type iCoreWebView2PermissionRequestedEventHandler struct {
	vtbl *iCoreWebView2PermissionRequestedEventHandlerVtbl
	impl iCoreWebView2PermissionRequestedEventHandlerImpl
}

func _ICoreWebView2PermissionRequestedEventHandlerIUnknownQueryInterface(this *iCoreWebView2PermissionRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2PermissionRequestedEventHandlerIUnknownAddRef(this *iCoreWebView2PermissionRequestedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2PermissionRequestedEventHandlerIUnknownRelease(this *iCoreWebView2PermissionRequestedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2PermissionRequestedEventHandlerInvoke(this *iCoreWebView2PermissionRequestedEventHandler, sender *ICoreWebView2, args *iCoreWebView2PermissionRequestedEventArgs) uintptr {
	return this.impl.PermissionRequested(sender, args)
}

var iCoreWebView2PermissionRequestedEventHandlerFn = iCoreWebView2PermissionRequestedEventHandlerVtbl{
	_IUnknownVtbl{
		jchv.NewComProc(_ICoreWebView2PermissionRequestedEventHandlerIUnknownQueryInterface),
		jchv.NewComProc(_ICoreWebView2PermissionRequestedEventHandlerIUnknownAddRef),
		jchv.NewComProc(_ICoreWebView2PermissionRequestedEventHandlerIUnknownRelease),
	},
	jchv.NewComProc(_ICoreWebView2PermissionRequestedEventHandlerInvoke),
}

func newICoreWebView2PermissionRequestedEventHandler(impl iCoreWebView2PermissionRequestedEventHandlerImpl) *iCoreWebView2PermissionRequestedEventHandler {
	return &iCoreWebView2PermissionRequestedEventHandler{
		vtbl: &iCoreWebView2PermissionRequestedEventHandlerFn,
		impl: impl,
	}
}

func (i *ICoreWebView2) AddWebResourceRequestedFilter(uri string, resourceContext jchv.COREWEBVIEW2_WEB_RESOURCE_CONTEXT) error {
	var err error
	// Convert string 'uri' to *uint16
	_uri, err := windows.UTF16PtrFromString(uri)
	if err != nil {
		return err
	}
	_, _, err = i.vtbl.AddWebResourceRequestedFilter.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
		uintptr(resourceContext),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
func (i *ICoreWebView2) AddNavigationCompleted(eventHandler *jchv.ICoreWebView2NavigationCompletedEventHandler, token *_EventRegistrationToken) error {
	var err error
	_, _, err = i.vtbl.AddNavigationCompleted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
