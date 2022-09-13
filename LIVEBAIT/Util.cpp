#include "Util.h"
#include <locale>
#include <codecvt>
#include <Windows.h>


namespace Util {

    std::wstring s2ws(const std::string& str)
    {
        using convert_typeX = std::codecvt_utf8<wchar_t>;
        std::wstring_convert<convert_typeX, wchar_t> converterX;

        return converterX.from_bytes(str);
    }

    std::string ws2s(const std::wstring& wstr)
    {
        using convert_typeX = std::codecvt_utf8<wchar_t>;
        std::wstring_convert<convert_typeX, wchar_t> converterX;

        return converterX.to_bytes(wstr);
    }

    std::wstring GetUsername() {
        DWORD bufSize = 100;
        WCHAR username[100];
        HANDLE hMutex;
        BOOL success;

        success = GetUserNameW(username, &bufSize);
        if (!success) {
            return L"";
        }

        return std::wstring(username);
    }

    void LoadResourceData(Common::EmbeddedResource* lpRsrc, HMODULE hProcess, HRSRC hRsrc) {
        LPVOID lpData;
        DWORD dRsrcSize;

        lpData = LoadResource(hProcess, hRsrc);
        if (!lpData) {
            return;
        }

        dRsrcSize = SizeofResource(hProcess, hRsrc);
        if (!dRsrcSize) {
            return;
        }

        lpRsrc->Payload = VirtualAlloc(0, dRsrcSize, MEM_COMMIT, PAGE_READWRITE);
        lpRsrc->PayloadSize = dRsrcSize;
        RtlMoveMemory(lpRsrc->Payload, lpData, lpRsrc->PayloadSize);
    }

    Common::EmbeddedResource Load(std::wstring& name) {
        HMODULE hSelfProcess;
        HRSRC hResource;
        DWORD dResourceSize;
        LPVOID lpData;

        Common::EmbeddedResource payload{ NULL, -1 };
        hSelfProcess = GetModuleHandle(NULL);
        hResource = FindResourceW(hSelfProcess, name.c_str(), L"BINARY");
        if (!hResource) {
            return payload;
        }

        LoadResourceData(&payload, hSelfProcess, hResource);
        FreeResource(hResource);
        return payload;
    }

    Common::EmbeddedResource Load(DWORD id)
    {
        HMODULE hSelfProcess;
        HRSRC hResource;
        DWORD dResourceSize;
        LPVOID lpShellcode;
        Common::EmbeddedResource payload{ NULL, -1 };


        hSelfProcess = GetModuleHandle(NULL);
        hResource = FindResourceW(hSelfProcess, MAKEINTRESOURCE(id), L"BINARY");
        if (!hResource) {
            return payload;
        }

        LoadResourceData(&payload, hSelfProcess, hResource);
        FreeResource(hResource);
        return payload;
    }
}
