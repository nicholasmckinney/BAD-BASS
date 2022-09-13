#pragma once
#include <string>
#include "common.h"

namespace Util {
	std::wstring s2ws(const std::string& str);
	std::string ws2s(const std::wstring& wstr);
	std::wstring GetUsername();
	Common::EmbeddedResource Load(std::wstring& name);
	Common::EmbeddedResource Load(DWORD id);
}
