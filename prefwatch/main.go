package main

/*
#cgo darwin CFLAGS: -DDARWIN -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

// import exported cgo callback from the other go file.
extern void GoCallback(char * prefPlist);


// implement a handler for userdefaults
// references:
// http://nshipster.com/key-value-observing/
// https://github.com/shurcooL/trayhost/blob/master/platform/darwin/tray.m
// https://github.com/google/santa/blob/363826502fe592e0b5695dbb422d7b8de36f0dc6/Source/SantaGUI/SNTMessageWindowController.m#L44-L77
@interface ManageHandler : NSObject<NSUserNotificationCenterDelegate>
@property (retain) NSUserDefaults *userDefaults;
@end

ManageHandler * uncDelegate;

@implementation ManageHandler
-(instancetype) init {
	self = [super init];
	if (self) {
		_userDefaults = [[NSUserDefaults alloc] initWithSuiteName:@"ManagedInstalls"];
		[_userDefaults addObserver:self
			forKeyPath:@"LastCheckDate"
			options:NSKeyValueObservingOptionNew
			context:NULL];

	}
	return self;
}

- (void)dealloc
{
  [_userDefaults removeObserver:self forKeyPath:@"LastCheckDate"];
  [super dealloc];
}

- (void)observeValueForKeyPath:(NSString *)keyPath
                      ofObject:(id)object
                        change:(NSDictionary *)change
                       context:(void *)context {

	  // read the ManagedInstalls preferences
	  NSUserDefaults * userDefaults = [[NSUserDefaults alloc] initWithSuiteName:@"ManagedInstalls"];
	  NSDictionary *defaultAsDic = [userDefaults dictionaryRepresentation];

	  // serialize the dictionary to an xml plist
	  NSError *error;
	  NSPropertyListFormat format = NSPropertyListXMLFormat_v1_0;
	  NSData * data = [NSPropertyListSerialization dataWithPropertyList: defaultAsDic format:format options: 0 error: &error];

	  // pass the plist to the Go callback
	  const char * cString = [data bytes];
	  GoCallback((char*)cString);
}

@end

void ObservePreferences() {
	  uncDelegate = [[ManageHandler alloc] init];
	  [[NSUserNotificationCenter defaultUserNotificationCenter] setDelegate: uncDelegate];

	  [[NSRunLoop mainRunLoop] run];
}

*/
import "C"
import "fmt"

func main() {
	fmt.Println("start loop")
	C.ObservePreferences()
	fmt.Println("done")
}
