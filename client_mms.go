package iec61850

/*
#include <iec61850_client.h>
#include <mms_client_connection.h>
#include <mms_value.h>
#include <stdlib.h>

static void destroyMmsValueLinkedList(LinkedList L) {
	if (L) LinkedList_destroyDeep(L, (LinkedListValueDeleteFunction)MmsValue_delete);
}

static void destroyJournalEntryLinkedList(LinkedList L) {
	if (L) LinkedList_destroyDeep(L, (LinkedListValueDeleteFunction)MmsJournalEntry_destroy);
}
*/
import "C"
import (
	"unsafe"
)

// VariableAccessSpec describes a single variable for defining a named variable list.
// Use ArrayIndex -1 and ComponentName "" for simple variables; set them for array/component access.
type VariableAccessSpec struct {
	DomainID      string
	ItemID        string
	ArrayIndex    int32  // -1 for no array index
	ComponentName string // optional component of array element
}

// VariableListEntry is one entry from ReadNamedVariableListDirectory.
type VariableListEntry struct {
	DomainID      string
	ItemID        string
	ArrayIndex    int32
	ComponentName string
}

// getMmsConnection returns the underlying MmsConnection for the client. Must be connected.
func (c *Client) getMmsConnection() C.MmsConnection {
	return C.IedConnection_getMmsConnection(c.conn)
}

// GetConnectionParameters returns the MMS connection parameters (max outstanding calls, PDU size, etc.).
// The client must have an established connection (conn may be set; parameters are updated after connect).
func (c *Client) GetConnectionParameters() (*MmsConnectionParameters, error) {
	if c.conn == nil {
		return nil, NotConnected
	}
	p := C.MmsConnection_getMmsConnectionParameters(c.getMmsConnection())
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
	}, nil
}

// GetServerStatus returns the MMS server status (VMD logical and physical status).
// extendedDerivation instructs the server to run self-diagnosis to determine status.
func (c *Client) GetServerStatus(extendedDerivation bool) (*MmsServerStatus, error) {
	var cError C.MmsError
	var vmdLogical, vmdPhysical C.int
	C.MmsConnection_getServerStatus(c.getMmsConnection(), &cError, &vmdLogical, &vmdPhysical, C.bool(extendedDerivation))
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	return &MmsServerStatus{
		VmdLogicalStatus:  int32(vmdLogical),
		VmdPhysicalStatus: int32(vmdPhysical),
		LocalDetail:       0,
	}, nil
}

// ObtainFile sends an obtainFile request: the server will read the file from the client.
// sourceFile is the local (client) path, destinationFile is the remote (server) path.
// This is the same as uploading a file to the server (MMS obtainFile service).
func (c *Client) ObtainFile(sourceFile, destinationFile string) error {
	cSrc := C.CString(sourceFile)
	defer C.free(unsafe.Pointer(cSrc))
	cDst := C.CString(destinationFile)
	defer C.free(unsafe.Pointer(cDst))
	var cError C.MmsError
	C.MmsConnection_obtainFile(c.getMmsConnection(), &cError, cSrc, cDst)
	return GetMmsError(cError)
}

// RenameFile renames a file on the server (currentName -> newName).
func (c *Client) RenameFile(currentName, newName string) error {
	cCur := C.CString(currentName)
	defer C.free(unsafe.Pointer(cCur))
	cNew := C.CString(newName)
	defer C.free(unsafe.Pointer(cNew))
	var cError C.MmsError
	C.MmsConnection_fileRename(c.getMmsConnection(), &cError, cCur, cNew)
	return GetMmsError(cError)
}

// GetVariableAccessAttributes returns the MMS variable access attributes (type specification) for a named variable.
// The caller must call Free() on the returned ref when done.
func (c *Client) GetVariableAccessAttributes(domainID, itemID string) (*MmsVariableSpecificationRef, error) {
	cDomain := C.CString(domainID)
	defer C.free(unsafe.Pointer(cDomain))
	cItem := C.CString(itemID)
	defer C.free(unsafe.Pointer(cItem))
	var cError C.MmsError
	spec := C.MmsConnection_getVariableAccessAttributes(c.getMmsConnection(), &cError, cDomain, cItem)
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if spec == nil {
		return nil, NullPointer
	}
	return &MmsVariableSpecificationRef{c: spec, owned: true, libraryOwned: true}, nil
}

// GetDomainNames returns the list of MMS domain names on the server.
func (c *Client) GetDomainNames() ([]string, error) {
	var cError C.MmsError
	list := C.MmsConnection_getDomainNames(c.getMmsConnection(), &cError)
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}
	defer C.LinkedList_destroy(list)
	return linkedListToStrings(list), nil
}

