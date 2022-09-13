#pragma once
#include <Windows.h>

namespace ResourceServer {
	DWORD WINAPI ServeWebResources(LPVOID parameter);
	HANDLE Start();
}

