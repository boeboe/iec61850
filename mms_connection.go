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
		C.MmsConnection_destroy(conn)
	}
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

	host := C.CString(hostname)
	defer C.free(unsafe.Pointer(host))
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
	cDomain := C.CString(domainID)
	defer C.free(unsafe.Pointer(cDomain))
	cItem := C.CString(itemID)
	defer C.free(unsafe.Pointer(cItem))
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
	cDomain := C.CString(domainID)
	defer C.free(unsafe.Pointer(cDomain))
	cItem := C.CString(itemID)
	defer C.free(unsafe.Pointer(cItem))
	ctx := &writeVarAsyncCtx{callback: callback}
	var cError C.MmsError
	C.MmsConnection_writeVariableAsync(conn, nil, &cError, cDomain, cItem, value.c, (C.MmsConnection_WriteVariableHandler)(C.writeVariableAsyncBridge), unsafe.Pointer(ctx))
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
	cDomain := C.CString(domainID)
	defer C.free(unsafe.Pointer(cDomain))
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
	if domainID != "" {
		cDomain = C.CString(domainID)
		defer C.free(unsafe.Pointer(cDomain))
	}
	cList := C.CString(listName)
	defer C.free(unsafe.Pointer(cList))
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
	if domainID != "" {
		cDomain = C.CString(domainID)
		defer C.free(unsafe.Pointer(cDomain))
	}
	cList := C.CString(listName)
	defer C.free(unsafe.Pointer(cList))
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
	if domainID != "" {
		cDomain = C.CString(domainID)
		defer C.free(unsafe.Pointer(cDomain))
	}
	cList := C.CString(listName)
	defer C.free(unsafe.Pointer(cList))
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

// GetDomainVariableListNames returns the names of named variable lists in the given domain. Pass domainID as "" for VMD scope.
func (c *MmsConnection) GetDomainVariableListNames(domainID string) ([]string, error) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	if c.c == nil {
		return nil, NotConnected
	}
	var cDomain *C.char
	if domainID != "" {
		cDomain = C.CString(domainID)
		defer C.free(unsafe.Pointer(cDomain))
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
	if domainID != "" {
		cDomain = C.CString(domainID)
		defer C.free(unsafe.Pointer(cDomain))
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
	cSrc := C.CString(sourceFile)
	defer C.free(unsafe.Pointer(cSrc))
	cDst := C.CString(destinationFile)
	defer C.free(unsafe.Pointer(cDst))
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
	cCur := C.CString(currentName)
	defer C.free(unsafe.Pointer(cCur))
	cNew := C.CString(newName)
	defer C.free(unsafe.Pointer(cNew))
	var cError C.MmsError
	C.MmsConnection_fileRename(c.c, &cError, cCur, cNew)
	return GetMmsError(cError)
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
	cDomain := C.CString(domainID)
	defer C.free(unsafe.Pointer(cDomain))
	cJournal := C.CString(journalName)
	defer C.free(unsafe.Pointer(cJournal))
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
	cDomain := C.CString(domainID)
	defer C.free(unsafe.Pointer(cDomain))
	cJournal := C.CString(journalName)
	defer C.free(unsafe.Pointer(cJournal))
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
