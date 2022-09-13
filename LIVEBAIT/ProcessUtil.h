#pragma once
#include <Windows.h>
#include <unordered_set>
#include "xor.h"

namespace Process {
	typedef DWORD PID;

	DWORD FindExplorerPID(HANDLE hSnapshot);
	std::unordered_set<PID> EnumerateTopLevelBrowsers();
}