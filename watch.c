#ifdef _WIN32
#include "watch.h"
#include <windows.h>
#include <tchar.h>
#include <stdio.h>

const DWORD dwFilter = REG_NOTIFY_CHANGE_NAME |
                       REG_NOTIFY_CHANGE_ATTRIBUTES |
                       REG_NOTIFY_CHANGE_LAST_SET |
                       REG_NOTIFY_CHANGE_SECURITY;


#define HANDLE_ERR(val) if (val != ERROR_SUCCESS) return val

long watcher_create(void* hive, char* regPath, Watcher* out) {
	HKEY hMainKey = (HKEY) hive;
	HKEY hKey;

	DEBUG_PRINT("RegOpenKeyEx '%s'\n", regPath);
	LONG lErrorCode = RegOpenKeyEx(hMainKey, regPath, 0, KEY_NOTIFY, &hKey);
	HANDLE_ERR(lErrorCode);

	DEBUG_PRINT("CreateEvent\n");
	HANDLE hEvent = CreateEvent(NULL, TRUE, FALSE, NULL);
	if (hEvent == NULL) 
		return GetLastError();

	out->hKey = hKey;
	out->hEvent = hEvent;
	return 0;
}

long watcher_await(Watcher* out, long timeout, boolean* changed) {
	LONG lErrorCode;
	DEBUG_PRINT("RegNotifyChangeKeyValue\n");
	lErrorCode = RegNotifyChangeKeyValue(out->hKey, TRUE, dwFilter, out->hEvent, TRUE);
	HANDLE_ERR(lErrorCode);
	DEBUG_PRINT("WaitForSingleObject\n");
	DWORD dwEvent = WaitForSingleObject(out->hEvent, timeout);
	ResetEvent(out->hEvent);
	switch (dwEvent) {
		case WAIT_TIMEOUT:
			*changed = FALSE;
			break;
		case WAIT_OBJECT_0:
			*changed = TRUE;
			break;
		default:
			return dwEvent; 
	}

	return 0;
}

long watcher_close(Watcher* out) {
	DEBUG_PRINT("RegCloseKey\n");
	LONG lErrorCode = RegCloseKey(out->hKey);
	HANDLE_ERR(lErrorCode);
	DEBUG_PRINT("CloseHandle\n");
	CloseHandle(out->hEvent);
}

#endif