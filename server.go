package iec61850

// #include <iec61850_server.h>
import "C"

import (
	"unsafe"
)

// ClientConnection represents a client connection to the server. Use from connection indication or access handlers.
type ClientConnection struct {
	c C.ClientConnection
}

// PeerAddress returns the client peer address. Valid only while the connection exists.
func (c *ClientConnection) PeerAddress() string {
	if c == nil || c.c == nil {
		return ""
	}
	s := C.ClientConnection_getPeerAddress(c.c)
	if s == nil {
		return ""
	}
	return C.GoString(s)
}

// LocalAddress returns the local server address. Valid only while the connection exists.
func (c *ClientConnection) LocalAddress() string {
	if c == nil || c.c == nil {
		return ""
	}
	s := C.ClientConnection_getLocalAddress(c.c)
	if s == nil {
		return ""
	}
	return C.GoString(s)
}

// SecurityToken returns the security token set by the authenticator, or nil.
func (c *ClientConnection) SecurityToken() unsafe.Pointer {
	if c == nil || c.c == nil {
		return nil
	}
	return C.ClientConnection_getSecurityToken(c.c)
}

// Abort closes the client connection. The connection must not be used after calling Abort.
func (c *ClientConnection) Abort() bool {
	if c == nil || c.c == nil {
		return false
	}
	return bool(C.ClientConnection_abort(c.c))
}

// ClaimOwnership claims the connection for use outside the callback; call Release when done.
func (c *ClientConnection) ClaimOwnership() *ClientConnection {
	if c == nil || c.c == nil {
		return nil
	}
	claimed := C.ClientConnection_claimOwnership(c.c)
	if claimed == nil {
		return nil
	}
	return &ClientConnection{c: claimed}
}

// Release releases ownership claimed with ClaimOwnership.
func (c *ClientConnection) Release() {
	if c == nil || c.c == nil {
		return
	}
	C.ClientConnection_release(c.c)
}

// ConnectionIndicationHandler is called when a client connects or disconnects.
type ConnectionIndicationHandler func(connection *ClientConnection, connected bool)

type IedServer struct {
	server                      C.IedServer
	serverConfig                ServerConfig
	tlsConfig                   C.TLSConfiguration
	clientAuthenticator         ClientAuthenticator
	connectionIndicationHandler ConnectionIndicationHandler
	filestoreBasepath           string // last set via SetFilestoreBasepath; C API has no getter
}

func NewServerWithTlsSupport(serverConfig ServerConfig, tlsConfig *TLSConfig, iedModel *IedModel) (*IedServer, error) {
	cTlsConfig, err := tlsConfig.createCTlsConfig()
	if err != nil {
		return nil, err
	}

	config := serverConfig.createIedServerConfig(serverConfig)
	defer C.IedServerConfig_destroy(config)
	return &IedServer{
		server:       C.IedServer_createWithConfig(iedModel.Model, cTlsConfig, config),
		serverConfig: serverConfig,
		tlsConfig:    cTlsConfig,
	}, nil
}

func NewServerWithConfig(serverConfig ServerConfig, iedModel *IedModel) *IedServer {
	config := serverConfig.createIedServerConfig(serverConfig)
	defer C.IedServerConfig_destroy(config)
	return &IedServer{
		server:       C.IedServer_createWithConfig(iedModel.Model, nil, config),
		serverConfig: serverConfig,
	}
}

// NewServer creates a new instance of the IedServer using the provided _iedModel.
func NewServer(iedModel *IedModel) *IedServer {
	return &IedServer{
		server: C.IedServer_create(iedModel.Model),
	}
}

// Start initiates the IedServer on the provided port.
func (is *IedServer) Start(port int) {
	C.IedServer_start(is.server, C.int(port))
	// If there's another way to detect the error, handle it here.
}

// StartThreadless starts the server in non-threaded mode; call ProcessIncomingData and optionally WaitReady periodically.
func (is *IedServer) StartThreadless(port int) {
	C.IedServer_startThreadless(is.server, C.int(port))
}

// StopThreadless stops the server when running in threadless mode.
func (is *IedServer) StopThreadless() {
	C.IedServer_stopThreadless(is.server)
}

// WaitReady waits until a connection has data or the timeout expires. For use with StartThreadless.
// Returns non-zero if at least one connection is ready (then call ProcessIncomingData).
func (is *IedServer) WaitReady(timeoutMs uint) int {
	return int(C.IedServer_waitReady(is.server, C.uint(timeoutMs)))
}

// ProcessIncomingData processes incoming TCP data. Call periodically when using StartThreadless.
func (is *IedServer) ProcessIncomingData() {
	C.IedServer_processIncomingData(is.server)
}

