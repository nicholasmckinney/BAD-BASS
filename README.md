# BAD BASS

BAD BASS is a proof-of-concept for a browser-in-the-browser phishing technique involving WebView injection. Essentially, a new window is injected over the browser's child window that handles the rendered web-content. This webview can then be loaded with phishing pages to collect user information.

[Read More About It](https://malware.tech/posts/bad-bass/)

# Requirements

* Visual Studio 2022
* Go >= 1.18.3
* [Donut](https://github.com/TheWover/donut)

# Building

1) Compile WEBPHISH to DLL using the build.bat script
2) Compile WEBPHISH DLL to shellcode using Donut
3) Place shellcode binary file into LIVEBAIT/payload/loader.bin
4) Compile LIVEBAIT with Visual Studio
5) Package a web-inject archive with PHISHROD (see -h)
6) Embed the web-inject archive from (5) using PHISHROD into the LIVEBAIT executable (see PHISHROD -h for assistance)

# Credits

Huge credits to jchv for the [go-webview2](https://github.com/jchv/go-webview2) project that wraps the WebView interfaces. Virtually all of the browser code was pulled from that project. Small modifications were made that required some of the files be pulled into the WEBPHISH sources.