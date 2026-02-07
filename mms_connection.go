package iec61850

/*
#include <stdlib.h>
#include <mms_client_connection.h>
#include <mms_value.h>
#include <iso_connection_parameters.h>
#include <tls_config.h>
#include <string.h>

extern void mmsConnectionStateChangedBridge(MmsConnection connection, void* parameter, MmsConnectionState newState);
extern void readVariableAsyncBridge(uint32_t invokeId, void* parameter, MmsError mmsError, MmsValue* value);
extern void writeVariableAsyncBridge(uint32_t invokeId, void* parameter, MmsError mmsError, MmsDataAccessError accessError);
extern void rawMessageHandlerBridge(void* parameter, uint8_t* message, int messageLength, _Bool received);
extern void fileDirectoryAsyncBridge(uint32_t invokeId, void* parameter, MmsError mmsError, char* filename, uint32_t size, uint64_t lastModified, _Bool moreFollows);
extern void genericServiceAsyncBridge(uint32_t invokeId, void* parameter, MmsError mmsError, _Bool success);
extern void readNVLDirectoryAsyncBridge(uint32_t invokeId, void* parameter, MmsError mmsError, LinkedList specs, _Bool deletable);
extern void getVariableAccessAttributesAsyncBridge(uint32_t invokeId, void* parameter, MmsError mmsError, MmsVariableSpecification* spec);
extern void identifyAsyncBridge(uint32_t invokeId, void* parameter, MmsError mmsError, char* vendorName, char* modelName, char* revision);
extern void readJournalAsyncBridge(uint32_t invokeId, void* parameter, MmsError mmsError, LinkedList journalEntries, _Bool moreFollows);
extern void concludeAbortBridge(void* parameter, MmsError mmsError, _Bool success);
extern void getNameListAsyncBridge(uint32_t invokeId, void* parameter, MmsError mmsError, LinkedList nameList, _Bool moreFollows);
extern void informationReportBridge(void* parameter, char* domainName, char* variableListName, MmsValue* value, _Bool isVariableListName);

static void destroyMmsValueLinkedListLocal(LinkedList L) {
	if (L) LinkedList_destroyDeep(L, (LinkedListValueDeleteFunction)MmsValue_delete);
}
static void destroyJournalEntryLinkedListLocal(LinkedList L) {
	if (L) LinkedList_destroyDeep(L, (LinkedListValueDeleteFunction)MmsJournalEntry_destroy);
}
static void destroyCharPtrLinkedList(LinkedList L) {
	if (L) LinkedList_destroyDeep(L, (LinkedListValueDeleteFunction)free);
}
*/
import "C"
import (
	"sync"
	"unsafe"
)

// TLSConfiguration holds TLS settings and in-memory certificates for secure MMS connections.
type TLSConfiguration struct {
	ChainValidation      bool     // Enable certificate chain validation
	AllowOnlyKnownCerts  bool     // Allow only known (explicitly added) certificates
	CACertificates       [][]byte // CA certificates (PEM or DER)
	OwnCertificate       []byte   // Own certificate (PEM or DER)
	OwnKey               []byte   // Own private key (PEM; optional password empty)
}

// MmsConnection wraps a C MmsConnection for standalone MMS client use (without IedConnection).
// Call Destroy when done.
type MmsConnection struct {
	c      C.MmsConnection
	connMu sync.Mutex
}

// rawMessageCtx holds the Go callback for raw message handling; stored in rawMessageRegistry.
type rawMessageCtx struct {
	callback func(message []byte, received bool)
}

var (
	rawMessageRegistry   = make(map[C.MmsConnection]*rawMessageCtx)
	rawMessageRegistryMu sync.Mutex

	infoReportRegistry   = make(map[C.MmsConnection]func(domainName, variableListName string, value *MmsValue, isVariableListName bool))
	infoReportRegistryMu sync.Mutex
)

// fileDirAsyncCtx holds state for an in-flight FileDirectoryAsync; keyed by connection in fileDirAsyncRegistry.
type fileDirAsyncCtx struct {
	entries  []MmsFileDirectoryEntryEx
	callback func(entries []MmsFileDirectoryEntryEx, moreFollows bool, err error)
}

var (
	fileDirAsyncRegistry   = make(map[C.MmsConnection]*fileDirAsyncCtx)
	fileDirAsyncRegistryMu sync.Mutex
)

// NewMmsConnection creates a new non-TLS MmsConnection. Call Destroy when done.
func NewMmsConnection() *MmsConnection {
	return &MmsConnection{c: C.MmsConnection_create()}
}

// NewMmsConnectionSecure creates a new TLS-enabled MmsConnection using the given TLS configuration.
// Call Destroy when done.
func NewMmsConnectionSecure(tlsConfig *TLSConfiguration) *MmsConnection {
	if tlsConfig == nil {
		return &MmsConnection{c: C.MmsConnection_create()}
	}
	cTls := buildCTLSConfiguration(tlsConfig)
	if cTls == nil {
		return nil
	}
	defer C.TLSConfiguration_destroy(cTls)
	return &MmsConnection{c: C.MmsConnection_createSecure(cTls)}
}

// NewMmsConnectionNonThreaded creates an MmsConnection that does not use a background thread; call Tick() periodically. tlsConfig may be nil for non-TLS.
func NewMmsConnectionNonThreaded(tlsConfig *TLSConfiguration) *MmsConnection {
	var cTls C.TLSConfiguration
	if tlsConfig != nil {
		cTls = buildCTLSConfiguration(tlsConfig)
		if cTls == nil {
			return nil
		}
		defer C.TLSConfiguration_destroy(cTls)
	}
	return &MmsConnection{c: C.MmsConnection_createNonThreaded(cTls)}
}

func buildCTLSConfiguration(t *TLSConfiguration) C.TLSConfiguration {
	cTls := C.TLSConfiguration_create()
	if cTls == nil {
		return nil
	}
	C.TLSConfiguration_setClientMode(cTls)
	C.TLSConfiguration_setChainValidation(cTls, C.bool(t.ChainValidation))
	C.TLSConfiguration_setAllowOnlyKnownCertificates(cTls, C.bool(t.AllowOnlyKnownCerts))
	if len(t.OwnCertificate) > 0 {
		var certPtr *C.uint8_t
		if len(t.OwnCertificate) > 0 {
			certPtr = (*C.uint8_t)(unsafe.Pointer(&t.OwnCertificate[0]))
		}
		C.TLSConfiguration_setOwnCertificate(cTls, certPtr, C.int(len(t.OwnCertificate)))
	}
	if len(t.OwnKey) > 0 {
		var keyPtr *C.uint8_t
		if len(t.OwnKey) > 0 {
			keyPtr = (*C.uint8_t)(unsafe.Pointer(&t.OwnKey[0]))
		}
		C.TLSConfiguration_setOwnKey(cTls, keyPtr, C.int(len(t.OwnKey)), nil)
	}
	for _, ca := range t.CACertificates {
		if len(ca) == 0 {
			continue
		}
		caPtr := (*C.uint8_t)(unsafe.Pointer(&ca[0]))
		C.TLSConfiguration_addCACertificate(cTls, caPtr, C.int(len(ca)))
	}
	return cTls
}

// Destroy releases the C MmsConnection. Safe to call multiple times.
func (c *MmsConnection) Destroy() {
	c.connMu.Lock()
	conn := c.c
	c.c = nil
	c.connMu.Unlock()
	if conn != nil {
		mmsConnAsyncRegistryLock.Lock()
		delete(mmsConnAsyncRegistry, conn)
		mmsConnAsyncRegistryLock.Unlock()
		rawMessageRegistryMu.Lock()
		delete(rawMessageRegistry, conn)
		rawMessageRegistryMu.Unlock()
		fileDirAsyncRegistryMu.Lock()
		delete(fileDirAsyncRegistry, conn)
		fileDirAsyncRegistryMu.Unlock()
		C.MmsConnection_destroy(conn)
	}
}