// GetDomainVariableNames returns the names of variables in an MMS domain.
func (c *Client) GetDomainVariableNames(domainID string) ([]string, error) {
	cDomain := C.CString(domainID)
	defer C.free(unsafe.Pointer(cDomain))
	var cError C.MmsError
	list := C.MmsConnection_getDomainVariableNames(c.getMmsConnection(), &cError, cDomain)
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}
	defer C.LinkedList_destroy(list)
	return linkedListToStrings(list), nil
}

// GetDomainVariableListNames returns the names of named variable lists in a domain, or VMD scope if domainID is "".
func (c *Client) GetDomainVariableListNames(domainID string) ([]string, error) {
	var cDomain *C.char
	if domainID != "" {
		cDomain = C.CString(domainID)
		defer C.free(unsafe.Pointer(cDomain))
	}
	var cError C.MmsError
	list := C.MmsConnection_getDomainVariableListNames(c.getMmsConnection(), &cError, cDomain)
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}
	defer C.LinkedList_destroy(list)
	return linkedListToStrings(list), nil
}

// GetVariableListNamesAssociationSpecific returns the names of association-specific named variable lists.
func (c *Client) GetVariableListNamesAssociationSpecific() ([]string, error) {
	var cError C.MmsError
	list := C.MmsConnection_getVariableListNamesAssociationSpecific(c.getMmsConnection(), &cError)
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}
	defer C.LinkedList_destroy(list)
	return linkedListToStrings(list), nil
}

func linkedListToStrings(list C.LinkedList) []string {
	var out []string
	for node := list; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data != nil {
			out = append(out, C.GoString((*C.char)(data)))
		}
	}
	return out
}

// WriteMultipleVariablesFromSpecs writes multiple variables in one request using variable access specs.
// All items must be in the same domain (domainID). values must have the same length as items; each value is written to the variable specified by the corresponding VariableAccessSpec. Returns one MmsDataAccessError per variable.
func (c *Client) WriteMultipleVariablesFromSpecs(domainID string, items []VariableAccessSpec, values []*MmsValueRef) ([]MmsDataAccessError, error) {
	if len(items) != len(values) {
		return nil, UserProvidedInvalidArgument
	}
	if len(items) == 0 {
		return nil, nil
	}
	itemIDs := make([]string, len(items))
	for i := range items {
		itemIDs[i] = items[i].ItemID
	}
	return c.WriteMultipleVariables(domainID, itemIDs, values)
}

// WriteMultipleVariables writes multiple variables in one request. itemIDs and values must have the same length; each value is written to the variable named by the corresponding itemID in the given domain. Returns one MmsDataAccessError per variable.
func (c *Client) WriteMultipleVariables(domainID string, itemIDs []string, values []*MmsValueRef) ([]MmsDataAccessError, error) {
	if len(itemIDs) != len(values) {
		return nil, UserProvidedInvalidArgument
	}
	if len(itemIDs) == 0 {
		return nil, nil
	}
	cDomain := C.CString(domainID)
	defer C.free(unsafe.Pointer(cDomain))
	itemsList := C.LinkedList_create()
	defer C.LinkedList_destroyDeep(itemsList, (C.LinkedListValueDeleteFunction)(C.free))
	valuesList := C.LinkedList_create()
	for _, id := range itemIDs {
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
	C.MmsConnection_writeMultipleVariables(c.getMmsConnection(), &cError, cDomain, itemsList, valuesList, &cResults)
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if cResults == nil {
		return nil, nil
	}
	defer C.destroyMmsValueLinkedList(cResults)
	var results []MmsDataAccessError
	for node := cResults; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data != nil {
			val := (*C.MmsValue)(data)
			results = append(results, MmsDataAccessError(C.MmsValue_getDataAccessError(val)))
		}
	}
	return results, nil
}

// ReadNamedVariableListValues reads all values from a domain or VMD scoped named variable list.
// Pass domainID as "" for VMD scope. specWithResult should be true for IEC 61850 compliant requests.
// The returned slice contains one MmsValue per list entry; the caller does not own the underlying C memory (it is freed by the library after the call returns), so values are converted to Go types.
func (c *Client) ReadNamedVariableListValues(domainID, listName string, specWithResult bool) ([]*MmsValue, error) {
	var cDomain *C.char
	if domainID != "" {
		cDomain = C.CString(domainID)
		defer C.free(unsafe.Pointer(cDomain))
	}
	cList := C.CString(listName)
	defer C.free(unsafe.Pointer(cList))
	var cError C.MmsError
	result := C.MmsConnection_readNamedVariableListValues(c.getMmsConnection(), &cError, cDomain, cList, C.bool(specWithResult))
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	defer C.MmsValue_delete(result)
	return toGoStructure(result, Array)
}

