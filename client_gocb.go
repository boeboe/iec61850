package iec61850

// #include <iec61850_client.h>
// #include <iec61850_common.h>
import "C"

import (
	"errors"
)

// GOOSE Control Block element masks
const (
	GoCBElementGoEna      = 1   // GOCB_ELEMENT_GO_ENA
	GoCBElementGoID       = 2   // GOCB_ELEMENT_GO_ID
	GoCBElementDatSet     = 4   // GOCB_ELEMENT_DATSET
	GoCBElementConfRev    = 8   // GOCB_ELEMENT_CONF_REV
	GoCBElementNdsComm    = 16  // GOCB_ELEMENT_NDS_COMM
	GoCBElementDstAddress = 32  // GOCB_ELEMENT_DST_ADDRESS
	GoCBElementMinTime    = 64  // GOCB_ELEMENT_MIN_TIME
	GoCBElementMaxTime    = 128 // GOCB_ELEMENT_MAX_TIME
	GoCBElementFixedOffs  = 256 // GOCB_ELEMENT_FIXED_OFFS
	GoCBElementAll        = 511 // GOCB_ELEMENT_ALL
)

// PhyComAddress represents Ethernet address and VLAN attributes for GOOSE/SV
type PhyComAddress struct {
	Addr         [6]byte // MAC address (6 bytes)
	VlanPriority uint8   // VLAN priority (0-7)
	VlanId       uint16  // VLAN ID
	AppId        uint16  // Application ID (APPID)
}

// ClientGooseControlBlockValues holds the values of a GOOSE control block
type ClientGooseControlBlockValues struct {
	GoEna      bool          // GOOSE enable
	GoID       string        // GOOSE identifier
	DatSet     string        // Dataset reference
	ConfRev    uint32        // Configuration revision
	NdsComm    bool          // Needs commissioning
	MinTime    uint32        // Minimum time between GOOSE messages (ms)
	MaxTime    uint32        // Maximum time between GOOSE messages (ms)
	FixedOffs  bool          // Fixed offsets flag
	DstAddress PhyComAddress // Destination physical communication address
}

// GetGoCBValues reads the GOOSE control block values from the server
//
// Parameters:
//   - goCBReference: IEC 61850-7-2 ACSI object reference (e.g., "simpleIOGenericIO/LLN0.gcbEvents")
//
// Returns:
//   - ClientGooseControlBlockValues containing the current GoCB values
//   - error if the operation fails
func (c *Client) GetGoCBValues(goCBReference string) (*ClientGooseControlBlockValues, error) {
	var clientError C.IedClientError
	cGoCBRef, freeCGoCBRef := allocCString(goCBReference)
	defer freeCGoCBRef()

	goCB := C.IedConnection_getGoCBValues(c.conn, &clientError, cGoCBRef, nil)
	if goCB == nil {
		return nil, GetIedClientError(clientError)
	}
	defer C.ClientGooseControlBlock_destroy(goCB)

	// Get destination address
	cDstAddr := C.ClientGooseControlBlock_getDstAddress(goCB)
	dstAddr := PhyComAddress{
		VlanPriority: uint8(cDstAddr.vlanPriority),
		VlanId:       uint16(cDstAddr.vlanId),
		AppId:        uint16(cDstAddr.appId),
	}
	// Copy MAC address
	for i := 0; i < 6; i++ {
		dstAddr.Addr[i] = byte(cDstAddr.dstAddress[i])
	}

	return &ClientGooseControlBlockValues{
		GoEna:      bool(C.ClientGooseControlBlock_getGoEna(goCB)),
		GoID:       C.GoString(C.ClientGooseControlBlock_getGoID(goCB)),
		DatSet:     C.GoString(C.ClientGooseControlBlock_getDatSet(goCB)),
		ConfRev:    uint32(C.ClientGooseControlBlock_getConfRev(goCB)),
		NdsComm:    bool(C.ClientGooseControlBlock_getNdsComm(goCB)),
		MinTime:    uint32(C.ClientGooseControlBlock_getMinTime(goCB)),
		MaxTime:    uint32(C.ClientGooseControlBlock_getMaxTime(goCB)),
		FixedOffs:  bool(C.ClientGooseControlBlock_getFixedOffs(goCB)),
		DstAddress: dstAddr,
	}, nil
}

// SetGoCBValues writes GOOSE control block values to the server
//
// Parameters:
//   - goCBReference: IEC 61850-7-2 ACSI object reference (e.g., "simpleIOGenericIO/LLN0.gcbEvents")
//   - values: The GoCB values to write
//   - parametersMask: Bitmask specifying which parameters to write (use GoCBElement* constants)
//   - singleRequest: If true, use single MMS write request; if false, use multiple requests
//
// Returns:
//   - error if the operation fails
//
// Note: Only GoEna, GoID, DatSet, and DstAddress are typically writable on most servers.
// Other attributes are usually read-only.
func (c *Client) SetGoCBValues(goCBReference string, values *ClientGooseControlBlockValues, parametersMask uint32, singleRequest bool) error {
	var clientError C.IedClientError
	cGoCBRef, freeCGoCBRef := allocCString(goCBReference)
	defer freeCGoCBRef()

	// Create a ClientGooseControlBlock instance
	goCB := C.ClientGooseControlBlock_create(cGoCBRef)
	if goCB == nil {
		return errors.New("failed to create ClientGooseControlBlock")
	}
	defer C.ClientGooseControlBlock_destroy(goCB)

	// Set the values according to the parameter mask
	if parametersMask&GoCBElementGoEna != 0 {
		C.ClientGooseControlBlock_setGoEna(goCB, C.bool(values.GoEna))
	}
	if parametersMask&GoCBElementGoID != 0 {
		cGoID, freeCGoID := allocCString(values.GoID)
		defer freeCGoID()
		C.ClientGooseControlBlock_setGoID(goCB, cGoID)
	}
	if parametersMask&GoCBElementDatSet != 0 {
		cDatSet, freeCDatSet := allocCString(values.DatSet)
		defer freeCDatSet()
		C.ClientGooseControlBlock_setDatSet(goCB, cDatSet)
	}
	if parametersMask&GoCBElementDstAddress != 0 {
		// Create C PhyComAddress from Go struct
		var cDstAddr C.PhyComAddress
		cDstAddr.vlanPriority = C.uint8_t(values.DstAddress.VlanPriority)
		cDstAddr.vlanId = C.uint16_t(values.DstAddress.VlanId)
		cDstAddr.appId = C.uint16_t(values.DstAddress.AppId)
		for i := 0; i < 6; i++ {
			cDstAddr.dstAddress[i] = C.uint8_t(values.DstAddress.Addr[i])
		}
		C.ClientGooseControlBlock_setDstAddress(goCB, cDstAddr)
	}

	// Write to server
	C.IedConnection_setGoCBValues(c.conn, &clientError, goCB, C.uint32_t(parametersMask), C.bool(singleRequest))

	if clientError != C.IED_ERROR_OK {
		return GetIedClientError(clientError)
	}

	return nil
}
