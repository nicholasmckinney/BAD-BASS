package resource

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"unsafe"
)

const BinaryResourceType = "BINARY"

func Embed(targetFilePath, resourceFilePath, name string) error {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	BeginUpdateResourceW := kernel32.NewProc("BeginUpdateResourceW")
	EndUpdateResourceW := kernel32.NewProc("EndUpdateResourceW")
	UpdateResourceW := kernel32.NewProc("UpdateResourceW")

	resourceData, err := os.ReadFile(resourceFilePath)
	if err != nil {
		return fmt.Errorf("error while reading resource file (%s): %w", resourceFilePath, err)
	}

	targetFilePathW := windows.StringToUTF16(targetFilePath)
	hResource, _, _ := BeginUpdateResourceW.Call(
		uintptr(unsafe.Pointer(&targetFilePathW[0])),
		uintptr(0), // FALSE: Don't delete existing resources
	)

	if int(hResource) == 0 {
		return fmt.Errorf("error occurred while opening resource handle to target file: %s", targetFilePath)
	}

	resourceType := windows.StringToUTF16(BinaryResourceType)
	resourceName := windows.StringToUTF16(name)

	ret, _, _ := UpdateResourceW.Call(
		hResource,
		uintptr(unsafe.Pointer(&resourceType[0])),
		uintptr(unsafe.Pointer(&resourceName[0])),
		windows.LANG_ENGLISH,
		uintptr(unsafe.Pointer(&resourceData[0])),
		uintptr(len(resourceData)),
	)

	if int(ret) == 0 {
		return fmt.Errorf("error occurred while updating resource in target file: %s", targetFilePath)
	}

	ret, _, _ = EndUpdateResourceW.Call(
		hResource,
		uintptr(0), // discard changes = FALSE (i.e. write modifications to file)
	)

	if int(ret) == 0 {
		return fmt.Errorf("error occurred while flushing modifications to file: %s", targetFilePath)
	}

	return nil
}