// SetConnectTimeout sets the connect timeout in milliseconds.
func (c *MmsConnection) SetConnectTimeout(timeoutMs uint32) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c != nil {
		C.MmsConnection_setConnectTimeout(c.c, C.uint32_t(timeoutMs))
	}
}

// GetConnectTimeout returns the connect timeout in milliseconds. The C library does not expose a getter; returns 0.
func (c *MmsConnection) GetConnectTimeout() uint32 {
	return 0
}

// SetRequestTimeout sets the request timeout in milliseconds.
func (c *MmsConnection) SetRequestTimeout(timeoutMs uint32) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c != nil {
		C.MmsConnection_setRequestTimeout(c.c, C.uint32_t(timeoutMs))
	}
}

// GetRequestTimeout returns the request timeout in milliseconds.
func (c *MmsConnection) GetRequestTimeout() uint32 {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return 0
	}
	return uint32(C.MmsConnection_getRequestTimeout(c.c))
}

// SetMaxOutstandingCalls sets the maximum outstanding calls (calling and called) for this connection.
func (c *MmsConnection) SetMaxOutstandingCalls(calling, called int) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c != nil {
		C.MmsConnnection_setMaxOutstandingCalls(c.c, C.int(calling), C.int(called))
	}
}

// SetLocalDetail sets the maximum MMS PDU size (local detail) for this connection.
func (c *MmsConnection) SetLocalDetail(localDetail int32) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c != nil {
		C.MmsConnection_setLocalDetail(c.c, C.int32_t(localDetail))
	}
}

// GetLocalDetail returns the maximum MMS PDU size (local detail) for this connection.
func (c *MmsConnection) GetLocalDetail() int32 {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return 0
	}
	return int32(C.MmsConnection_getLocalDetail(c.c))
}

// GetIsoConnectionParameters returns a copy of the ISO connection parameters (selectors, AP titles). Returns nil if connection is invalid.
func (c *MmsConnection) GetIsoConnectionParameters() *IsoConnectionParameters {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return nil
	}
	p := C.MmsConnection_getIsoConnectionParameters(c.c)
	if p == nil {
		return nil
	}
	out := &IsoConnectionParameters{}
	out.LocalAeQualifier = int32(p.localAEQualifier)
	out.RemoteAeQualifier = int32(p.remoteAEQualifier)
	out.LocalApTitle = copyApTitle(p.localApTitle, p.localApTitleLen)
	out.RemoteApTitle = copyApTitle(p.remoteApTitle, p.remoteApTitleLen)
	out.LocalTSelector = tSelectorToSlice(p.localTSelector)
	out.LocalSSelector = sSelectorToSlice(p.localSSelector)
	out.LocalPSelector = pSelectorToSlice(p.localPSelector)
	out.RemoteTSelector = tSelectorToSlice(p.remoteTSelector)
	out.RemoteSSelector = sSelectorToSlice(p.remoteSSelector)
	out.RemotePSelector = pSelectorToSlice(p.remotePSelector)
	return out
}

// SetFilestoreBasepath sets the virtual filestore basepath for MMS file services (client side). Requires CONFIG_SET_FILESTORE_BASEPATH_AT_RUNTIME in the C library.
func (c *MmsConnection) SetFilestoreBasepath(basepath string) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return
	}
	cPath, freecPath := allocCString(basepath)
	defer freecPath()
	C.MmsConnection_setFilestoreBasepath(c.c, cPath)
}

// SetInformationReportHandler sets the handler for MMS information reports (unsolicited updates). Pass nil to clear.
func (c *MmsConnection) SetInformationReportHandler(callback func(domainName, variableListName string, value *MmsValue, isVariableListName bool)) {
	c.connMu.Lock()
	conn := c.c
	if conn != nil && callback != nil {
		infoReportRegistryMu.Lock()
		infoReportRegistry[conn] = callback
		infoReportRegistryMu.Unlock()
		C.MmsConnection_setInformationReportHandler(conn, (C.MmsInformationReportHandler)(C.informationReportBridge), unsafe.Pointer(conn))
	} else if conn != nil {
		C.MmsConnection_setInformationReportHandler(conn, nil, nil)
		infoReportRegistryMu.Lock()
		delete(infoReportRegistry, conn)
		infoReportRegistryMu.Unlock()
	}
	c.connMu.Unlock()
}

// SetRawMessageHandler sets a callback that receives every raw MMS message (sent or received). Pass nil to clear.
// The callback may be invoked from a different goroutine; received is true for incoming, false for outgoing.
func (c *MmsConnection) SetRawMessageHandler(callback func(message []byte, received bool)) {
	c.connMu.Lock()
	conn := c.c
	c.connMu.Unlock()
	if conn == nil {
		return
	}
	rawMessageRegistryMu.Lock()
	if callback == nil {
		delete(rawMessageRegistry, conn)
		rawMessageRegistryMu.Unlock()
		C.MmsConnection_setRawMessageHandler(conn, nil, nil)
		return
	}
	rawMessageRegistry[conn] = &rawMessageCtx{callback: callback}
	rawMessageRegistryMu.Unlock()
	C.MmsConnection_setRawMessageHandler(conn, (C.MmsRawMessageHandler)(C.rawMessageHandlerBridge), unsafe.Pointer(conn))
}

func copyApTitle(arr [10]C.uint8_t, len C.int) []byte {
	n := int(len)
	if n <= 0 {
		return nil
	}
	if n > 10 {
		n = 10
	}
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = byte(arr[i])
	}
	return out
}

func tSelectorToSlice(s C.TSelector) []byte {
	n := int(s.size)
	if n <= 0 {
		return nil
	}
	if n > 4 {
		n = 4
	}
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = byte(s.value[i])
	}
	return out
}

func sSelectorToSlice(s C.SSelector) []byte {
	n := int(s.size)
	if n <= 0 {
		return nil
	}
	if n > 16 {
		n = 16
	}
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = byte(s.value[i])
	}
	return out
}

func pSelectorToSlice(s C.PSelector) []byte {
	n := int(s.size)
	if n <= 0 {
		return nil
	}
	if n > 16 {
		n = 16
	}
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = byte(s.value[i])
	}
	return out
}

// ConnectAsync starts a non-blocking connection. The callback is invoked with nil when connected or with an error on failure/close.
func (c *MmsConnection) ConnectAsync(hostname string, port int, callback func(error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	// Register callback for state changes before starting connect
	ctx := &mmsConnAsyncContext{callback: callback}
	mmsConnAsyncRegistryLock.Lock()
	mmsConnAsyncRegistry[conn] = ctx
	mmsConnAsyncRegistryLock.Unlock()
	C.MmsConnection_setConnectionStateChangedHandler(conn, (C.MmsConnectionStateChangedHandler)(C.mmsConnectionStateChangedBridge), unsafe.Pointer(ctx))
	c.connMu.Unlock()

	host, freehost := allocCString(hostname)
	defer freehost()
	var cError C.MmsError
	C.MmsConnection_connectAsync(conn, &cError, host, C.int(port))
	if err := GetMmsError(cError); err != nil {
		mmsConnAsyncRegistryLock.Lock()
		delete(mmsConnAsyncRegistry, conn)
		mmsConnAsyncRegistryLock.Unlock()
		if callback != nil {
			callback(err)
		}
		return err
	}
	return nil
}

type mmsConnAsyncContext struct {
	callback func(error)
}

var (
	mmsConnAsyncRegistry     = make(map[C.MmsConnection]*mmsConnAsyncContext)
	mmsConnAsyncRegistryLock sync.Mutex
)

//export mmsConnectionStateChangedBridge
func mmsConnectionStateChangedBridge(connection C.MmsConnection, parameter unsafe.Pointer, newState C.MmsConnectionState) {
	if parameter == nil {
		return
	}
	ctx := (*mmsConnAsyncContext)(parameter)
	mmsConnAsyncRegistryLock.Lock()
	_, ok := mmsConnAsyncRegistry[connection]
	if ok {
		delete(mmsConnAsyncRegistry, connection)
	}
	mmsConnAsyncRegistryLock.Unlock()
	cb := ctx.callback
	if cb == nil {
		return
	}
	switch newState {
	case C.MMS_CONNECTION_STATE_CONNECTED:
		cb(nil)
	case C.MMS_CONNECTION_STATE_CLOSED, C.MMS_CONNECTION_STATE_CLOSING:
		cb(ConnectionLost)
	}
}

// SetIsoConnectionParameters sets the ISO layer selectors (T, S, P) for local and remote.
// Pass nil or empty slice for a selector to leave it unchanged or empty.
// T selector max 4 bytes, S and P selectors max 16 bytes.
func (c *MmsConnection) SetIsoConnectionParameters(
	localTSelector, localSSelector, localPSelector []byte,
	remoteTSelector, remoteSSelector, remotePSelector []byte,
) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return
	}
	isoParams := C.MmsConnection_getIsoConnectionParameters(c.c)
	if isoParams == nil {
		return
	}
	localT := fillTSelector(localTSelector)
	localS := fillSSelector(localSSelector)
	localP := fillPSelector(localPSelector)
	remoteT := fillTSelector(remoteTSelector)
	remoteS := fillSSelector(remoteSSelector)
	remoteP := fillPSelector(remotePSelector)
	C.IsoConnectionParameters_setLocalAddresses(isoParams, localP, localS, localT)
	C.IsoConnectionParameters_setRemoteAddresses(isoParams, remoteP, remoteS, remoteT)
}

