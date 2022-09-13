package console

import (
	"syscall"
	"unsafe"
)

const (
	ATTACH_PARENT_PROCESS = ^uint32(0) // (DWORD)-1
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procAttachConsole = modkernel32.NewProc("AttachConsole")
)

func AttachConsole(dwParentProcess uint32) (ok bool) {
	r0, _, _ := syscall.Syscall(procAttachConsole.Addr(), 1, uintptr(dwParentProcess), 0, 0)
	ok = bool(r0 != 0)
	return
}

// MessageBox of Win32 API.
func MessageBox(hwnd uintptr, caption, title string, flags uint) int {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(flags))

	return int(ret)
}

// MessageBoxPlain of Win32 API.
func MessageBoxPlain(title, caption string) int {
	const (
		NULL  = 0
		MB_OK = 0
	)
	return MessageBox(NULL, caption, title, MB_OK)
}
