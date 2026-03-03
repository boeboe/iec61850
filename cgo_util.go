package iec61850

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// C2GoStr converts a C string to a Go string (UTF-8).
func C2GoStr(str *C.char) string {
	return C.GoString(str)
}

// Go2CStr converts a Go string (UTF-8) to a C string; caller must not free the result when using allocGo2CStr.
func Go2CStr(str string) *C.char {
	return C.CString(str)
}

func C2GoBool(i C.int) bool { return i != 0 }

func Go2CBool(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

// allocCString allocates a C string and returns a cleanup function
// Usage: cStr, free := allocCString("hello"); defer free()
func allocCString(s string) (*C.char, func()) {
	cStr := C.CString(s)
	return cStr, func() {
		C.free(unsafe.Pointer(cStr))
	}
}

// allocCMalloc allocates C memory and returns a cleanup function
// Usage: ptr, free := allocCMalloc(size); defer free()
func allocCMalloc(size C.size_t) (unsafe.Pointer, func()) {
	ptr := C.malloc(size)
	return ptr, func() {
		C.free(ptr)
	}
}

// allocGo2CStr allocates a C string (UTF-8) and returns a cleanup function.
// Usage: cStr, free := allocGo2CStr("hello"); defer free()
func allocGo2CStr(s string) (*C.char, func()) {
	return allocCString(s)
}