func fillTSelector(src []byte) C.TSelector {
	var s C.TSelector
	if len(src) > 0 {
		n := len(src)
		if n > 4 {
			n = 4
		}
		s.size = C.uint8_t(n)
		for i := 0; i < n; i++ {
			s.value[i] = C.uint8_t(src[i])
		}
	}
	return s
}

func fillSSelector(src []byte) C.SSelector {
	var s C.SSelector
	if len(src) > 0 {
		n := len(src)
		if n > 16 {
			n = 16
		}
		s.size = C.uint8_t(n)
		for i := 0; i < n; i++ {
			s.value[i] = C.uint8_t(src[i])
		}
	}
	return s
}

func fillPSelector(src []byte) C.PSelector {
	var p C.PSelector
	if len(src) > 0 {
		n := len(src)
		if n > 16 {
			n = 16
		}
		p.size = C.uint8_t(n)
		for i := 0; i < n; i++ {
			p.value[i] = C.uint8_t(src[i])
		}
	}
	return p
}

// GetMmsConnectionParameters returns the MMS connection parameters (max outstanding calls, PDU size, etc.).
// Returns nil if the connection is invalid.
func (c *MmsConnection) GetMmsConnectionParameters() *MmsConnectionParameters {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return nil
	}
	p := C.MmsConnection_getMmsConnectionParameters(c.c)
	var sv [11]uint8
	for i := 0; i < 11; i++ {
		sv[i] = uint8(p.servicesSupported[i])
	}
	return &MmsConnectionParameters{
		MaxServOutstandingCalling: int32(p.maxServOutstandingCalling),
		MaxServOutstandingCalled:  int32(p.maxServOutstandingCalled),
		DataStructureNestingLevel: int32(p.dataStructureNestingLevel),
		MaxPduSize:                int32(p.maxPduSize),
		ServicesSupported:        sv,
	}
}

// Conclude sends the MMS conclude service to orderly close the association with the server.
// The connection is closed; use Connect again to reconnect.
func (c *MmsConnection) Conclude() error {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return NotConnected
	}
	var cError C.MmsError
	C.MmsConnection_conclude(c.c, &cError)
	return GetMmsError(cError)
}

// Tick processes connection events for non-threaded mode. Call periodically. Returns true if more work is pending.
func (c *MmsConnection) Tick() bool {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return false
	}
	return bool(C.MmsConnection_tick(c.c))
}

// AbortAsync sends the MMS abort service asynchronously. The C API does not provide a completion callback; it returns immediately.
func (c *MmsConnection) AbortAsync() error {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return NotConnected
	}
	var cError C.MmsError
	C.MmsConnection_abortAsync(c.c, &cError)
	return GetMmsError(cError)
}

