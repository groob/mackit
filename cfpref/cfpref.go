package cfpref

/*
#cgo LDFLAGS: -framework IOKit -framework ApplicationServices
#include <CoreFoundation/CoreFoundation.h>
#include <stdlib.h>

// cfstring_utf8_length returns the number of characters successfully converted to UTF-8 and
// the bytes required to store them.
static inline CFIndex cfstring_utf8_length(CFStringRef str, CFIndex *need) {
  CFIndex n, usedBufLen;
  CFRange rng = CFRangeMake(0, CFStringGetLength(str));
  return CFStringGetBytes(str, rng, kCFStringEncodingUTF8, 0, 0, NULL, 0, need);
}
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type CFTypeID uint

type CFPropertyListRef struct {
	ref C.CFPropertyListRef
}

func SetAppValue(key, value, appID string) {
	cKey := cfstring(key)
	cVal := cfstring(value)
	cAppID := cfstring(appID)

	C.CFPreferencesSetAppValue(cKey, cVal, cAppID)
}

func CopyAppValue(key, appID string) CFPropertyListRef {
	cKey := cfstring(key)
	cAppID := cfstring(appID)

	appValue := C.CFPreferencesCopyAppValue(cKey, cAppID)
	return CFPropertyListRef{ref: appValue}
}

func cfstringGo(cfs C.CFStringRef) string {
	var usedBufLen C.CFIndex
	n := C.cfstring_utf8_length(cfs, &usedBufLen)
	if n <= 0 {
		return ""
	}
	rng := C.CFRange{location: C.CFIndex(0), length: n}
	buf := make([]byte, int(usedBufLen))

	bufp := unsafe.Pointer(&buf[0])
	C.CFStringGetBytes(cfs, rng, C.kCFStringEncodingUTF8, 0, 0, (*C.UInt8)(bufp), C.CFIndex(len(buf)), &usedBufLen)

	sh := &reflect.StringHeader{
		Data: uintptr(bufp),
		Len:  int(usedBufLen),
	}
	return *(*string)(unsafe.Pointer(sh))
}

// cfstring efficiently creates a CFString from a Go String.
func cfstring(s string) C.CFStringRef {
	n := C.CFIndex(len(s))
	return C.CFStringCreateWithBytes(nil, *(**C.UInt8)(unsafe.Pointer(&s)), n, C.kCFStringEncodingUTF8, 0)
}

func (plisRef CFPropertyListRef) CFTypeID() CFTypeID {
	typeId := C.CFGetTypeID(C.CFTypeRef(plisRef.ref))
	return CFTypeID(typeId)
}
