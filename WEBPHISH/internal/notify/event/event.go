package event

import (
	"Webphish/internal"
	"Webphish/internal/w32"
	"fmt"
	"syscall"
	"unsafe"
)

const (
	FALSE = 0
)

type WinEventNotifier struct {
	event uintptr
}

func New(name string) (*WinEventNotifier, internal.ErrorCode) {
	name = fmt.Sprintf("ev%s", name)
	eventName, _ := syscall.UTF16FromString(name)
	// this event should have already been created by the loader
	hEvent, _, _ := w32.Kernel32CreateEvent.Call(
		uintptr(0),
		uintptr(0),
		uintptr(FALSE), // auto-reset state as non-signaled after a waiting thread is notified (i.e. C++ loader thread)
		uintptr(unsafe.Pointer(&eventName[0])),
	)

	if int(hEvent) == 0 { // failure
		return nil, internal.ERR_CREATE_EVENT
	}
	return &WinEventNotifier{event: hEvent}, internal.GenericSuccess
}

func (w *WinEventNotifier) Notify() bool {
	success, _, _ := w32.Kernel32SetEvent.Call(w.event)
	return int(success) != 0
}
