#ifndef WATCH_H
#define WATCH_H

#define DEBUG 3
#include <windows.h>
#include <stdlib.h>

#if defined(DEBUG) && DEBUG > 0
 #define DEBUG_PRINT(fmt, ...) fprintf(stderr, "DEBUG: %s:%d:%s(): " fmt, \
    __FILE__, __LINE__, __func__, ##__VA_ARGS__);

#else
 #define DEBUG_PRINT(fmt, ...) /* Don't do anything in release builds */
#endif


typedef boolean (*callback)(char*, long);

long reg_add_listener(callback cb, HKEY hMainKey, char* regPath);
long reg_listen(HKEY hMainKey, char* regPath);

#endif