// ReadNamedVariableListValuesAssociationSpecific reads values from an association-specific named variable list.
func (c *Client) ReadNamedVariableListValuesAssociationSpecific(listName string, specWithResult bool) ([]*MmsValue, error) {
	cList := C.CString(listName)
	defer C.free(unsafe.Pointer(cList))
	var cError C.MmsError
	result := C.MmsConnection_readNamedVariableListValuesAssociationSpecific(c.getMmsConnection(), &cError, cList, C.bool(specWithResult))
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	defer C.MmsValue_delete(result)
	return toGoStructure(result, Array)
}

// ReadNamedVariableListDirectory returns the directory (list of variable references) and whether the list is deletable.
func (c *Client) ReadNamedVariableListDirectory(domainID, listName string) (entries []VariableListEntry, deletable bool, err error) {
	var cDomain *C.char
	if domainID != "" {
		cDomain = C.CString(domainID)
		defer C.free(unsafe.Pointer(cDomain))
	}
	cList := C.CString(listName)
	defer C.free(unsafe.Pointer(cList))
	var cError C.MmsError
	var cDeletable C.bool
	list := C.MmsConnection_readNamedVariableListDirectory(c.getMmsConnection(), &cError, cDomain, cList, &cDeletable)
	if err = GetMmsError(cError); err != nil {
		return nil, false, err
	}
	deletable = bool(cDeletable)
	if list == nil {
		return nil, deletable, nil
	}
	defer C.LinkedList_destroyDeep(list, (C.LinkedListValueDeleteFunction)(C.MmsVariableAccessSpecification_destroy))
	for node := list; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data == nil {
			continue
		}
		spec := (*C.MmsVariableAccessSpecification)(data)
		entries = append(entries, VariableListEntry{
			DomainID:      C.GoString(spec.domainId),
			ItemID:        C.GoString(spec.itemId),
			ArrayIndex:    int32(spec.arrayIndex),
			ComponentName: C.GoString(spec.componentName),
		})
	}
	return entries, deletable, nil
}

// ReadNamedVariableListDirectoryAssociationSpecific returns the directory for an association-specific list.
func (c *Client) ReadNamedVariableListDirectoryAssociationSpecific(listName string) (entries []VariableListEntry, deletable bool, err error) {
	cList := C.CString(listName)
	defer C.free(unsafe.Pointer(cList))
	var cError C.MmsError
	var cDeletable C.bool
	list := C.MmsConnection_readNamedVariableListDirectoryAssociationSpecific(c.getMmsConnection(), &cError, cList, &cDeletable)
	if err = GetMmsError(cError); err != nil {
		return nil, false, err
	}
	deletable = bool(cDeletable)
	if list == nil {
		return nil, deletable, nil
	}
	defer C.LinkedList_destroyDeep(list, (C.LinkedListValueDeleteFunction)(C.MmsVariableAccessSpecification_destroy))
	for node := list; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data == nil {
			continue
		}
		spec := (*C.MmsVariableAccessSpecification)(data)
		entries = append(entries, VariableListEntry{
			DomainID:      C.GoString(spec.domainId),
			ItemID:        C.GoString(spec.itemId),
			ArrayIndex:    int32(spec.arrayIndex),
			ComponentName: C.GoString(spec.componentName),
		})
	}
	return entries, deletable, nil
}