// PerformPeriodicTasks runs periodic background tasks (e.g. report timeouts). Call when using StartThreadless.
func (is *IedServer) PerformPeriodicTasks() {
	C.IedServer_performPeriodicTasks(is.server)
}

// IsRunning checks if the IedServer is currently running.
func (is *IedServer) IsRunning() bool {
	return bool(C.IedServer_isRunning(is.server))
}

// Stop terminates the IedServer.
func (is *IedServer) Stop() {
	C.IedServer_stop(is.server)
}

// Destroy frees all resources associated with the IedServer.
func (is *IedServer) Destroy() {
	C.IedServer_destroy(is.server)
}

// LockDataModel locks the data _iedModel of the IedServer.
func (is *IedServer) LockDataModel() {
	C.IedServer_lockDataModel(is.server)
}

// UnlockDataModel unlocks the data _iedModel of the IedServer.
func (is *IedServer) UnlockDataModel() {
	C.IedServer_unlockDataModel(is.server)
}

// UpdateUTCTimeAttributeValue updates a DataAttribute with a UTC time value.
func (is *IedServer) UpdateUTCTimeAttributeValue(node *ModelNode, value int64) {
	if node == nil || node._modelNode == nil {
		return
	}
	C.IedServer_updateUTCTimeAttributeValue(is.server, (*C.DataAttribute)(node._modelNode), C.uint64_t(value))
}

// UpdateFloatAttributeValue updates a DataAttribute with a float value.
func (is *IedServer) UpdateFloatAttributeValue(node *ModelNode, value float32) {
	if node == nil || node._modelNode == nil {
		return
	}
	C.IedServer_updateFloatAttributeValue(is.server, (*C.DataAttribute)(node._modelNode), C.float(value))
}

// UpdateInt32AttributeValue updates a DataAttribute with an Int32 value.
func (is *IedServer) UpdateInt32AttributeValue(node *ModelNode, value int32) {
	if node == nil || node._modelNode == nil {
		return
	}
	C.IedServer_updateInt32AttributeValue(is.server, (*C.DataAttribute)(node._modelNode), C.int32_t(value))
}

// UpdateVisibleStringAttributeValue updates a DataAttribute with a visible string value.
func (is *IedServer) UpdateVisibleStringAttributeValue(attr *DataAttribute, value string) {
	cValue, freeCValue := allocCString(value)
	defer freeCValue()
	C.IedServer_updateVisibleStringAttributeValue(is.server, attr.attribute, cValue)
}

// UpdateQuality updates the quality attribute with an UInt16 value
func (is *IedServer) UpdateQuality(node *ModelNode, quality uint16) {
	C.IedServer_updateQuality(is.server, (*C.DataAttribute)(node._modelNode), C.ushort(quality))
}

// GetAttributeValue reads the value of the attribute in the server
func (is *IedServer) GetAttributeValue(node *ModelNode) (*MmsValue, error) {
	mmsValue := C.IedServer_getAttributeValue(is.server, (*C.DataAttribute)(node._modelNode))
	mmsType := MmsType(C.MmsValue_getType(mmsValue))

	value, err := toGoValue(mmsValue, mmsType)
	if err != nil {
		return nil, err
	}
	return &MmsValue{mmsType, value}, nil
}

// GetUTCTimeAttributeValue reads the value of a time attribute in the server
func (is *IedServer) GetUTCTimeAttributeValue(node *ModelNode) int64 {
	timestamp := C.IedServer_getUTCTimeAttributeValue(is.server, (*C.DataAttribute)(node._modelNode))
	return int64(timestamp)
}

// GetNumberOfOpenConnections reads the amount of connections with the server
func (is *IedServer) GetNumberOfOpenConnections() int {
	return int(C.IedServer_getNumberOfOpenConnections(is.server))
}

// SetServerIdentity updates the server identity of the IedServer
func (is *IedServer) SetServerIdentity(vendor string, model string, version string) {
	cVendor, freeCVendor := allocCString(vendor)
	cModel, freeCModel := allocCString(model)
	cVersion, freeCVersion := allocCString(version)

	defer func() {
		freeCVendor()
		freeCModel()
		freeCVersion()
	}()

	C.IedServer_setServerIdentity(is.server, cVendor, cModel, cVersion)
}

// SetMmsLocalIpAddress sets the local IP address the MMS server will bind to. Call before Start.
func (is *IedServer) SetMmsLocalIpAddress(localIpAddress string) error {
	cAddr, freeCAddr := allocCString(localIpAddress)
	defer freeCAddr()
	C.IedServer_setLocalIpAddress(is.server, cAddr)
	return nil
}

// SetMmsTcpPort is not supported by the library; the TCP port is set when calling Start(port).
func (is *IedServer) SetMmsTcpPort(tcpPort int) error {
	_ = tcpPort
	return UserProvidedInvalidArgument
}
