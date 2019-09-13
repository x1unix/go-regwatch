#ifndef WATCH_H
#define WATCH_H

#include <windows.h>
#include <stdlib.h>

#if defined(DEBUG) && DEBUG > 0
 #define DEBUG_PRINT(fmt, ...) fprintf(stderr, "DEBUG: %s:%d:%s(): " fmt, \
    __FILE__, __LINE__, __func__, ##__VA_ARGS__);

#else
 #define DEBUG_PRINT(fmt, ...) /* Don't do anything in release builds */
#endif

typedef struct Watcher {
	HKEY hKey;
	HANDLE hEvent;
} Watcher;

long watcher_create(void* hive, char* regPath, Watcher* out);
long watcher_await(Watcher* out, long timeout, boolean* changed);
long watcher_close(Watcher* out);
#endif