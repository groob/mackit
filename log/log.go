package log

/*
#include <os/log.h>

void oslog(char* message) {
	os_log(OS_LOG_DEFAULT, "%{public}s",message);
}

*/
import "C"

func Log(message string) {
	C.oslog(C.CString(message))
}
