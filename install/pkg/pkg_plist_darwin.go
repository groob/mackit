package pkg

/*
#cgo darwin CFLAGS: -DDARWIN -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
const char * RemoveIFPkgPathMappingsKey(char *path)
{
	NSData *plistData = [NSData dataWithContentsOfFile:[NSString stringWithFormat:@"%s", path]];
	NSError *error;
	NSPropertyListFormat format;
	NSMutableDictionary *plist = [NSPropertyListSerialization propertyListWithData:plistData options:NSPropertyListMutableContainers format:&format error:&error];
	[plist removeObjectForKey:@"IFPkgPathMappings"];

	NSData *data = [NSPropertyListSerialization dataWithPropertyList: plist format:format options: 0 error: &error];
	if (data == nil) {
		const char * cerr = [[NSString stringWithFormat:@"ERROR: %@", error.localizedDescription] UTF8String];
		return cerr;
	}

	const char * cString = [data bytes];
	return cString;
}
*/
import "C"
import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

// if the pkg is a bundle, remove the IFPkgPathMappings key from the plist.
func rmIFPkgPathMappingsFromPlist(pkgpath string) error {
	path := filepath.Join(pkgpath, "Contents/Info.plist")
	fi, err := os.Stat(path)
	if err != nil {
		return nil
	}
	data := C.RemoveIFPkgPathMappingsKey(C.CString(path))
	defer C.free(unsafe.Pointer(data))
	strData := strings.TrimSpace(C.GoString(data))
	if strings.Contains(strData, "ERROR:") {
		return errors.New(strData)
	}

	return ioutil.WriteFile(path, []byte(strData), fi.Mode())
}
