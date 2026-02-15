package iec61850

/*
#include "goose_receiver.h"
#include "goose_subscriber.h"
#include <stdint.h>

extern void cgoReportCallbackBridgeDispatcher(GooseSubscriber subscriber, void *parameter);

static void goose_report_proxy_handler(GooseSubscriber subscriber, void* parameter) {
	cgoReportCallbackBridgeDispatcher(subscriber, parameter);
}

static void simple_goose_subscriber_set_listener(GooseSubscriber subscriber, uintptr_t parameter) {
	GooseSubscriber_setListener(subscriber, goose_report_proxy_handler, (void *)parameter);
}
*/
import "C"
import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type (
	// GooseReceiverSocket is the opaque handle returned by StartThreadless.
	// It represents the Ethernet socket used for receiving; the receiver parses
	// frames when you call HandleMessage with a buffer.
	GooseReceiverSocket struct {
		c C.EthernetSocket
	}

	GooseReceiver struct {
		noCopy        struct{}
		gooseReceiver *C.struct_sGooseReceiver
		refs          map[GooseCallbackHandlerID]struct{}
		keepBuffer    []byte // optional buffer passed to createEx; keep alive for receiver lifetime
	}
)

var (
	gooseCallbackLocker struct {
		noCopy struct{}
		sync.RWMutex
		idOffset     atomic.Uintptr
		callbackRefs map[GooseCallbackHandlerID]struct {
			handler    GooseReportCallback
			subscriber *GooseSubscriber
		}
	}
)

func init() {
	gooseCallbackLocker.Lock()
	defer gooseCallbackLocker.Unlock()

	gooseCallbackLocker.idOffset.Add(1000)
	gooseCallbackLocker.callbackRefs = map[GooseCallbackHandlerID]struct {
		handler    GooseReportCallback
		subscriber *GooseSubscriber
	}{}
}

//export cgoReportCallbackBridgeDispatcher
func cgoReportCallbackBridgeDispatcher(_ *C.struct_sGooseSubscriber, parameter unsafe.Pointer) {
	refID := GooseCallbackHandlerID(parameter)
	gooseCallbackLocker.RLock()
	defer gooseCallbackLocker.RUnlock()

	if fetch, ok := gooseCallbackLocker.callbackRefs[refID]; ok {
		fetch.handler(&GooseReport{
			parameter:       parameter,
			GooseSubscriber: fetch.subscriber,
		})
	}
}

func NewGooseReceiver() *GooseReceiver {
	return newGooseReceiverWithBuffer(nil)
}

// NewGooseReceiverEx creates a GOOSE receiver that uses the given buffer for message handling
// instead of allocating its own. Pass nil or an empty slice to use the default (library-allocated) buffer.
// When buffer is non-nil, the receiver keeps a reference to it for its lifetime; do not modify
// the buffer while the receiver is in use.
func NewGooseReceiverEx(buffer []byte) *GooseReceiver {
	return newGooseReceiverWithBuffer(buffer)
}

func newGooseReceiverWithBuffer(buffer []byte) *GooseReceiver {
	var cReceiver *C.struct_sGooseReceiver
	var keepBuf []byte
	if len(buffer) > 0 {
		cReceiver = C.GooseReceiver_createEx((*C.uint8_t)(unsafe.Pointer(&buffer[0])))
		keepBuf = buffer
	} else {
		cReceiver = C.GooseReceiver_create()
	}
	return &GooseReceiver{
		gooseReceiver: cReceiver,
		refs:          make(map[GooseCallbackHandlerID]struct{}),
		keepBuffer:    keepBuf,
	}
}

func (receiver *GooseReceiver) AddSubscriber(subscriber *GooseSubscriber) *GooseReceiver {
	gooseCallbackLocker.Lock()
	defer gooseCallbackLocker.Unlock()

	gooseCallbackLocker.callbackRefs[subscriber.HandlerID] = struct {
		handler    GooseReportCallback
		subscriber *GooseSubscriber
	}{
		handler:    subscriber.Conf.ReportHandler,
		subscriber: subscriber,
	}
	receiver.refs[subscriber.HandlerID] = struct{}{}
	C.simple_goose_subscriber_set_listener(
		subscriber.subscriber,
		C.uintptr_t(subscriber.HandlerID),
	)
	C.GooseReceiver_addSubscriber(receiver.gooseReceiver, subscriber.subscriber)

	return receiver
}

func (receiver *GooseReceiver) RemoveSubscriber(subscriber *GooseSubscriber) *GooseReceiver {
	gooseCallbackLocker.Lock()
	defer gooseCallbackLocker.Unlock()

	C.GooseReceiver_removeSubscriber(receiver.gooseReceiver, subscriber.subscriber)
	delete(gooseCallbackLocker.callbackRefs, subscriber.HandlerID)
	delete(receiver.refs, subscriber.HandlerID)

	return receiver
}

func (receiver *GooseReceiver) SetInterfaceID(interfaceID string) *GooseReceiver {
	tmp, freeTmp := allocCString(interfaceID)
	defer freeTmp()
	C.GooseReceiver_setInterfaceId(receiver.gooseReceiver, tmp)

	return receiver
}

func (receiver *GooseReceiver) GetInterfaceID() string {
	return C.GoString(C.GooseReceiver_getInterfaceId(receiver.gooseReceiver))
}

func (receiver *GooseReceiver) Start() *GooseReceiver {
	C.GooseReceiver_start(receiver.gooseReceiver)

	return receiver
}

func (receiver *GooseReceiver) IsRunning() bool {
	return bool(C.GooseReceiver_isRunning(receiver.gooseReceiver))
}

func (receiver *GooseReceiver) Tick() bool {
	return bool(C.GooseReceiver_tick(receiver.gooseReceiver))
}

// StartThreadless starts the GOOSE receiver in non-threaded mode. The returned socket
// handle can be used with external read loops; call HandleMessage with each received
// Ethernet frame. Call StopThreadless to stop.
func (receiver *GooseReceiver) StartThreadless() *GooseReceiverSocket {
	sock := C.GooseReceiver_startThreadless(receiver.gooseReceiver)
	if sock == nil {
		return nil
	}
	return &GooseReceiverSocket{c: sock}
}

// StopThreadless stops the receiver when running in threadless mode (after StartThreadless).
func (receiver *GooseReceiver) StopThreadless() {
	C.GooseReceiver_stopThreadless(receiver.gooseReceiver)
}

// HandleMessage parses a GOOSE message from a raw Ethernet frame. Use this when driving
// reception yourself (e.g. with StartThreadless or custom socket reads). buffer must
// contain the complete Ethernet frame.
func (receiver *GooseReceiver) HandleMessage(buffer []byte) {
	if len(buffer) == 0 {
		return
	}
	C.GooseReceiver_handleMessage(receiver.gooseReceiver, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), C.int(len(buffer)))
}

func (receiver *GooseReceiver) Stop() *GooseReceiver {
	C.GooseReceiver_stop(receiver.gooseReceiver)

	return receiver
}

func (receiver *GooseReceiver) Destroy() {
	gooseCallbackLocker.Lock()
	defer gooseCallbackLocker.Unlock()
	for id := range receiver.refs {
		delete(gooseCallbackLocker.callbackRefs, id)
	}
	C.GooseReceiver_destroy(receiver.gooseReceiver)
	receiver.refs = nil
	receiver.gooseReceiver = nil
	receiver.keepBuffer = nil
}