// DefineNamedVariableList defines a new domain or VMD scoped named variable list. Pass domainID as "" for VMD scope.
func (c *Client) DefineNamedVariableList(domainID, listName string, variableSpecs []VariableAccessSpec) error {
	if len(variableSpecs) == 0 {
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
	defer C.LinkedList_destroyDeep(clist, (C.LinkedListValueDeleteFunction)(C.MmsVariableAccessSpecification_destroy))
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
	var cError C.MmsError
	C.MmsConnection_defineNamedVariableList(c.getMmsConnection(), &cError, cDomain, cList, clist)
	return GetMmsError(cError)
}

// DefineNamedVariableListAssociationSpecific defines a new association-specific named variable list.
func (c *Client) DefineNamedVariableListAssociationSpecific(listName string, variableSpecs []VariableAccessSpec) error {
	if len(variableSpecs) == 0 {
		return UserProvidedInvalidArgument
	}
	cList := C.CString(listName)
	defer C.free(unsafe.Pointer(cList))
	clist := C.LinkedList_create()
	defer C.LinkedList_destroyDeep(clist, (C.LinkedListValueDeleteFunction)(C.MmsVariableAccessSpecification_destroy))
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
	var cError C.MmsError
	C.MmsConnection_defineNamedVariableListAssociationSpecific(c.getMmsConnection(), &cError, cList, clist)
	return GetMmsError(cError)
}

// DeleteNamedVariableList deletes a domain or VMD scoped named variable list. Pass domainID as "" for VMD scope.
func (c *Client) DeleteNamedVariableList(domainID, listName string) (bool, error) {
	var cDomain *C.char
	if domainID != "" {
		cDomain = C.CString(domainID)
		defer C.free(unsafe.Pointer(cDomain))
	}
	cList := C.CString(listName)
	defer C.free(unsafe.Pointer(cList))
	var cError C.MmsError
	ok := C.MmsConnection_deleteNamedVariableList(c.getMmsConnection(), &cError, cDomain, cList)
	if err := GetMmsError(cError); err != nil {
		return false, err
	}
	return bool(ok), nil
}

// DeleteAssociationSpecificNamedVariableList deletes an association-specific named variable list.
func (c *Client) DeleteAssociationSpecificNamedVariableList(listName string) (bool, error) {
	cList := C.CString(listName)
	defer C.free(unsafe.Pointer(cList))
	var cError C.MmsError
	ok := C.MmsConnection_deleteAssociationSpecificNamedVariableList(c.getMmsConnection(), &cError, cList)
	if err := GetMmsError(cError); err != nil {
		return false, err
	}
	return bool(ok), nil
}

// SetNamedVariableListValues is an alias for WriteNamedVariableList: writes values to a domain or VMD scoped named variable list.
func (c *Client) SetNamedVariableListValues(domainID, listName string, values []*MmsValueRef) ([]MmsDataAccessError, error) {
	return c.WriteNamedVariableList(domainID, listName, values)
}

// WriteNamedVariableList writes values to a domain or VMD scoped named variable list.
// Pass domainID as "" for VMD scope. values must contain one MmsValueRef per list entry; refs are not consumed.
// Returns the data access error result for each variable write.
func (c *Client) WriteNamedVariableList(domainID, listName string, values []*MmsValueRef) ([]MmsDataAccessError, error) {
	if len(values) == 0 {
		return nil, UserProvidedInvalidArgument
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
	C.MmsConnection_writeNamedVariableList(c.getMmsConnection(), &cError, C.bool(false), cDomain, cList, clist, &cResults)
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if cResults == nil {
		return nil, nil
	}
	defer C.destroyMmsValueLinkedList(cResults)
	var results []MmsDataAccessError
	for node := cResults; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data != nil {
			val := (*C.MmsValue)(data)
			results = append(results, MmsDataAccessError(C.MmsValue_getDataAccessError(val)))
		}
	}
	return results, nil
}

// MmsServerIdentity holds the result of MMS Identify (VMD identity).
type MmsServerIdentity struct {
	VendorName string
	ModelName  string
	Revision   string
}

// Identify returns the server identity (vendor name, model name, revision) via the MMS identify service.
func (c *Client) Identify() (*MmsServerIdentity, error) {
	var cError C.MmsError
	ident := C.MmsConnection_identify(c.getMmsConnection(), &cError)
	if err := GetMmsError(cError); err != nil {
		return nil, err
	}
	if ident == nil {
		return nil, NullPointer
	}
	defer C.MmsServerIdentity_destroy(ident)
	return &MmsServerIdentity{
		VendorName: C.GoString(ident.vendorName),
		ModelName:  C.GoString(ident.modelName),
		Revision:   C.GoString(ident.revision),
	}, nil
}

// JournalVariable is one variable in a journal entry.
type JournalVariable struct {
	Tag   string
	Value *MmsValue
}

// JournalEntry is one entry from a journal read.
type JournalEntry struct {
	EntryID        *MmsValue        // Octet string
	OccurrenceTime *MmsValue        // Binary time
	Variables      []JournalVariable
}

func convertJournalEntry(entry C.MmsJournalEntry) JournalEntry {
	je := JournalEntry{}
	je.EntryID = nil
	je.OccurrenceTime = nil
	if eid := C.MmsJournalEntry_getEntryID(entry); eid != nil {
		t := MmsType(C.MmsValue_getType(eid))
		if v, err := toGoValue(eid, t); err == nil {
			je.EntryID = &MmsValue{Type: t, Value: v}
		}
	}
	if ot := C.MmsJournalEntry_getOccurenceTime(entry); ot != nil {
		t := MmsType(C.MmsValue_getType(ot))
		if v, err := toGoValue(ot, t); err == nil {
			je.OccurrenceTime = &MmsValue{Type: t, Value: v}
		}
	}
	varsList := C.MmsJournalEntry_getJournalVariables(entry)
	for node := varsList; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data == nil {
			continue
		}
		jv := (C.MmsJournalVariable)(data)
		tag := ""
		if t := C.MmsJournalVariable_getTag(jv); t != nil {
			tag = C.GoString(t)
		}
		val := C.MmsJournalVariable_getValue(jv)
		var goVal *MmsValue
		if val != nil {
			mt := MmsType(C.MmsValue_getType(val))
			if v, err := toGoValue(val, mt); err == nil {
				goVal = &MmsValue{Type: mt, Value: v}
			}
		}
		je.Variables = append(je.Variables, JournalVariable{Tag: tag, Value: goVal})
	}
	return je
}

