package iec61850

/*
#include <logging_api.h>
#include <iec61850_server.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// LogStorageRef wraps a C LogStorage pointer. Obtain one from C code (e.g. SqliteLogStorage_createInstance
// when the library is built with sqlite support) and wrap it with NewLogStorageRef. Call Destroy when done.
type LogStorageRef struct {
	c C.LogStorage
}

// NewLogStorageRef wraps a C LogStorage pointer. The pointer must be a valid LogStorage* from the C API
// (e.g. from SqliteLogStorage_createInstance in a build with sqlite support). Caller retains ownership
// until Destroy is called.
func NewLogStorageRef(ptr unsafe.Pointer) *LogStorageRef {
	if ptr == nil {
		return nil
	}
	return &LogStorageRef{c: (C.LogStorage)(ptr)}
}

// Destroy frees the LogStorage. Safe to call multiple times.
func (r *LogStorageRef) Destroy() {
	if r != nil && r.c != nil {
		C.LogStorage_destroy(r.c)
		r.c = nil
	}
}

// SetMaxLogEntries sets the maximum number of log entries for this log.
func (r *LogStorageRef) SetMaxLogEntries(maxEntries int) {
	if r != nil && r.c != nil {
		C.LogStorage_setMaxLogEntries(r.c, C.int(maxEntries))
	}
}

// GetMaxLogEntries returns the maximum number of log entries.
func (r *LogStorageRef) GetMaxLogEntries() int {
	if r == nil || r.c == nil {
		return 0
	}
	return int(C.LogStorage_getMaxLogEntries(r.c))
}

// AddEntry adds a log entry with the given timestamp (milliseconds since Unix epoch). Returns the entry ID.
func (r *LogStorageRef) AddEntry(timestampMs uint64) uint64 {
	if r == nil || r.c == nil {
		return 0
	}
	return uint64(C.LogStorage_addEntry(r.c, C.uint64_t(timestampMs)))
}

// AddEntryData adds data to an existing log entry.
func (r *LogStorageRef) AddEntryData(entryID uint64, dataRef string, data []byte, reasonCode uint8) bool {
	if r == nil || r.c == nil {
		return false
	}
	cRef := C.CString(dataRef)
	defer C.free(unsafe.Pointer(cRef))
	var cData *C.uint8_t
	cSize := C.int(0)
	if len(data) > 0 {
		cData = (*C.uint8_t)(unsafe.Pointer(&data[0]))
		cSize = C.int(len(data))
	}
	return bool(C.LogStorage_addEntryData(r.c, C.uint64_t(entryID), cRef, cData, cSize, C.uint8_t(reasonCode)))
}

// SetLogStorage assigns a LogStorage to a log reference (e.g. "GenericIO/LLN0$EventLog").
// The server must have a Log object at that reference. storage can be nil to clear.
func (is *IedServer) SetLogStorage(logRef string, storage *LogStorageRef) {
	if is == nil || is.server == nil {
		return
	}
	cRef := C.CString(logRef)
	defer C.free(unsafe.Pointer(cRef))
	var cStorage C.LogStorage
	if storage != nil && storage.c != nil {
		cStorage = storage.c
	}
	C.IedServer_setLogStorage(is.server, cRef, cStorage)
}
