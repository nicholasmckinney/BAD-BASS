//go:build windows
// +build windows

package edge

import (
	"unsafe"

	"Webphish/internal/w32"
)

func (e *Chromium) Resize() {
	if e.controller == nil {
		return
	}
	var bounds w32.Rect
	_, _, _ = w32.User32GetClientRect.Call(e.hwnd, uintptr(unsafe.Pointer(&bounds)))
	_, _, _ = e.controller.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(e.controller)),
		uintptr(unsafe.Pointer(&bounds)),
	)
}