// ReadJournalTimeRange reads journal entries in the given time range. startTimeMs and endTimeMs are milliseconds since Unix epoch (binary time). The caller does not own the returned entries.
func (c *Client) ReadJournalTimeRange(domainID, itemID string, startTimeMs, endTimeMs uint64) (entries []JournalEntry, moreFollows bool, err error) {
	cDomain := C.CString(domainID)
	defer C.free(unsafe.Pointer(cDomain))
	cItem := C.CString(itemID)
	defer C.free(unsafe.Pointer(cItem))
	startV := C.MmsValue_newBinaryTime(C.bool(false))
	defer C.MmsValue_delete(startV)
	C.MmsValue_setBinaryTime(startV, C.uint64_t(startTimeMs))
	endV := C.MmsValue_newBinaryTime(C.bool(false))
	defer C.MmsValue_delete(endV)
	C.MmsValue_setBinaryTime(endV, C.uint64_t(endTimeMs))
	var cMore C.bool
	var cError C.MmsError
	list := C.MmsConnection_readJournalTimeRange(c.getMmsConnection(), &cError, cDomain, cItem, startV, endV, &cMore)
	if err = GetMmsError(cError); err != nil {
		return nil, false, err
	}
	moreFollows = bool(cMore)
	if list == nil {
		return nil, moreFollows, nil
	}
	defer C.destroyJournalEntryLinkedList(list)
	for node := list; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data != nil {
			entries = append(entries, convertJournalEntry(C.MmsJournalEntry(data)))
		}
	}
	return entries, moreFollows, nil
}

// ReadJournalStartAfter reads journal entries starting after the given time and entry specification.
// timeSpecification and entrySpecification are MMS binary time and octet string respectively; pass nil for entrySpecification to start from the time.
func (c *Client) ReadJournalStartAfter(domainID, itemID string, timeSpecificationMs uint64, entrySpecification []byte) (entries []JournalEntry, moreFollows bool, err error) {
	cDomain := C.CString(domainID)
	defer C.free(unsafe.Pointer(cDomain))
	cItem := C.CString(itemID)
	defer C.free(unsafe.Pointer(cItem))
	timeV := C.MmsValue_newBinaryTime(C.bool(false))
	defer C.MmsValue_delete(timeV)
	C.MmsValue_setBinaryTime(timeV, C.uint64_t(timeSpecificationMs))
	var entryV *C.MmsValue
	if len(entrySpecification) > 0 {
		entryV = C.MmsValue_newOctetString(C.int(len(entrySpecification)), C.int(len(entrySpecification)))
		defer C.MmsValue_delete(entryV)
		for i, b := range entrySpecification {
			C.MmsValue_setOctetStringOctet(entryV, C.int(i), C.uint8_t(b))
		}
	}
	var cMore C.bool
	var cError C.MmsError
	list := C.MmsConnection_readJournalStartAfter(c.getMmsConnection(), &cError, cDomain, cItem, timeV, entryV, &cMore)
	if err = GetMmsError(cError); err != nil {
		return nil, false, err
	}
	moreFollows = bool(cMore)
	if list == nil {
		return nil, moreFollows, nil
	}
	defer C.destroyJournalEntryLinkedList(list)
	for node := list; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data != nil {
			entries = append(entries, convertJournalEntry(C.MmsJournalEntry(data)))
		}
	}
	return entries, moreFollows, nil
}
