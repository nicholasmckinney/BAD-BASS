package pipe

import (
	"Webphish/internal"
	"Webphish/internal/w32"
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

const (
	PIPE_ACCESS_DUPLEX = 0x00000003
	PIPE_TYPE_MESSAGE  = 0x00000004
	PIPE_NO_WAIT       = 0x00000001
	MAX_INSTANCES      = 255
	BUFFER_SIZE        = 5000
	DEFAULT_WAIT       = 0 // ends up as 50 ms (https://docs.microsoft.com/en-us/windows/win32/api/namedpipeapi/nf-namedpipeapi-createnamedpipew)
)

type NamedPipeWriter struct {
	path       string
	PipeHandle uintptr
}

func New(name string) *NamedPipeWriter {
	result := NamedPipeWriter{
		path: fmt.Sprintf("\\\\.\\pipe\\%s", name),
	}
	return &result
}

func (p *NamedPipeWriter) Capture(msg string, ready chan bool) internal.ErrorCode {
	pipeName, _ := windows.UTF16FromString(p.path)

	// https://docs.microsoft.com/en-us/windows/win32/api/namedpipeapi/nf-namedpipeapi-createnamedpipew
	hPipe, _, _ := w32.Kernel32CreateNamedPipe.Call(
		uintptr(unsafe.Pointer(&pipeName[0])),
		uintptr(PIPE_ACCESS_DUPLEX),
		uintptr(PIPE_TYPE_MESSAGE),
		uintptr(MAX_INSTANCES),
		uintptr(BUFFER_SIZE),
		uintptr(BUFFER_SIZE),
		uintptr(DEFAULT_WAIT),
		uintptr(0), // security attributes NULL
	)

	if int(hPipe) <= 0 { // failure
		ready <- false
		return internal.ERR_CREATE_NAMED_PIPE
	}

	p.PipeHandle = hPipe
	ready <- true

	success, _, _ := w32.Kernel32ConnectNamedPipe.Call(
		p.PipeHandle,
		uintptr(0),
	)

	if int(success) <= 0 { // failure
		w32.Kernel32CloseHandle.Call(p.PipeHandle)
		return internal.ERR_PIPE_CLIENT_CONNECT
	}

	output, _ := windows.UTF16FromString(msg)
	outputLen := len(output) * 2
	var bytesWritten w32.DWORD
	success, _, _ = w32.Kernel32WriteFile.Call(
		p.PipeHandle,
		uintptr(unsafe.Pointer(&output[0])),
		uintptr(w32.DWORD(outputLen)),
		uintptr(unsafe.Pointer(&bytesWritten)),
		uintptr(0),
	)

	if success <= 0 {
		w32.Kernel32CloseHandle.Call(p.PipeHandle)
		return internal.ERR_WRITE_FAIL
	}

	w32.Kernel32CloseHandle.Call(p.PipeHandle)
	return 0
}
