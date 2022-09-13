#include "ProcessUtil.h"
#include <Psapi.h>
#include <iostream>
#include <tlhelp32.h>
#include "Util.h"
#include <unordered_set>
#include "xor.h"

namespace Process {

    const std::string& CHROME_BROWSER = Xor("chrome.exe");
    const std::string& FIREFOX_BROWSER = Xor("firefox.exe");
    const std::string& EDGE_BROWSER = Xor("edge.exe");
    const std::string& EXPLORER = Xor("explorer.exe");

    std::unordered_set<std::string> BROWSERS = { CHROME_BROWSER, FIREFOX_BROWSER, EDGE_BROWSER };

    DWORD FindExplorerPID(HANDLE hSnapshot) {
        PROCESSENTRY32W pe = { 0 };
        pe.dwSize = sizeof(PROCESSENTRY32W);

        if (!SUCCEEDED(Process32FirstW(hSnapshot, &pe))) {
            return -1;
        }

        do {
            std::wstring exe(pe.szExeFile);
            std::string comparable = Util::ws2s(exe);
            //std::cout << "Comparing exe: " << comparable << " to " << EXPLORER.c_str() << std::endl;

            if (comparable.compare(EXPLORER.c_str()) == 0) {
                return pe.th32ProcessID;
            }
        } while (Process32NextW(hSnapshot, &pe));


        return -1;
    }

    std::unordered_set<PID> EnumerateTopLevelBrowsers() {
        HANDLE hProcListSnapshot;
        auto result = std::unordered_set<PID>();

        hProcListSnapshot = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
        if (hProcListSnapshot == INVALID_HANDLE_VALUE) {
            std::cerr << "Failed to get snapshot" << std::endl;
            return result;
        }

        DWORD explorerPid = FindExplorerPID(hProcListSnapshot);
        if (explorerPid == -1) {
            return result;
        }

        PROCESSENTRY32W pe = { 0 };
        pe.dwSize = sizeof(PROCESSENTRY32W);

        if (!SUCCEEDED(Process32FirstW(hProcListSnapshot, &pe))) {
            CloseHandle(hProcListSnapshot);
            return result;
        }

        do {
            std::wstring exe(pe.szExeFile);
            auto comparable = Util::ws2s(exe);

            for (auto browser : BROWSERS) {
                if (comparable.compare(browser.c_str()) == 0) {
                    // process name matches one of edge.exe, firefox.exe, chrome.exe
                    if (pe.th32ParentProcessID == explorerPid) {
                        result.insert(pe.th32ProcessID);
                    }
                }
            }
        } while (Process32Next(hProcListSnapshot, &pe));
        return result;
    }
}