#pragma once
#include <Windows.h>

namespace Common {
	typedef struct {
		LPVOID  Payload;
		DWORD   PayloadSize;
	} EmbeddedResource;
}