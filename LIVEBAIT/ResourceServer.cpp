#include "ResourceServer.h"
#include <iostream>
#include <string>
#include "common.h"
#include "Util.h"
#include "xor.h"

namespace ResourceServer {
	const std::string& MAILSLOT_PREFIX =Xor( "\\\\.\\mailslot\\");

	DWORD WINAPI ServeWebResources(LPVOID parameter) {
		HANDLE hMyMailslot;
		HANDLE hDestinationMailslot;
		DWORD dNextMessageSize;
		DWORD dMessageCount;
		DWORD dNumBytesRead;
		DWORD dBytesWritten;
		LPVOID lpMessage;

		std::wstring rsrc_name = L"500";
		Common::EmbeddedResource webResources = Util::Load(rsrc_name);
		BYTE* payload = (BYTE*)webResources.Payload;
		
		std::wstring username = Util::GetUsername();
		
		std::wstring mailSlotName = Util::s2ws(MAILSLOT_PREFIX.c_str()) + L"mx" + username;
		
		hMyMailslot = CreateMailslotW(
			mailSlotName.c_str(),
			0, // ALLOWS ANY SIZE MESSAGE
			MAILSLOT_WAIT_FOREVER,
			NULL
		);

		if (hMyMailslot == INVALID_HANDLE_VALUE) {
			return -1;
		}

		while (true) {
			Sleep(5000);
			if (!GetMailslotInfo(hMyMailslot, NULL, &dNextMessageSize, &dMessageCount, NULL)) {
				break;
			}

			if (dNextMessageSize == MAILSLOT_NO_MESSAGE) {
				continue;
			}

			if (!dMessageCount) {
				continue;
			}

			std::wcout << "Number of messages in mailslot: " << dMessageCount << std::endl;

			lpMessage = VirtualAlloc(NULL, dNextMessageSize, MEM_COMMIT, PAGE_READWRITE);

			if (!lpMessage) {
				std::wcerr << L"Failed to allocate memory for mailslot message" << std::endl;
				break;
			}

			if (!ReadFile(hMyMailslot, lpMessage, dNextMessageSize, &dNumBytesRead, NULL)) {
				std::wcerr << "Failed to read message from mailslot" << std::endl;
				VirtualFree(lpMessage, 0, MEM_DECOMMIT);
				break;
			}

			if (dNextMessageSize != dNumBytesRead) {
				std::wcerr << "Failed to read entire message from mailslot" << std::endl;
				VirtualFree(lpMessage, 0, MEM_DECOMMIT);
				break;
			}
			
			std::wstring destinationMailslot(reinterpret_cast<wchar_t*>(lpMessage), dNextMessageSize / (sizeof(wchar_t)));

			hDestinationMailslot = CreateFileW(
				destinationMailslot.c_str(),
				GENERIC_WRITE, 
				FILE_SHARE_READ, 
				NULL, 
				OPEN_EXISTING, 
				FILE_ATTRIBUTE_NORMAL, 
				NULL
			);

			if (hDestinationMailslot == INVALID_HANDLE_VALUE) {
				VirtualFree(lpMessage, 0, MEM_DECOMMIT);
				continue;
			}
			WriteFile(hDestinationMailslot, (LPCVOID)payload, webResources.PayloadSize, &dBytesWritten, NULL);
			std::wcout << L"Wrote " << dBytesWritten << L" to destination mailslot " << destinationMailslot << std::endl;
			
			CloseHandle(hDestinationMailslot);
			VirtualFree(lpMessage, 0, MEM_DECOMMIT);
		}

		VirtualFree(webResources.Payload, 0, MEM_DECOMMIT);
		CloseHandle(hMyMailslot);
		return 0;
	}

	HANDLE Start() {
		return CreateThread(NULL, 0, ServeWebResources, NULL, 0, NULL);
	}
}
