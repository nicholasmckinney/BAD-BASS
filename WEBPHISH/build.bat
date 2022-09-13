ECHO "Cleaning build directory..."
RD /S /Q build
ECHO "Making build directory"
MD build
ECHO "Downloading packages..."
go get github.com/jchv/go-webview2
github.com/hidez8891/shm
ECHO "Building DLL"
go build -buildmode=c-shared -ldflags="-s -w" -o build/webphish.dll .