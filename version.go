package iec61850

/*
#include "iec61850_common.h"
*/
import "C"

// GetLibraryVersion returns the version string of the underlying libiec61850 C library
func GetLibraryVersion() string {
	version := C.LibIEC61850_getVersionString()
	return C.GoString(version)
}
