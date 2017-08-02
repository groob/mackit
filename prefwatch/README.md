experimental example that watches a `NSUserDefaults` domain, and calls a Go callback whenever a certain key changes.

build/run:

```
go build -ldflags -s
./prefwatch
```

prefwatch watches the `LastCheckDate` key for the `ManagedInstalls` domain. Let's change it.
```
sudo defaults write /Library/Preferences/ManagedInstalls LastCheckDate foo
sudo defaults write /Library/Preferences/ManagedInstalls LastCheckDate bar
sudo defaults write /Library/Preferences/ManagedInstalls LastCheckDate baz
```

every time you change key, `prefwatch` will print out NSUserDefaults dictionary serialized to XML. The function is: 

```
//export GoCallback
func GoCallback(s *C.char) {
	fmt.Printf(C.GoString(s))
}
```


