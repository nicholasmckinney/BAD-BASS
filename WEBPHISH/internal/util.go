package internal

import (
	"Webphish/internal/w32"
	"crypto/rand"
	"encoding/binary"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ErrorCode int

const (
	GenericSuccess = 0
	GenericError   = -1
)

const (
	ERR_CONSTRUCTOR               = -10
	ERR_WRITE_FAIL                = -20
	ERR_CREATE_NAMED_PIPE         = -101
	ERR_PIPE_CLIENT_CONNECT       = -102
	ERR_CREATE_MAILSLOT           = -150
	ERR_READ_MAILSLOT_MSG         = -151
	ERR_CONNECT_MAILSLOT          = -152
	ERR_CREATE_EVENT              = -201
	ERR_FAILED_NOTIFY             = -301
	ERR_DECODER_CIPHER_CREATION   = -400
	ERR_DECODER_UNPACK            = -401
	ERR_DECODER_READ_CONTENT      = -402
	ERR_LOADER_LOAD_CONFIGURATION = -500
	ERR_LOADER_PUT_FILE           = -501
)

const (
	Alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

func GetUserName() string {
	var maxSize uint32
	maxSize = 1024
	username := make([]uint16, maxSize)

	ret, _, _ := w32.Advapi32GetUserName.Call(uintptr(unsafe.Pointer(&username[0])), uintptr(unsafe.Pointer(&maxSize)))
	if int(ret) > 0 {
		return syscall.UTF16ToString(username)
	}
	return ""
}

func BytesToWideString(input []byte) string {
	DataReceived := make([]uint16, (len(input)/2)+1)
	var pos int
	for idx := 0; idx < len(input)-1; idx += 2 {
		DataReceived[pos] = binary.LittleEndian.Uint16(input[idx : idx+2])
		pos++
	}
	return windows.UTF16ToString(DataReceived)
}

func RandomString(length int) string {
	ll := len(Alphanumeric)
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = Alphanumeric[int(b[i])%ll]
	}
	return string(b)
}
