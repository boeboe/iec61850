package iec61850

/*
#include "goose_publisher.h"
#include "mms_value.h"

static bool is_publisher_not_null(GoosePublisher p) {
	return p != NULL;
}

static void destroy_linked_list_val(LinkedList value) {
	LinkedList_destroyDeep(value, (LinkedListValueDeleteFunction)MmsValue_delete);
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

type (
	LinkedListValue struct {
		internalLinkedList *C.struct_sLinkedList
	}

	// CommParameters is the Go equivalent of C's struct sCommParameters (goose_publisher.h).
	// It holds VLAN, APPID, and destination MAC for GOOSE/SV. Embedded in GoosePublisherConf.
	CommParameters struct {
		VlanPriority uint8
		VlanID       uint16
		AppID        uint16
		DstAddr      [6]uint8
	}

	// GoosePublisherConf configures a GOOSE publisher. InterfaceID is the Ethernet interface (e.g. "eth0");
	// the remaining fields are the C CommParameters (embedded).
	GoosePublisherConf struct {
		InterfaceID string
		CommParameters
	}

	GoosePublisher struct {
		internalPublisher *C.struct_sGoosePublisher
	}
)

var (
	ErrCreateGoosePublisher = errors.New("can not create goose publisher")
	ErrSendGooseValue       = errors.New("can not send goose value")
)

func NewGoosePublisher(conf GoosePublisherConf) (publisher *GoosePublisher, err error) {
	return NewGoosePublisherEx(conf, true)
}

// NewGoosePublisherEx creates a GOOSE publisher with optional VLAN tag. useVlanTag false disables
// VLAN tags in sent frames when not needed. Otherwise equivalent to NewGoosePublisher.
func NewGoosePublisherEx(conf GoosePublisherConf, useVlanTag bool) (publisher *GoosePublisher, err error) {
	parameters := C.struct_sCommParameters{}
	parameters.appId = C.uint16_t(conf.AppID)
	parameters.vlanId = C.uint16_t(conf.VlanID)
	parameters.vlanPriority = C.uint8_t(conf.VlanPriority)
	for i := 0; i < len(conf.DstAddr); i++ {
		parameters.dstAddress[i] = C.uint8_t(conf.DstAddr[i])
	}
	ether, freeEther := allocCString(conf.InterfaceID)
	defer freeEther()

	cGoosePublisher := C.GoosePublisher_createEx(&parameters, ether, C.bool(useVlanTag))
	if !bool(C.is_publisher_not_null(cGoosePublisher)) {
		err = ErrCreateGoosePublisher
		return
	}

	publisher = &GoosePublisher{
		internalPublisher: cGoosePublisher,
	}

	return
}

func (receiver *GoosePublisher) SetGoCbRef(goCbRef string) {
	ref, freeRef := allocCString(goCbRef)
	defer freeRef()

	C.GoosePublisher_setGoCbRef(receiver.internalPublisher, ref)
}

// SetGoID sets the GOOSE identifier string sent in GOOSE messages (e.g. when it differs from GoCbRef).
func (receiver *GoosePublisher) SetGoID(goID string) {
	ref, freeRef := allocCString(goID)
	defer freeRef()

	C.GoosePublisher_setGoID(receiver.internalPublisher, ref)
}

func (receiver *GoosePublisher) SetDataSetRef(dataSetRef string) {
	ref, freeRef := allocCString(dataSetRef)
	defer freeRef()

	C.GoosePublisher_setDataSetRef(receiver.internalPublisher, ref)
}

func (receiver *GoosePublisher) SetConfRev(confRef uint32) {
	C.GoosePublisher_setConfRev(receiver.internalPublisher, C.uint32_t(confRef))
}

func (receiver *GoosePublisher) SetTimeAllowedToLive(timeAllowedToLive uint32) {
	C.GoosePublisher_setTimeAllowedToLive(receiver.internalPublisher, C.uint32_t(timeAllowedToLive))
}

func (receiver *GoosePublisher) SetSimulation(simulation bool) {
	C.GoosePublisher_setSimulation(receiver.internalPublisher, C.bool(simulation))
}

func (receiver *GoosePublisher) SetStNum(stNum uint32) {
	C.GoosePublisher_setStNum(receiver.internalPublisher, C.uint32_t(stNum))
}

func (receiver *GoosePublisher) SetSqNum(sqNum uint32) {
	C.GoosePublisher_setSqNum(receiver.internalPublisher, C.uint32_t(sqNum))
}

func (receiver *GoosePublisher) SetNeedsCommission(ndsCom bool) {
	C.GoosePublisher_setNeedsCommission(receiver.internalPublisher, C.bool(ndsCom))
}

func (receiver *GoosePublisher) IncreaseStNum() {
	C.GoosePublisher_increaseStNum(receiver.internalPublisher)
}

func (receiver *GoosePublisher) Reset() {
	C.GoosePublisher_reset(receiver.internalPublisher)
}

func (receiver *GoosePublisher) Publish(dataSet *LinkedListValue) error {
	if int(C.GoosePublisher_publish(receiver.internalPublisher, dataSet.internalLinkedList)) == -1 {
		return ErrSendGooseValue
	}

	return nil
}

// PublishAndDump publishes a GOOSE message and copies the raw encoded payload into msgBuf.
// Returns the number of bytes written into msgBuf (use msgBuf[:msgLen]) or an error if publish failed.
// msgBuf must be non-nil and have positive length; use a large enough buffer (e.g. 1500+ bytes for Ethernet).
func (receiver *GoosePublisher) PublishAndDump(dataSet *LinkedListValue, msgBuf []byte) (msgLen int, err error) {
	if len(msgBuf) == 0 {
		return 0, ErrSendGooseValue
	}
	var cMsgLen C.int32_t
	rc := C.GoosePublisher_publishAndDump(receiver.internalPublisher, dataSet.internalLinkedList, (*C.char)(unsafe.Pointer(&msgBuf[0])), &cMsgLen, C.int32_t(len(msgBuf)))
	if rc == -1 {
		return 0, ErrSendGooseValue
	}
	return int(cMsgLen), nil
}

func (receiver *GoosePublisher) Close() {
	C.GoosePublisher_destroy(receiver.internalPublisher)
}

func NewLinkedListValue() *LinkedListValue {
	return &LinkedListValue{
		internalLinkedList: C.LinkedList_create(),
	}
}

func (receiver *LinkedListValue) Add(value *MmsValue) error {
	rawVal, err := toMmsValue(value.Type, value.Value)
	if err != nil {
		return err
	}
	C.LinkedList_add(receiver.internalLinkedList, unsafe.Pointer(rawVal))

	return nil
}

func (receiver *LinkedListValue) Size() int {
	return int(C.LinkedList_size(receiver.internalLinkedList))
}

func (receiver *LinkedListValue) Destroy() {
	C.destroy_linked_list_val(receiver.internalLinkedList)
}
