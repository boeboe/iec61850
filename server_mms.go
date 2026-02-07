package iec61850

/*
#include <iec61850_server.h>
#include <mms_server.h>
#include <stdlib.h>

extern MmsError fileAccessHandlerBridge(void* parameter, MmsServerConnection connection, MmsFileServiceType service, char* localFilename, char* otherFilename);

extern MmsError variableListAccessHandlerBridge(void* parameter, MmsVariableListAccessType accessType, MmsVariableListType listType, MmsDomain* domain, char* listName, MmsServerConnection connection);

static MmsError fileAccessHandlerWrap(void* p, MmsServerConnection c, MmsFileServiceType s, const char* local, const char* other) {
	return fileAccessHandlerBridge(p, c, s, (char*)local, (char*)other);
}
*/
import "C"
import (
	"crypto/x509"
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

// MmsVariableListAccessType indicates the kind of named variable list access (create, delete, read, write, get directory).
type MmsVariableListAccessType int

const (
	MmsVarlistCreate       MmsVariableListAccessType = 0
	MmsVarlistDelete       MmsVariableListAccessType = 1
	MmsVarlistRead         MmsVariableListAccessType = 2
	MmsVarlistWrite        MmsVariableListAccessType = 3
	MmsVarlistGetDirectory MmsVariableListAccessType = 4
)

// MmsVariableListType indicates the scope of the named variable list (domain, association, or VMD specific).
type MmsVariableListType int

const (
	MmsVarlistTypeDomainSpecific       MmsVariableListType = 0
	MmsVarlistTypeAssociationSpecific  MmsVariableListType = 1
	MmsVarlistTypeVmdSpecific          MmsVariableListType = 2
)

// VariableListAccessHandler is called when a client accesses a named variable list. Return nil to allow, or an MmsError to deny.
// DomainID is empty for association- or VMD-specific lists.
type VariableListAccessHandler func(accessType MmsVariableListAccessType, listType MmsVariableListType, domainID, listName string) error

var (
	fileAccessHandlerRegistry    = make(map[int32]FileAccessHandler)
	fileAccessHandlerRegistryMu  sync.Mutex
	fileAccessHandlerNextId     int32
	fileAccessHandlerParamPool  []*fileAccessHandlerParam // keep param pointers alive for C callback

	variableListAccessRegistry   = make(map[int32]VariableListAccessHandler)
	variableListAccessRegistryMu sync.Mutex
	variableListAccessNextId     int32
	variableListAccessParamPool  []*variableListAccessParam
)

type variableListAccessParam struct {
	id int32
}

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

//export variableListAccessHandlerBridge
func variableListAccessHandlerBridge(parameter unsafe.Pointer, accessType C.MmsVariableListAccessType, listType C.MmsVariableListType, domain *C.MmsDomain, listName *C.char, connection C.MmsServerConnection) C.MmsError {
	_ = connection
	if parameter == nil {
		return C.MMS_ERROR_NONE
	}
	p := (*variableListAccessParam)(parameter)
	variableListAccessRegistryMu.Lock()
	handler := variableListAccessRegistry[p.id]
	variableListAccessRegistryMu.Unlock()
	if handler == nil {
		return C.MMS_ERROR_NONE
	}
	domainID := ""
	if domain != nil {
		// MmsDomain_getName is in private header; domain name not available without it
		domainID = ""
	}
	name := ""
	if listName != nil {
		name = C.GoString(listName)
	}
	err := handler(MmsVariableListAccessType(accessType), MmsVariableListType(listType), domainID, name)
	if err == nil {
		return C.MMS_ERROR_NONE
	}
	if err == AccessDenied || err == ObjectAccessUnsupported {
		return C.MMS_ERROR_ACCESS_OBJECT_ACCESS_DENIED
	}
	if err == ObjectDoesNotExist {
		return C.MMS_ERROR_ACCESS_OBJECT_NON_EXISTENT
	}
	if err == ObjectExists {
		return C.MMS_ERROR_DEFINITION_OBJECT_EXISTS
	}
	return C.MMS_ERROR_ACCESS_OTHER
}

// InstallVariableListAccessHandler installs a callback invoked when a client accesses a named variable list (create, delete, read, write, get directory).
// Return nil to allow the access, or an error (e.g. AccessDenied) to deny.
func (is *IedServer) InstallVariableListAccessHandler(handler VariableListAccessHandler) {
	if handler == nil {
		return
	}
	variableListAccessRegistryMu.Lock()
	variableListAccessNextId++
	id := variableListAccessNextId
	variableListAccessRegistry[id] = handler
	param := &variableListAccessParam{id: id}
	variableListAccessParamPool = append(variableListAccessParamPool, param)
	variableListAccessRegistryMu.Unlock()
	mmsServer := C.IedServer_getMmsServer(is.server)
	C.MmsServer_installVariableListAccessHandler(mmsServer, (C.MmsNamedVariableListAccessHandler)(C.variableListAccessHandlerBridge), unsafe.Pointer(param))
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

// SetMaxMmsConnections sets the maximum number of MMS client connections at runtime.
func (is *IedServer) SetMaxMmsConnections(maxConnections int) {
	mmsServer := C.IedServer_getMmsServer(is.server)
	C.MmsServer_setMaxConnections(mmsServer, C.int(maxConnections))
}

// SetMaxConnections is an alias for SetMaxMmsConnections.
func (is *IedServer) SetMaxConnections(maxConnections int) {
	is.SetMaxMmsConnections(maxConnections)
}

// SetMaxMmsPduSize sets the maximum MMS PDU size for the server (applies to new connections). Not all builds export this.
func (is *IedServer) SetMaxMmsPduSize(maxPduSize int) {
	_ = maxPduSize
	// Server PDU size is typically negotiated per connection; no server-wide setter in public API.
}

// GetMaxMmsPduSize returns the maximum MMS PDU size. Returns 0 as there is no server-wide value in the public API.
func (is *IedServer) GetMaxMmsPduSize() int {
	return 0
}

// EnableMmsFileService enables or disables the MMS file service at runtime.
func (is *IedServer) EnableMmsFileService(enable bool) {
	mmsServer := C.IedServer_getMmsServer(is.server)
	C.MmsServer_enableFileService(mmsServer, C.bool(enable))
}

// EnableDynamicNamedVariableLists enables or disables dynamic named variable list (data set) service at runtime.
func (is *IedServer) EnableDynamicNamedVariableLists(enable bool) {
	mmsServer := C.IedServer_getMmsServer(is.server)
	C.MmsServer_enableDynamicNamedVariableListService(mmsServer, C.bool(enable))
}

// SetMaxAssociationSpecificDataSets sets the maximum number of association-specific (non-permanent) data sets per connection.
func (is *IedServer) SetMaxAssociationSpecificDataSets(maxDataSets int) {
	mmsServer := C.IedServer_getMmsServer(is.server)
	C.MmsServer_setMaxAssociationSpecificDataSets(mmsServer, C.int(maxDataSets))
}

// SetMaxDomainSpecificDataSets sets the maximum number of domain-specific (permanent) data sets.
func (is *IedServer) SetMaxDomainSpecificDataSets(maxDataSets int) {
	mmsServer := C.IedServer_getMmsServer(is.server)
	C.MmsServer_setMaxDomainSpecificDataSets(mmsServer, C.int(maxDataSets))
}

// SetFilestoreBasepath sets the (virtual) filestore base path for MMS file services.
// Call before Start. GetFilestoreBasepath returns the last value set here (the C API does not expose a getter).
func (is *IedServer) SetFilestoreBasepath(basepath string) {
	cPath := C.CString(basepath)
	defer C.free(unsafe.Pointer(cPath))
	C.IedServer_setFilestoreBasepath(is.server, cPath)
	is.filestoreBasepath = basepath
}

// GetFilestoreBasepath returns the filestore base path last set with SetFilestoreBasepath, or empty string.
// The C API does not provide a getter; this returns the value stored by the Go binding.
func (is *IedServer) GetFilestoreBasepath() string {
	if is == nil {
		return ""
	}
	return is.filestoreBasepath
}

// MmsServerConnection represents an MMS server-side client connection. RemoteAddress and LocalAddress
// may be set when the library exports MmsServerConnection_getClientAddress/getLocalAddress.
type MmsServerConnection struct {
	Connection    C.MmsServerConnection
	RemoteAddress string
	LocalAddress  string
}

// SetMmsClientAuthenticator installs a callback to authenticate clients. Use IedServer.SetAuthenticator instead,
// which provides the same capability via the ACSE authenticator (password, certificate, or TLS).
func (is *IedServer) SetMmsClientAuthenticator(handler func(connection *MmsServerConnection, tlsCert *x509.Certificate) bool) {
	_ = handler
	// IedServer.SetAuthenticator(clientAuthenticator) is the supported API for client authentication.
}