// ConcludeAsync sends the MMS conclude service asynchronously. The callback is invoked when the operation completes (may run from another goroutine).
func (c *MmsConnection) ConcludeAsync(callback func(error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	if callback == nil {
		return UserProvidedInvalidArgument
	}
	ctx := &concludeAbortCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_concludeAsync(conn, &cError, (C.MmsConnection_ConcludeAbortHandler)(C.concludeAbortBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// --- Async read/write context and bridges ---

type readVarAsyncCtx struct {
	callback func(*MmsValue, error)
}

type writeVarAsyncCtx struct {
	callback func(error)
}

//export readVariableAsyncBridge
func readVariableAsyncBridge(invokeId C.uint32_t, parameter unsafe.Pointer, mmsError C.MmsError, value *C.MmsValue) {
	if parameter == nil {
		return
	}
	ctx := (*readVarAsyncCtx)(parameter)
	cb := ctx.callback
	if cb == nil {
		if value != nil {
			C.MmsValue_delete(value)
		}
		return
	}
	err := GetMmsError(mmsError)
	var goVal *MmsValue
	if value != nil {
		goVal = CMmsValueToMmsValue(value)
		C.MmsValue_delete(value)
	}
	cb(goVal, err)
}

//export writeVariableAsyncBridge
func writeVariableAsyncBridge(invokeId C.uint32_t, parameter unsafe.Pointer, mmsError C.MmsError, accessError C.MmsDataAccessError) {
	_ = accessError
	if parameter == nil {
		return
	}
	ctx := (*writeVarAsyncCtx)(parameter)
	cb := ctx.callback
	if cb != nil {
		cb(GetMmsError(mmsError))
	}
}

type genericServiceAsyncCtx struct {
	callback func(error)
}

//export genericServiceAsyncBridge
func genericServiceAsyncBridge(invokeId C.uint32_t, parameter unsafe.Pointer, mmsError C.MmsError, success C._Bool) {
	_ = invokeId
	if parameter == nil {
		return
	}
	ctx := (*genericServiceAsyncCtx)(parameter)
	if ctx.callback != nil {
		err := GetMmsError(mmsError)
		if err == nil && !bool(success) {
			err = Unknown
		}
		ctx.callback(err)
	}
}

type concludeAbortCtx struct {
	callback func(error)
}

//export concludeAbortBridge
func concludeAbortBridge(parameter unsafe.Pointer, mmsError C.MmsError, success C._Bool) {
	_ = success
	if parameter == nil {
		return
	}
	ctx := (*concludeAbortCtx)(parameter)
	if ctx.callback != nil {
		ctx.callback(GetMmsError(mmsError))
	}
}

type getNameListAsyncCtx struct {
	callback func(names []string, moreFollows bool, err error)
}

//export informationReportBridge
func informationReportBridge(parameter unsafe.Pointer, domainName *C.char, variableListName *C.char, value *C.MmsValue, isVariableListName C._Bool) {
	if parameter == nil {
		if value != nil {
			C.MmsValue_delete(value)
		}
		return
	}
	conn := (C.MmsConnection)(parameter)
	infoReportRegistryMu.Lock()
	cb := infoReportRegistry[conn]
	infoReportRegistryMu.Unlock()
	if cb == nil {
		if value != nil {
			C.MmsValue_delete(value)
		}
		return
	}
	d, v := "", ""
	if domainName != nil {
		d = C.GoString(domainName)
	}
	if variableListName != nil {
		v = C.GoString(variableListName)
	}
	var goVal *MmsValue
	if value != nil {
		goVal = CMmsValueToMmsValue(value)
		C.MmsValue_delete(value)
	}
	cb(d, v, goVal, bool(isVariableListName))
}

//export getNameListAsyncBridge
func getNameListAsyncBridge(invokeId C.uint32_t, parameter unsafe.Pointer, mmsError C.MmsError, nameList C.LinkedList, moreFollows C._Bool) {
	_ = invokeId
	if parameter == nil {
		if nameList != nil {
			C.destroyCharPtrLinkedList(nameList)
		}
		return
	}
	ctx := (*getNameListAsyncCtx)(parameter)
	cb := ctx.callback
	if cb == nil {
		if nameList != nil {
			C.destroyCharPtrLinkedList(nameList)
		}
		return
	}
	err := GetMmsError(mmsError)
	var names []string
	if err == nil && nameList != nil {
		for node := nameList; node != nil; node = C.LinkedList_getNext(node) {
			data := C.LinkedList_getData(node)
			if data != nil {
				names = append(names, C.GoString((*C.char)(data)))
			}
		}
		C.destroyCharPtrLinkedList(nameList)
	}
	cb(names, bool(moreFollows), err)
}

type readNVLDirectoryAsyncCtx struct {
	callback func(specs []MmsVariableAccessSpec, deletable bool, err error)
}

//export readNVLDirectoryAsyncBridge
func readNVLDirectoryAsyncBridge(invokeId C.uint32_t, parameter unsafe.Pointer, mmsError C.MmsError, specs C.LinkedList, deletable C._Bool) {
	_ = invokeId
	if parameter == nil {
		if specs != nil {
			C.LinkedList_destroyDeep(specs, (C.LinkedListValueDeleteFunction)(C.MmsVariableAccessSpecification_destroy))
		}
		return
	}
	ctx := (*readNVLDirectoryAsyncCtx)(parameter)
	var out []MmsVariableAccessSpec
	if specs != nil {
		for node := specs; node != nil; node = C.LinkedList_getNext(node) {
			data := C.LinkedList_getData(node)
			if data != nil {
				spec := (*C.MmsVariableAccessSpecification)(data)
				out = append(out, MmsVariableAccessSpec{
					DomainID: C.GoString(spec.domainId),
					ItemID:   C.GoString(spec.itemId),
				})
			}
		}
		C.LinkedList_destroyDeep(specs, (C.LinkedListValueDeleteFunction)(C.MmsVariableAccessSpecification_destroy))
	}
	if ctx.callback != nil {
		ctx.callback(out, bool(deletable), GetMmsError(mmsError))
	}
}

type getVariableAccessAttributesAsyncCtx struct {
	callback func(*MmsVariableSpecificationRef, error)
}

//export getVariableAccessAttributesAsyncBridge
func getVariableAccessAttributesAsyncBridge(invokeId C.uint32_t, parameter unsafe.Pointer, mmsError C.MmsError, spec *C.MmsVariableSpecification) {
	_ = invokeId
	if parameter == nil {
		if spec != nil {
			C.MmsVariableSpecification_destroy(spec)
		}
		return
	}
	ctx := (*getVariableAccessAttributesAsyncCtx)(parameter)
	if ctx.callback != nil {
		err := GetMmsError(mmsError)
		if err != nil {
			if spec != nil {
				C.MmsVariableSpecification_destroy(spec)
			}
			ctx.callback(nil, err)
			return
		}
		var ref *MmsVariableSpecificationRef
		if spec != nil {
			ref = &MmsVariableSpecificationRef{c: spec, owned: true, libraryOwned: true}
		}
		ctx.callback(ref, nil)
	}
}

type identifyAsyncCtx struct {
	callback func(vendorName, modelName, revision string, err error)
}

type readJournalAsyncCtx struct {
	callback func(entries []*MmsJournalEntry, moreFollows bool, err error)
}

//export readJournalAsyncBridge
func readJournalAsyncBridge(invokeId C.uint32_t, parameter unsafe.Pointer, mmsError C.MmsError, journalEntries C.LinkedList, moreFollows C._Bool) {
	_ = invokeId
	if parameter == nil {
		if journalEntries != nil {
			C.destroyJournalEntryLinkedListLocal(journalEntries)
		}
		return
	}
	ctx := (*readJournalAsyncCtx)(parameter)
	cb := ctx.callback
	if cb == nil {
		if journalEntries != nil {
			C.destroyJournalEntryLinkedListLocal(journalEntries)
		}
		return
	}
	err := GetMmsError(mmsError)
	var entries []*MmsJournalEntry
	if err == nil && journalEntries != nil {
		for node := journalEntries; node != nil; node = C.LinkedList_getNext(node) {
			data := C.LinkedList_getData(node)
			if data != nil {
				e := convertCJournalEntryToMms(C.MmsJournalEntry(data))
				entries = append(entries, &e)
			}
		}
		C.destroyJournalEntryLinkedListLocal(journalEntries)
	}
	cb(entries, bool(moreFollows), err)
}

//export identifyAsyncBridge
func identifyAsyncBridge(invokeId C.uint32_t, parameter unsafe.Pointer, mmsError C.MmsError, vendorName *C.char, modelName *C.char, revision *C.char) {
	_ = invokeId
	if parameter == nil {
		return
	}
	ctx := (*identifyAsyncCtx)(parameter)
	if ctx.callback != nil {
		v, m, r := "", "", ""
		if vendorName != nil {
			v = C.GoString(vendorName)
		}
		if modelName != nil {
			m = C.GoString(modelName)
		}
		if revision != nil {
			r = C.GoString(revision)
		}
		ctx.callback(v, m, r, GetMmsError(mmsError))
	}
}

//export rawMessageHandlerBridge
func rawMessageHandlerBridge(parameter unsafe.Pointer, message *C.uint8_t, messageLength C.int, received C._Bool) {
	if parameter == nil {
		return
	}
	conn := C.MmsConnection(parameter)
	rawMessageRegistryMu.Lock()
	ctx := rawMessageRegistry[conn]
	rawMessageRegistryMu.Unlock()
	if ctx == nil || ctx.callback == nil {
		return
	}
	n := int(messageLength)
	var msg []byte
	if message != nil && n > 0 {
		msg = C.GoBytes(unsafe.Pointer(message), messageLength)
	}
	ctx.callback(msg, bool(received))
}

//export fileDirectoryAsyncBridge
func fileDirectoryAsyncBridge(invokeId C.uint32_t, parameter unsafe.Pointer, mmsError C.MmsError, filename *C.char, size C.uint32_t, lastModified C.uint64_t, moreFollows C._Bool) {
	_ = invokeId
	if parameter == nil {
		return
	}
	conn := C.MmsConnection(parameter)
	fileDirAsyncRegistryMu.Lock()
	ctx := fileDirAsyncRegistry[conn]
	fileDirAsyncRegistryMu.Unlock()
	if ctx == nil || ctx.callback == nil {
		return
	}
	err := GetMmsError(mmsError)
	if err != nil {
		fileDirAsyncRegistryMu.Lock()
		delete(fileDirAsyncRegistry, conn)
		fileDirAsyncRegistryMu.Unlock()
		ctx.callback(nil, false, err)
		return
	}
	if filename != nil {
		ctx.entries = append(ctx.entries, MmsFileDirectoryEntryEx{
			Filename:         C.GoString(filename),
			FileSize:         uint32(size),
			LastModifiedTime: uint64(lastModified),
		})
	}
	if filename == nil || !bool(moreFollows) {
		fileDirAsyncRegistryMu.Lock()
		delete(fileDirAsyncRegistry, conn)
		fileDirAsyncRegistryMu.Unlock()
		entries := make([]MmsFileDirectoryEntryEx, len(ctx.entries))
		copy(entries, ctx.entries)
		ctx.callback(entries, bool(moreFollows), nil)
	}
}

// ReadVariableAsync reads a single variable asynchronously. The callback receives the value and any error.
// On success the value is non-nil (caller does not free it). The callback may be invoked from a different goroutine.
func (c *MmsConnection) ReadVariableAsync(domainID, itemID string, callback func(*MmsValue, error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	cDomain, freecDomain := allocCString(domainID)
	defer freecDomain()
	cItem, freecItem := allocCString(itemID)
	defer freecItem()
	ctx := &readVarAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_readVariableAsync(conn, nil, &cError, cDomain, cItem, (C.MmsConnection_ReadVariableHandler)(C.readVariableAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// WriteVariableAsync writes a single variable asynchronously. value must be a valid MmsValueRef (C-backed).
// The callback is invoked with nil on success or an error. It may be invoked from a different goroutine.
func (c *MmsConnection) WriteVariableAsync(domainID, itemID string, value *MmsValueRef, callback func(error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	if value == nil || value.c == nil {
		if callback != nil {
			callback(UserProvidedInvalidArgument)
		}
		return UserProvidedInvalidArgument
	}
	cDomain, freecDomain := allocCString(domainID)
	defer freecDomain()
	cItem, freecItem := allocCString(itemID)
	defer freecItem()
	ctx := &writeVarAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_writeVariableAsync(conn, nil, &cError, cDomain, cItem, value.c, (C.MmsConnection_WriteVariableHandler)(C.writeVariableAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// DefineNamedVariableListAsync defines a new domain or VMD scoped named variable list asynchronously.
// Pass domainID as "" for VMD scope. The callback is invoked on completion; it may run from another goroutine.
func (c *MmsConnection) DefineNamedVariableListAsync(domainID, listName string, variableSpecs []VariableAccessSpec, callback func(error)) error {
	if len(variableSpecs) == 0 {
		if callback != nil {
			callback(UserProvidedInvalidArgument)
		}
		return UserProvidedInvalidArgument
	}
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	var cDomain *C.char
	var freeCDomain func()
	if domainID != "" {
		cDomain, freeCDomain = allocCString(domainID)
		defer freeCDomain()
	}
	cList, freecList := allocCString(listName)
	defer freecList()
	clist := C.LinkedList_create()
	for _, vs := range variableSpecs {
		cDom := C.CString(vs.DomainID)
		cItem := C.CString(vs.ItemID)
		var spec *C.MmsVariableAccessSpecification
		if vs.ArrayIndex >= 0 || vs.ComponentName != "" {
			var cComp *C.char
			if vs.ComponentName != "" {
				cComp = C.CString(vs.ComponentName)
			}
			spec = C.MmsVariableAccessSpecification_createAlternateAccess(cDom, cItem, C.int32_t(vs.ArrayIndex), cComp)
		} else {
			spec = C.MmsVariableAccessSpecification_create(cDom, cItem)
		}
		C.LinkedList_add(clist, unsafe.Pointer(spec))
	}
	// C layer takes ownership of clist and its elements for the async request
	ctx := &genericServiceAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_defineNamedVariableListAsync(conn, nil, &cError, cDomain, cList, clist, (C.MmsConnection_GenericServiceHandler)(C.genericServiceAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// ReadNamedVariableListDirectoryAsync reads the entry list of a named variable list asynchronously.
// Pass domainID as "" for VMD scope. The callback receives the variable specs and deletable flag; it may run from another goroutine.
func (c *MmsConnection) ReadNamedVariableListDirectoryAsync(domainID, listName string, callback func(specs []MmsVariableAccessSpec, deletable bool, err error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, false, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	var cDomain *C.char
	var freeCDomain func()
	if domainID != "" {
		cDomain, freeCDomain = allocCString(domainID)
		defer freeCDomain()
	}
	cList, freecList := allocCString(listName)
	defer freecList()
	ctx := &readNVLDirectoryAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_readNamedVariableListDirectoryAsync(conn, nil, &cError, cDomain, cList, (C.MmsConnection_ReadNVLDirectoryHandler)(C.readNVLDirectoryAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// GetVariableAccessAttributesAsync retrieves the variable access attributes (type specification) asynchronously.
// On success the callback receives a non-nil MmsVariableSpecificationRef; the caller must call Free() on it when done.
// The callback may run from another goroutine.
func (c *MmsConnection) GetVariableAccessAttributesAsync(domainID, itemID string, callback func(*MmsVariableSpecificationRef, error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	cDomain, freecDomain := allocCString(domainID)
	defer freecDomain()
	cItem, freecItem := allocCString(itemID)
	defer freecItem()
	ctx := &getVariableAccessAttributesAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_getVariableAccessAttributesAsync(conn, nil, &cError, cDomain, cItem, (C.MmsConnection_GetVariableAccessAttributesHandler)(C.getVariableAccessAttributesAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// IdentifyAsync retrieves the server identity (vendor, model, revision) asynchronously.
// The callback may run from another goroutine.
func (c *MmsConnection) IdentifyAsync(callback func(vendorName, modelName, revision string, err error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback("", "", "", NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	ctx := &identifyAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_identifyAsync(conn, nil, &cError, (C.MmsConnection_IdentifyHandler)(C.identifyAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// WriteMultipleVariables writes multiple variables in one request. items and values must have the same length; values are C-backed MmsValueRefs.
// If accessResults is non-nil it is filled with the per-variable data access error results.
func (c *MmsConnection) WriteMultipleVariables(domainID string, items []string, values []*MmsValueRef, accessResults *[]MmsDataAccessError) error {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return NotConnected
	}
	if len(items) != len(values) {
		return UserProvidedInvalidArgument
	}
	if len(items) == 0 {
		return nil
	}
	cDomain, freecDomain := allocCString(domainID)
	defer freecDomain()
	itemsList := C.LinkedList_create()
	defer C.LinkedList_destroyDeep(itemsList, (C.LinkedListValueDeleteFunction)(C.free))
	valuesList := C.LinkedList_create()
	for _, id := range items {
		cItem := C.CString(id)
		C.LinkedList_add(itemsList, unsafe.Pointer(cItem))
	}
	for _, v := range values {
		if v != nil && v.c != nil {
			C.LinkedList_add(valuesList, unsafe.Pointer(v.c))
		}
	}
	defer C.LinkedList_destroyStatic(valuesList)
	var cError C.MmsError
	var cResults C.LinkedList
	C.MmsConnection_writeMultipleVariables(c.c, &cError, cDomain, itemsList, valuesList, &cResults)
	if err := GetMmsError(cError); err != nil {
		return err
	}
	if accessResults != nil && cResults != nil {
		defer C.destroyMmsValueLinkedListLocal(cResults)
		*accessResults = (*accessResults)[:0]
		for node := cResults; node != nil; node = C.LinkedList_getNext(node) {
			data := C.LinkedList_getData(node)
			if data != nil {
				val := (*C.MmsValue)(data)
				*accessResults = append(*accessResults, MmsDataAccessError(C.MmsValue_getDataAccessError(val)))
			}
		}
	}
	return nil
}

// ReadNamedVariableListValues reads the values of a domain or VMD scoped named variable list.
// Pass domainID as "" for VMD scope. specification should be true for IEC 61850 compliant requests.
// Returns a single MmsValue of type Array (Value is []*MmsValue) or nil on empty/error. Caller does not free the result.
func (c *MmsConnection) ReadNamedVariableListValues(domainID, listName string, specification bool) (*MmsValue, error) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return nil, NotConnected
	}
	var cDomain *C.char
	var freeCDomain func()
	if domainID != "" {
		cDomain, freeCDomain = allocCString(domainID)
		defer freeCDomain()
	}
	cList, freecList := allocCString(listName)
	defer freecList()
	var cError C.MmsError
	result := C.MmsConnection_readNamedVariableListValues(c.c, &cError, cDomain, cList, C.bool(specification))
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	defer C.MmsValue_delete(result)
	return CMmsValueToMmsValue(result), nil
}

// ReadArrayElements reads one or more elements of an array variable. startIndex is the first element; numberOfElements is the count (0 = single element at startIndex).
// Returns the value (single element or array of elements) or nil on error. Caller does not free the result.
func (c *MmsConnection) ReadArrayElements(domainID, itemID string, startIndex, numberOfElements uint32) (*MmsValue, error) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return nil, NotConnected
	}
	cDomain, freecDomain := allocCString(domainID)
	defer freecDomain()
	cItem, freecItem := allocCString(itemID)
	defer freecItem()
	var cError C.MmsError
	result := C.MmsConnection_readArrayElements(c.c, &cError, cDomain, cItem, C.uint32_t(startIndex), C.uint32_t(numberOfElements))
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	defer C.MmsValue_delete(result)
	return CMmsValueToMmsValue(result), nil
}

// ReadNamedVariableListValuesAsync reads the values of a domain or VMD scoped named variable list asynchronously.
// Pass domainID as "" for VMD scope. specification should be true for IEC 61850 compliant requests.
// The callback receives a single MmsValue of type Array (Value is []*MmsValue) or nil on error. Caller does not own the value.
func (c *MmsConnection) ReadNamedVariableListValuesAsync(domainID, listName string, specification bool, callback func(*MmsValue, error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	var cDomain *C.char
	var freeCDomain func()
	if domainID != "" {
		cDomain, freeCDomain = allocCString(domainID)
		defer freeCDomain()
	}
	cList, freecList := allocCString(listName)
	defer freecList()
	ctx := &readVarAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_readNamedVariableListValuesAsync(conn, nil, &cError, cDomain, cList, C.bool(specification), (C.MmsConnection_ReadVariableHandler)(C.readVariableAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// WriteArrayElements writes one or more array elements. index is the first element; numberOfElements is the count (0 = single element). value is the data to write (MmsValueRef, not consumed).
func (c *MmsConnection) WriteArrayElements(domainID, itemID string, index, numberOfElements int, value *MmsValueRef) (MmsDataAccessError, error) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return 0, NotConnected
	}
	if value == nil || value.c == nil {
		return 0, UserProvidedInvalidArgument
	}
	cDomain, freecDomain := allocCString(domainID)
	defer freecDomain()
	cItem, freecItem := allocCString(itemID)
	defer freecItem()
	var cError C.MmsError
	accessErr := C.MmsConnection_writeArrayElements(c.c, &cError, cDomain, cItem, C.int(index), C.int(numberOfElements), value.c)
	if err := GetMmsError(cError); err != nil {
		return 0, err
	}
	return MmsDataAccessError(accessErr), nil
}

// WriteNamedVariableList writes values to a domain or VMD scoped named variable list.
// Pass domainID as "" for VMD scope. values must contain one MmsValueRef per list entry (not consumed).
// If accessResults is non-nil it is filled with the per-variable data access error results.
func (c *MmsConnection) WriteNamedVariableList(domainID, listName string, values []*MmsValueRef, accessResults *[]MmsDataAccessError) error {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return NotConnected
	}
	if len(values) == 0 {
		return UserProvidedInvalidArgument
	}
	var cDomain *C.char
	var freeCDomain func()
	if domainID != "" {
		cDomain, freeCDomain = allocCString(domainID)
		defer freeCDomain()
	}
	cList, freecList := allocCString(listName)
	defer freecList()
	clist := C.LinkedList_create()
	for _, v := range values {
		if v != nil && v.c != nil {
			C.LinkedList_add(clist, unsafe.Pointer(v.c))
		}
	}
	defer C.LinkedList_destroyStatic(clist)
	var cError C.MmsError
	var cResults C.LinkedList
	C.MmsConnection_writeNamedVariableList(c.c, &cError, C.bool(false), cDomain, cList, clist, &cResults)
	if err := GetMmsError(cError); err != nil {
		return err
	}
	if accessResults != nil && cResults != nil {
		defer C.destroyMmsValueLinkedListLocal(cResults)
		*accessResults = (*accessResults)[:0]
		for node := cResults; node != nil; node = C.LinkedList_getNext(node) {
			data := C.LinkedList_getData(node)
			if data != nil {
				val := (*C.MmsValue)(data)
				*accessResults = append(*accessResults, MmsDataAccessError(C.MmsValue_getDataAccessError(val)))
			}
		}
	}
	return nil
}

// GetNamedVariableListAttributes returns the attributes of a named variable list (deletable, variable specs). Pass domainID as "" for VMD scope.
func (c *MmsConnection) GetNamedVariableListAttributes(domainID, listName string) (*MmsNamedVariableListAttributes, error) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return nil, NotConnected
	}
	var cDomain *C.char
	var freeCDomain func()
	if domainID != "" {
		cDomain, freeCDomain = allocCString(domainID)
		defer freeCDomain()
	}
	cList, freecList := allocCString(listName)
	defer freecList()
	var cError C.MmsError
	var cDeletable C.bool
	list := C.MmsConnection_readNamedVariableListDirectory(c.c, &cError, cDomain, cList, &cDeletable)
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	out := &MmsNamedVariableListAttributes{IsDeletable: bool(cDeletable)}
	if list == nil {
		return out, nil
	}
	defer C.LinkedList_destroyDeep(list, (C.LinkedListValueDeleteFunction)(C.MmsVariableAccessSpecification_destroy))
	for node := list; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data == nil {
			continue
		}
		spec := (*C.MmsVariableAccessSpecification)(data)
		out.Variables = append(out.Variables, MmsVariableAccessSpec{
			DomainID: C.GoString(spec.domainId),
			ItemID:   C.GoString(spec.itemId),
		})
	}
	return out, nil
}

// GetNamedVariableListAttributesAsync retrieves the attributes of a named variable list asynchronously.
// Pass domainID as "" for VMD scope. The callback may run from another goroutine.
func (c *MmsConnection) GetNamedVariableListAttributesAsync(domainID, listName string, callback func(*MmsNamedVariableListAttributes, error)) error {
	return c.ReadNamedVariableListDirectoryAsync(domainID, listName, func(specs []MmsVariableAccessSpec, deletable bool, err error) {
		if callback == nil {
			return
		}
		if err != nil {
			callback(nil, err)
			return
		}
		callback(&MmsNamedVariableListAttributes{IsDeletable: deletable, Variables: specs}, nil)
	})
}

// GetDomainVariableListNames returns the names of named variable lists in the given domain. Pass domainID as "" for VMD scope.
func (c *MmsConnection) GetDomainVariableListNames(domainID string) ([]string, error) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return nil, NotConnected
	}
	var cDomain *C.char
	var freeCDomain func()
	if domainID != "" {
		cDomain, freeCDomain = allocCString(domainID)
		defer freeCDomain()
	}
	var cError C.MmsError
	list := C.MmsConnection_getDomainVariableListNames(c.c, &cError, cDomain)
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}
	defer C.destroyCharPtrLinkedList(list)
	var names []string
	for node := list; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data != nil {
			names = append(names, C.GoString((*C.char)(data)))
		}
	}
	return names, nil
}

// GetDomainJournals returns the journal names in the given domain. Pass domainID as "" for VMD scope.
func (c *MmsConnection) GetDomainJournals(domainID string) ([]string, error) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return nil, NotConnected
	}
	var cDomain *C.char
	var freeCDomain func()
	if domainID != "" {
		cDomain, freeCDomain = allocCString(domainID)
		defer freeCDomain()
	}
	var cError C.MmsError
	list := C.MmsConnection_getDomainJournals(c.c, &cError, cDomain)
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}
	defer C.destroyCharPtrLinkedList(list)
	var names []string
	for node := list; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data != nil {
			names = append(names, C.GoString((*C.char)(data)))
		}
	}
	return names, nil
}

// GetVMDVariableNames returns the VMD-scope variable names. Pass "" for continueAfter in async version.
func (c *MmsConnection) GetVMDVariableNames() ([]string, error) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return nil, NotConnected
	}
	var cError C.MmsError
	list := C.MmsConnection_getVMDVariableNames(c.c, &cError)
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}
	defer C.destroyCharPtrLinkedList(list)
	var names []string
	for node := list; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data != nil {
			names = append(names, C.GoString((*C.char)(data)))
		}
	}
	return names, nil
}

// GetVMDVariableNamesAsync returns VMD variable names asynchronously. continueAfter is "" to start.
func (c *MmsConnection) GetVMDVariableNamesAsync(continueAfter string, callback func(names []string, moreFollows bool, err error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, false, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	if callback == nil {
		return UserProvidedInvalidArgument
	}
	var cCont *C.char
	var freeCCont func()
	if continueAfter != "" {
		cCont, freeCCont = allocCString(continueAfter)
		defer freeCCont()
	}
	ctx := &getNameListAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_getVMDVariableNamesAsync(conn, nil, &cError, cCont, (C.MmsConnection_GetNameListHandler)(C.getNameListAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// GetDomainNamesAsync returns domain names asynchronously. continueAfter is "" to start.
func (c *MmsConnection) GetDomainNamesAsync(continueAfter string, callback func(names []string, moreFollows bool, err error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, false, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	if callback == nil {
		return UserProvidedInvalidArgument
	}
	var cCont *C.char
	var freeCCont func()
	if continueAfter != "" {
		cCont, freeCCont = allocCString(continueAfter)
		defer freeCCont()
	}
	ctx := &getNameListAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_getDomainNamesAsync(conn, nil, &cError, cCont, nil, (C.MmsConnection_GetNameListHandler)(C.getNameListAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// GetDomainVariableNamesAsync returns variable names in a domain asynchronously. continueAfter is "" to start.
func (c *MmsConnection) GetDomainVariableNamesAsync(domainID, continueAfter string, callback func(names []string, moreFollows bool, err error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, false, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	if callback == nil {
		return UserProvidedInvalidArgument
	}
	var cDomain *C.char
	var freeCDomain func()
	if domainID != "" {
		cDomain, freeCDomain = allocCString(domainID)
		defer freeCDomain()
	}
	var cCont *C.char
	var freeCCont func()
	if continueAfter != "" {
		cCont, freeCCont = allocCString(continueAfter)
		defer freeCCont()
	}
	ctx := &getNameListAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_getDomainVariableNamesAsync(conn, nil, &cError, cDomain, cCont, nil, (C.MmsConnection_GetNameListHandler)(C.getNameListAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// GetDomainVariableListNamesAsync returns named variable list names in a domain asynchronously. continueAfter is "" to start.
func (c *MmsConnection) GetDomainVariableListNamesAsync(domainID, continueAfter string, callback func(names []string, moreFollows bool, err error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, false, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	if callback == nil {
		return UserProvidedInvalidArgument
	}
	var cDomain *C.char
	var freeCDomain func()
	if domainID != "" {
		cDomain, freeCDomain = allocCString(domainID)
		defer freeCDomain()
	}
	var cCont *C.char
	var freeCCont func()
	if continueAfter != "" {
		cCont, freeCCont = allocCString(continueAfter)
		defer freeCCont()
	}
	ctx := &getNameListAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_getDomainVariableListNamesAsync(conn, nil, &cError, cDomain, cCont, nil, (C.MmsConnection_GetNameListHandler)(C.getNameListAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// GetDomainJournalsAsync returns journal names in a domain asynchronously. continueAfter is "" to start.
func (c *MmsConnection) GetDomainJournalsAsync(domainID, continueAfter string, callback func(names []string, moreFollows bool, err error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, false, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	if callback == nil {
		return UserProvidedInvalidArgument
	}
	var cDomain *C.char
	var freeCDomain func()
	if domainID != "" {
		cDomain, freeCDomain = allocCString(domainID)
		defer freeCDomain()
	}
	var cCont *C.char
	var freeCCont func()
	if continueAfter != "" {
		cCont, freeCCont = allocCString(continueAfter)
		defer freeCCont()
	}
	ctx := &getNameListAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_getDomainJournalsAsync(conn, nil, &cError, cDomain, cCont, (C.MmsConnection_GetNameListHandler)(C.getNameListAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// GetVariableListNamesAssociationSpecificAsync returns association-specific variable list names asynchronously. continueAfter is "" to start.
func (c *MmsConnection) GetVariableListNamesAssociationSpecificAsync(continueAfter string, callback func(names []string, moreFollows bool, err error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, false, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	if callback == nil {
		return UserProvidedInvalidArgument
	}
	var cCont *C.char
	var freeCCont func()
	if continueAfter != "" {
		cCont, freeCCont = allocCString(continueAfter)
		defer freeCCont()
	}
	ctx := &getNameListAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_getVariableListNamesAssociationSpecificAsync(conn, nil, &cError, cCont, (C.MmsConnection_GetNameListHandler)(C.getNameListAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// GetServerStatus returns the MMS server status (VMD logical and physical status).
func (c *MmsConnection) GetServerStatus(extendedDerivation bool) (*MmsServerStatus, error) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return nil, NotConnected
	}
	var cError C.MmsError
	var vmdLogical, vmdPhysical C.int
	C.MmsConnection_getServerStatus(c.c, &cError, &vmdLogical, &vmdPhysical, C.bool(extendedDerivation))
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	return &MmsServerStatus{
		VmdLogicalStatus:  int32(vmdLogical),
		VmdPhysicalStatus: int32(vmdPhysical),
		LocalDetail:       0,
	}, nil
}

// ObtainFile requests the server to read a file from the client (upload: sourceFile local, destinationFile remote).
func (c *MmsConnection) ObtainFile(sourceFile, destinationFile string) error {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return NotConnected
	}
	cSrc, freecSrc := allocCString(sourceFile)
	defer freecSrc()
	cDst, freecDst := allocCString(destinationFile)
	defer freecDst()
	var cError C.MmsError
	C.MmsConnection_obtainFile(c.c, &cError, cSrc, cDst)
	return GetMmsError(cError)
}

// RenameFile renames a file on the server.
func (c *MmsConnection) RenameFile(currentName, newName string) error {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return NotConnected
	}
	cCur, freecCur := allocCString(currentName)
	defer freecCur()
	cNew, freecNew := allocCString(newName)
	defer freecNew()
	var cError C.MmsError
	C.MmsConnection_fileRename(c.c, &cError, cCur, cNew)
	return GetMmsError(cError)
}

// SendRawData sends raw data on the connection (for test purposes). buffer is the payload to send.
func (c *MmsConnection) SendRawData(buffer []byte) error {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return NotConnected
	}
	var cError C.MmsError
	var cBuf *C.uint8_t
	cSize := C.int(0)
	if len(buffer) > 0 {
		cBuf = (*C.uint8_t)(unsafe.Pointer(&buffer[0]))
		cSize = C.int(len(buffer))
	}
	C.MmsConnection_sendRawData(c.c, &cError, cBuf, cSize)
	return GetMmsError(cError)
}

// FileDirectoryAsync retrieves the file directory asynchronously. fileSpecification is the path/pattern; use "" for continueAfter to start from the beginning.
// The callback receives the list of entries, moreFollows (true if more data is available server-side), and any error.
func (c *MmsConnection) FileDirectoryAsync(fileSpecification, continueAfter string, callback func(entries []MmsFileDirectoryEntryEx, moreFollows bool, err error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, false, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	if callback == nil {
		return UserProvidedInvalidArgument
	}
	cSpec, freecSpec := allocCString(fileSpecification)
	defer freecSpec()
	var cCont *C.char
	var freeCCont func()
	if continueAfter != "" {
		cCont, freeCCont = allocCString(continueAfter)
		defer freeCCont()
	}
	ctx := &fileDirAsyncCtx{callback: callback}
	fileDirAsyncRegistryMu.Lock()
	fileDirAsyncRegistry[conn] = ctx
	fileDirAsyncRegistryMu.Unlock()
	var cError C.MmsError
	C.MmsConnection_getFileDirectoryAsync(conn, nil, &cError, cSpec, cCont, (C.MmsConnection_FileDirectoryHandler)(C.fileDirectoryAsyncBridge), unsafe.Pointer(conn))
	if err := GetMmsError(cError); err != nil {
		fileDirAsyncRegistryMu.Lock()
		delete(fileDirAsyncRegistry, conn)
		fileDirAsyncRegistryMu.Unlock()
		callback(nil, false, err)
		return err
	}
	return nil
}

func convertCJournalEntryToMms(entry C.MmsJournalEntry) MmsJournalEntry {
	je := MmsJournalEntry{}
	if eid := C.MmsJournalEntry_getEntryID(entry); eid != nil {
		if C.MmsValue_getType(eid) == C.MMS_OCTET_STRING {
			n := int(C.MmsValue_getOctetStringSize(eid))
			je.EntryID = make([]byte, n)
			for i := 0; i < n; i++ {
				je.EntryID[i] = byte(C.MmsValue_getOctetStringOctet(eid, C.int(i)))
			}
		}
	}
	if ot := C.MmsJournalEntry_getOccurenceTime(entry); ot != nil {
		je.OccurTime = uint64(C.MmsValue_getBinaryTimeAsUtcMs(ot))
	}
	varsList := C.MmsJournalEntry_getJournalVariables(entry)
	var parts []*MmsValue
	for node := varsList; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data == nil {
			continue
		}
		jv := (C.MmsJournalVariable)(data)
		val := C.MmsJournalVariable_getValue(jv)
		if val != nil {
			parts = append(parts, CMmsValueToMmsValue(val))
		}
	}
	if len(parts) > 0 {
		je.EntryContent = &MmsValue{Type: Structure, Value: parts}
	}
	return je
}

// ReadJournal reads journal entries with optional time range. If startingTime and endingTime are nil, 0 and max are used.
func (c *MmsConnection) ReadJournal(domainID, journalName string, startingTime, endingTime *uint64) ([]*MmsJournalEntry, bool, error) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return nil, false, NotConnected
	}
	cDomain, freecDomain := allocCString(domainID)
	defer freecDomain()
	cJournal, freecJournal := allocCString(journalName)
	defer freecJournal()
	var startV, endV *C.MmsValue
	if startingTime != nil {
		startV = C.MmsValue_newBinaryTime(C.bool(false))
		defer C.MmsValue_delete(startV)
		C.MmsValue_setBinaryTime(startV, C.uint64_t(*startingTime))
	}
	if endingTime != nil {
		endV = C.MmsValue_newBinaryTime(C.bool(false))
		defer C.MmsValue_delete(endV)
		C.MmsValue_setBinaryTime(endV, C.uint64_t(*endingTime))
	}
	if startV == nil {
		startV = C.MmsValue_newBinaryTime(C.bool(false))
		defer C.MmsValue_delete(startV)
		C.MmsValue_setBinaryTime(startV, 0)
	}
	if endV == nil {
		endV = C.MmsValue_newBinaryTime(C.bool(false))
		defer C.MmsValue_delete(endV)
		C.MmsValue_setBinaryTime(endV, 0xffffffffffff)
	}
	var cMore C.bool
	var cError C.MmsError
	list := C.MmsConnection_readJournalTimeRange(c.c, &cError, cDomain, cJournal, startV, endV, &cMore)
	if err := GetMmsError(cError); err != nil {
		return nil, false, err
	}
	if list == nil {
		return nil, bool(cMore), nil
	}
		defer C.destroyJournalEntryLinkedListLocal(list)
	var entries []*MmsJournalEntry
	for node := list; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data != nil {
			e := convertCJournalEntryToMms(C.MmsJournalEntry(data))
			entries = append(entries, &e)
		}
	}
	return entries, bool(cMore), nil
}

// ReadJournalTimeRangeAsync reads journal entries in the given time range asynchronously (milliseconds since Unix epoch).
// The callback may run from another goroutine.
func (c *MmsConnection) ReadJournalTimeRangeAsync(domainID, journalName string, startTime, endTime uint64, callback func(entries []*MmsJournalEntry, moreFollows bool, err error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, false, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	if callback == nil {
		return UserProvidedInvalidArgument
	}
	cDomain, freecDomain := allocCString(domainID)
	defer freecDomain()
	cJournal, freecJournal := allocCString(journalName)
	defer freecJournal()
	startV := C.MmsValue_newBinaryTime(C.bool(false))
	defer C.MmsValue_delete(startV)
	C.MmsValue_setBinaryTime(startV, C.uint64_t(startTime))
	endV := C.MmsValue_newBinaryTime(C.bool(false))
	defer C.MmsValue_delete(endV)
	C.MmsValue_setBinaryTime(endV, C.uint64_t(endTime))
	ctx := &readJournalAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_readJournalTimeRangeAsync(conn, nil, &cError, cDomain, cJournal, startV, endV, (C.MmsConnection_ReadJournalHandler)(C.readJournalAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// ReadJournalStartAfterAsync reads journal entries starting after the given time and entry specification asynchronously.
// The callback may run from another goroutine.
func (c *MmsConnection) ReadJournalStartAfterAsync(domainID, journalName string, entryID []byte, timeSpec *uint64, callback func(entries []*MmsJournalEntry, moreFollows bool, err error)) error {
	c.connMu.Lock()
	if c.c == nil {
		c.connMu.Unlock()
		if callback != nil {
			callback(nil, false, NotConnected)
		}
		return NotConnected
	}
	conn := c.c
	c.connMu.Unlock()
	if callback == nil {
		return UserProvidedInvalidArgument
	}
	cDomain, freecDomain := allocCString(domainID)
	defer freecDomain()
	cJournal, freecJournal := allocCString(journalName)
	defer freecJournal()
	timeV := C.MmsValue_newBinaryTime(C.bool(false))
	defer C.MmsValue_delete(timeV)
	if timeSpec != nil {
		C.MmsValue_setBinaryTime(timeV, C.uint64_t(*timeSpec))
	}
	var entryV *C.MmsValue
	if len(entryID) > 0 {
		entryV = C.MmsValue_newOctetString(C.int(len(entryID)), C.int(len(entryID)))
		defer C.MmsValue_delete(entryV)
		for i, b := range entryID {
			C.MmsValue_setOctetStringOctet(entryV, C.int(i), C.uint8_t(b))
		}
	}
	ctx := &readJournalAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_readJournalStartAfterAsync(conn, nil, &cError, cDomain, cJournal, timeV, entryV, (C.MmsConnection_ReadJournalHandler)(C.readJournalAsyncBridge), unsafe.Pointer(ctx))
	return GetMmsError(cError)
}

// ReadJournalTimeRange reads journal entries in the given time range (milliseconds since Unix epoch).
func (c *MmsConnection) ReadJournalTimeRange(domainID, journalName string, startTime, endTime uint64) ([]*MmsJournalEntry, bool, error) {
	return c.ReadJournal(domainID, journalName, &startTime, &endTime)
}

// ReadJournalStartAfter reads journal entries starting after the given entry and optional time.
func (c *MmsConnection) ReadJournalStartAfter(domainID, journalName string, entryID []byte, timeSpec *uint64) ([]*MmsJournalEntry, bool, error) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return nil, false, NotConnected
	}
	cDomain, freecDomain := allocCString(domainID)
	defer freecDomain()
	cJournal, freecJournal := allocCString(journalName)
	defer freecJournal()
	timeV := C.MmsValue_newBinaryTime(C.bool(false))
	defer C.MmsValue_delete(timeV)
	if timeSpec != nil {
		C.MmsValue_setBinaryTime(timeV, C.uint64_t(*timeSpec))
	}
	var entryV *C.MmsValue
	if len(entryID) > 0 {
		entryV = C.MmsValue_newOctetString(C.int(len(entryID)), C.int(len(entryID)))
		defer C.MmsValue_delete(entryV)
		for i, b := range entryID {
			C.MmsValue_setOctetStringOctet(entryV, C.int(i), C.uint8_t(b))
		}
	}
	var cMore C.bool
	var cError C.MmsError
	list := C.MmsConnection_readJournalStartAfter(c.c, &cError, cDomain, cJournal, timeV, entryV, &cMore)
	if err := GetMmsError(cError); err != nil {
		return nil, false, err
	}
	if list == nil {
		return nil, bool(cMore), nil
	}
	defer C.destroyJournalEntryLinkedListLocal(list)
	var entries []*MmsJournalEntry
	for node := list; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data != nil {
			e := convertCJournalEntryToMms(C.MmsJournalEntry(data))
			entries = append(entries, &e)
		}
	}
	return entries, bool(cMore), nil
}
