package iec61850

/*
#include <iec61850_server.h>
#include <mms_server.h>
#include <stdlib.h>

extern MmsError fileAccessHandlerBridge(void* parameter, MmsServerConnection connection, MmsFileServiceType service, char* localFilename, char* otherFilename);

static MmsError fileAccessHandlerWrap(void* p, MmsServerConnection c, MmsFileServiceType s, const char* local, const char* other) {
	return fileAccessHandlerBridge(p, c, s, (char*)local, (char*)other);
}
*/
import "C"
import (
	"sync"
	"unsafe"
)

// MmsFileServiceType indicates the type of MMS file service requested.
type MmsFileServiceType int

const (
	MmsFileAccessReadDirectory MmsFileServiceType = 0
	MmsFileAccessOpen          MmsFileServiceType = 1
	MmsFileAccessObtain        MmsFileServiceType = 2
	MmsFileAccessDelete        MmsFileServiceType = 3
	MmsFileAccessRename        MmsFileServiceType = 4
)

var (
	fileAccessHandlerRegistry    = make(map[int32]FileAccessHandler)
	fileAccessHandlerRegistryMu  sync.Mutex
	fileAccessHandlerNextId     int32
	fileAccessHandlerParamPool  []*fileAccessHandlerParam // keep param pointers alive for C callback
)

// fileAccessHandlerParam holds the callback id passed to C; must stay allocated for callback lifetime.
type fileAccessHandlerParam struct {
	id int32
}

// FileAccessHandler is called when a client requests an MMS file service. Return nil to allow, or an error (e.g. MMS_ERROR_FILE_FILE_ACCESS_DENIED) to deny.
type FileAccessHandler func(service MmsFileServiceType, localFilename, otherFilename string) error

//export fileAccessHandlerBridge
func fileAccessHandlerBridge(parameter unsafe.Pointer, connection C.MmsServerConnection, service C.MmsFileServiceType, localFilename, otherFilename *C.char) C.MmsError {
	if parameter == nil {
		return C.MMS_ERROR_NONE
	}
	p := (*fileAccessHandlerParam)(parameter)
	id := p.id
	fileAccessHandlerRegistryMu.Lock()
	handler := fileAccessHandlerRegistry[id]
	fileAccessHandlerRegistryMu.Unlock()
	if handler == nil {
		return C.MMS_ERROR_NONE
	}
	local := ""
	if localFilename != nil {
		local = C.GoString(localFilename)
	}
	other := ""
	if otherFilename != nil {
		other = C.GoString(otherFilename)
	}
	err := handler(MmsFileServiceType(service), local, other)
	if err == nil {
		return C.MMS_ERROR_NONE
	}
	// Map common errors to MmsError
	if err == AccessDenied || err == ObjectAccessUnsupported {
		return C.MMS_ERROR_FILE_FILE_ACCESS_DENIED
	}
	if err == ObjectDoesNotExist {
		return C.MMS_ERROR_FILE_FILE_NON_EXISTENT
	}
	return C.MMS_ERROR_FILE_OTHER
}

// SetFileAccessHandler installs a callback that is invoked when a client requests an MMS file service. Use it to allow or deny file access.
func (is *IedServer) SetFileAccessHandler(handler FileAccessHandler) {
	if handler == nil {
		return
	}
	fileAccessHandlerRegistryMu.Lock()
	fileAccessHandlerNextId++
	id := fileAccessHandlerNextId
	fileAccessHandlerRegistry[id] = handler
	param := &fileAccessHandlerParam{id: id}
	fileAccessHandlerParamPool = append(fileAccessHandlerParamPool, param)
	fileAccessHandlerRegistryMu.Unlock()
	mmsServer := C.IedServer_getMmsServer(is.server)
	C.MmsServer_installFileAccessHandler(mmsServer, (C.MmsFileAccessHandler)(C.fileAccessHandlerWrap), unsafe.Pointer(param))
}

// SetMaxConnections sets the maximum number of TCP client connections at runtime.
// Note: This requires the library to export MmsServer_setMaxConnections (may be internal in some builds).
func (is *IedServer) SetMaxConnections(maxConnections int) {
	mmsServer := C.IedServer_getMmsServer(is.server)
	C.MmsServer_setMaxConnections(mmsServer, C.int(maxConnections))
}
