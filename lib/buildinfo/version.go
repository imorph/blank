package buildinfo

// Version and AppNmae must be set via -ldflags '-X'
var Version, AppName string

// GetVersion returns version to caller
func GetVersion() string {
	return Version
}

// GetVersion returns version to caller
func GetAppName() string {
	return AppName
}
