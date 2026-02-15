package iec61850

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	sc "golang.org/x/text/encoding/simplifiedchinese"
)

func C2GoStr(str *C.char) string {
	utf8str, _ := sc.GB18030.NewDecoder().String(C.GoString(str))
	return utf8str
}

func Go2CStr(str string) *C.char {
	gbstr, _ := sc.GB18030.NewEncoder().String(str)
	return C.CString(gbstr)
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

// allocGo2CStr allocates a C string with GB18030 encoding and returns a cleanup function
// Usage: cStr, free := allocGo2CStr("hello"); defer free()
func allocGo2CStr(s string) (*C.char, func()) {
	cStr := Go2CStr(s)
	return cStr, func() {
		C.free(unsafe.Pointer(cStr))
	}
}
