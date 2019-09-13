#ifdef _WIN32
#include "watch.h"
#include <windows.h>
#include <tchar.h>
#include <stdio.h>

boolean reg_global_listener(char* regPath, long lEvent);

const DWORD dwFilter = REG_NOTIFY_CHANGE_NAME |
                       REG_NOTIFY_CHANGE_ATTRIBUTES |
                       REG_NOTIFY_CHANGE_LAST_SET |
                       REG_NOTIFY_CHANGE_SECURITY;

typedef struct Watcher {
	HKEY hKey;
	HANDLE hEvent;
} Watcher;

long reg_add_listener(callback cb, HKEY hMainKey, char* regPath) {
	HKEY hKey;

	DEBUG_PRINT("RegOpenKeyEx '%s'\n", regPath);
	LONG lErrorCode = RegOpenKeyEx(hMainKey, regPath, 0, KEY_NOTIFY, &hKey);
	if (lErrorCode != ERROR_SUCCESS) 
		return lErrorCode;

	DEBUG_PRINT("CreateEvent\n");
	HANDLE hEvent = CreateEvent(NULL, TRUE, FALSE, NULL);
	if (hEvent == NULL) 
		return GetLastError();

	while(TRUE) {
		DEBUG_PRINT("RegNotifyChangeKeyValue\n");
		lErrorCode = RegNotifyChangeKeyValue(hKey, TRUE, dwFilter, hEvent, TRUE);
		if(lErrorCode != ERROR_SUCCESS) break;

		DEBUG_PRINT("WaitForSingleObject\n");
		DWORD dwEvent = WaitForSingleObject(hEvent, INFINITE);
		ResetEvent(hEvent);
		if (dwEvent == WAIT_OBJECT_0) {
			boolean result = cb(regPath, dwEvent);
			DEBUG_PRINT("Callback returned %d\n", result);
			if (!result) {
				DEBUG_PRINT("break\n");
				break;
			}
			continue;
		}

		lErrorCode = dwEvent;
		break;
	}
	
	DEBUG_PRINT("Close watcher\n");
	RegCloseKey(hKey);
	CloseHandle(hEvent);
	free(regPath);	// Allocated in heap by C.GoString()
	return lErrorCode;
}

long reg_listen(HKEY hMainKey, char* regPath) {
  return reg_add_listener(reg_global_listener, hMainKey, regPath);
}
#endif