# IEC 61850 Go Bindings - Functions Reference

**Version**: 1.6.1  
**Generated**: February 7, 2026

This document provides comprehensive documentation for all exported Go functions in the iec61850 package, including their corresponding C functions, descriptions, and usage examples.

---

## Table of Contents

1. [Client Connection Functions](#client-connection-functions)
2. [Client Read/Write Functions](#client-readwrite-functions)
3. [MMS Connection Functions](#mms-connection-functions)
4. [MMS Value Functions](#mms-value-functions)
5. [Server Functions](#server-functions)
6. [Control Functions](#control-functions)
7. [Reporting Functions](#reporting-functions)
8. [GOOSE Functions](#goose-functions)
9. [Sampled Values (SV) Functions](#sampled-values-sv-functions)
10. [File Services Functions](#file-services-functions)
11. [Dataset Functions](#dataset-functions)
12. [Utility Functions](#utility-functions)

---

## Client Connection Functions

### NewClient

**Go Function**: `func NewClient(settings Settings) (*Client, error)`  
**C Function**: `IedConnection_create()`, `IedConnection_connect()`

**Description**: Creates a new IEC 61850 client connection with the specified settings.

**Parameters**:
- `settings` - Connection settings (host, port, timeouts)

**Returns**: Client instance and error

**Example**:
```go
settings := iec61850.Settings{
    Host:           "192.168.1.10",
    Port:           102,
    ConnectTimeout: 10000, // milliseconds
    RequestTimeout: 5000,  // milliseconds
}

client, err := iec61850.NewClient(settings)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

---

### NewClientWithTlsSupport

**Go Function**: `func NewClientWithTlsSupport(settings Settings, tlsConfig *TLSConfig) (*Client, error)`  
**C Function**: `IedConnection_createWithTlsSupport()`, `IedConnection_connect()`

**Description**: Creates a new IEC 61850 client connection with TLS encryption.

**Parameters**:
- `settings` - Connection settings
- `tlsConfig` - TLS configuration (certificates, keys, validation)

**Returns**: Client instance and error

**Example**:
```go
tlsConfig := &iec61850.TLSConfig{
    ChainValidation: true,
    OwnCertificate:  certPEM,
    OwnKey:          keyPEM,
    CACerts:         [][]byte{caCertPEM},
}

client, err := iec61850.NewClientWithTlsSupport(settings, tlsConfig)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

---

### NewClientWithDefaultSettings

**Go Function**: `func NewClientWithDefaultSettings() (*Client, error)`  
**C Function**: `IedConnection_create()`, `IedConnection_connect()`

**Description**: Creates a client with default settings (localhost:102).

**Example**:
```go
client, err := iec61850.NewClientWithDefaultSettings()
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

---

### Close

**Go Function**: `func (c *Client) Close()`  
**C Function**: `IedConnection_destroy()`

**Description**: Closes the client connection and releases all resources.

**Example**:
```go
client, _ := iec61850.NewClient(settings)
defer client.Close()
```

---

### GetState

**Go Function**: `func (c *Client) GetState() IedConnectionState`  
**C Function**: `IedConnection_getState()`

**Description**: Returns the current connection state.

**Returns**: Connection state (Closed, Connecting, Connected, Closing)

**Example**:
```go
state := client.GetState()
if state == iec61850.IedStateConnected {
    fmt.Println("Connected!")
}
```

---

### GetLastApplError

**Go Function**: `func (c *Client) GetLastApplError() LastApplError`  
**C Function**: `IedConnection_getLastApplError()`

**Description**: Returns the last application error from control operations.

**Returns**: LastApplError structure with error details

**Example**:
```go
lastError := client.GetLastApplError()
fmt.Printf("Error %d, AddCause: %d\n", lastError.Error, lastError.AddCause)
```

---

### GetRequestTimeout

**Go Function**: `func (c *Client) GetRequestTimeout() uint32`  
**C Function**: `IedConnection_getRequestTimeout()`

**Description**: Returns the current request timeout in milliseconds.

**Example**:
```go
timeout := client.GetRequestTimeout()
fmt.Printf("Request timeout: %d ms\n", timeout)
```

---

## Client Read/Write Functions

### Read

**Go Function**: `func (c *Client) Read(objectRef string, fc FC) (interface{}, error)`  
**C Function**: `IedConnection_readObject()`

**Description**: Reads an IEC 61850 data attribute value.

**Parameters**:
- `objectRef` - Object reference (e.g., "simpleIOGenericIO/GGIO1.AnIn1.mag.f")
- `fc` - Functional constraint (ST, MX, SP, etc.)

**Returns**: Value as interface{} and error

**Example**:
```go
value, err := client.Read("simpleIOGenericIO/GGIO1.AnIn1.mag.f", iec61850.MX)
if err != nil {
    log.Fatal(err)
}
floatValue := value.(float32)
fmt.Printf("Value: %f\n", floatValue)
```

---

### ReadBool

**Go Function**: `func (c *Client) ReadBool(objectRef string, fc FC) (bool, error)`  
**C Function**: `IedConnection_readObject()`

**Description**: Reads a boolean value from the server.

**Example**:
```go
stVal, err := client.ReadBool("simpleIOGenericIO/GGIO1.SPCSO1.stVal", iec61850.ST)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Switch state: %v\n", stVal)
```

---

### ReadInt32

**Go Function**: `func (c *Client) ReadInt32(objectRef string, fc FC) (int32, error)`  
**C Function**: `IedConnection_readObject()`

**Description**: Reads a 32-bit integer value.

**Example**:
```go
value, err := client.ReadInt32("simpleIOGenericIO/GGIO1.IntIn1.stVal", iec61850.ST)
```

---

### ReadInt64

**Go Function**: `func (c *Client) ReadInt64(objectRef string, fc FC) (int64, error)`  
**C Function**: `IedConnection_readObject()`

**Description**: Reads a 64-bit integer value.

---

### ReadFloat32

**Go Function**: `func (c *Client) ReadFloat32(objectRef string, fc FC) (float32, error)`  
**C Function**: `IedConnection_readObject()`

**Description**: Reads a 32-bit floating point value.

**Example**:
```go
magnitude, err := client.ReadFloat32("simpleIOGenericIO/GGIO1.AnIn1.mag.f", iec61850.MX)
```

---

### ReadFloat64

**Go Function**: `func (c *Client) ReadFloat64(objectRef string, fc FC) (float64, error)`  
**C Function**: `IedConnection_readObject()`

**Description**: Reads a 64-bit floating point value.

---

### ReadString

**Go Function**: `func (c *Client) ReadString(objectRef string, fc FC) (string, error)`  
**C Function**: `IedConnection_readObject()`

**Description**: Reads a string value (VisibleString or MmsString).

**Example**:
```go
description, err := client.ReadString("simpleIOGenericIO/LLN0.NamPlt.vendor", iec61850.DC)
```

---

### ReadBitString

**Go Function**: `func (c *Client) ReadBitString(objectRef string, fc FC) ([]byte, error)`  
**C Function**: `IedConnection_readObject()`

**Description**: Reads a bitstring value as byte array.

---

### ReadDataSet

**Go Function**: `func (c *Client) ReadDataSet(objectRef string) ([]*MmsValue, error)`  
**C Function**: `IedConnection_readDataSetValues()`

**Description**: Reads all values from a dataset.

**Example**:
```go
values, err := client.ReadDataSet("simpleIOGenericIO/LLN0.dataset1")
if err != nil {
    log.Fatal(err)
}
for _, val := range values {
    fmt.Printf("Value: %v\n", val)
}
```

---

### Write

**Go Function**: `func (c *Client) Write(objectRef string, fc FC, value interface{}) error`  
**C Function**: `IedConnection_writeObject()`

**Description**: Writes a value to an IEC 61850 data attribute.

**Parameters**:
- `objectRef` - Object reference
- `fc` - Functional constraint
- `value` - Value to write (bool, int, float, string, etc.)

**Example**:
```go
err := client.Write("simpleIOGenericIO/GGIO1.AnIn1.mag.f", iec61850.MX, float32(42.5))
if err != nil {
    log.Fatal(err)
}
```

---

## MMS Connection Functions

### NewMmsConnection

**Go Function**: `func NewMmsConnection() *MmsConnection`  
**C Function**: `MmsConnection_create()`

**Description**: Creates a new MMS connection (non-TLS, threaded mode).

**Example**:
```go
mmsConn := iec61850.NewMmsConnection()
defer mmsConn.Destroy()

err := mmsConn.ConnectAsync("192.168.1.10", 102, func(err error) {
    if err != nil {
        log.Println("Connection failed:", err)
    } else {
        log.Println("Connected!")
    }
})
```

---

### NewMmsConnectionSecure

**Go Function**: `func NewMmsConnectionSecure(tlsConfig *TLSConfiguration) *MmsConnection`  
**C Function**: `MmsConnection_createSecure()`

**Description**: Creates a TLS-secured MMS connection.

**Example**:
```go
tlsCfg := &iec61850.TLSConfiguration{
    ChainValidation: true,
    OwnCertificate:  certPEM,
    OwnKey:          keyPEM,
}
mmsConn := iec61850.NewMmsConnectionSecure(tlsCfg)
defer mmsConn.Destroy()
```

---

### NewMmsConnectionNonThreaded

**Go Function**: `func NewMmsConnectionNonThreaded(tlsConfig *TLSConfiguration) *MmsConnection`  
**C Function**: `MmsConnection_createNonThreaded()`

**Description**: Creates an MMS connection in non-threaded mode (requires calling `Tick()`).

**Example**:
```go
mmsConn := iec61850.NewMmsConnectionNonThreaded(nil)
defer mmsConn.Destroy()

// In event loop:
for {
    mmsConn.Tick()
    time.Sleep(10 * time.Millisecond)
}
```

---

### Destroy

**Go Function**: `func (c *MmsConnection) Destroy()`  
**C Function**: `MmsConnection_destroy()`

**Description**: Destroys the MMS connection and releases resources.

---

### SetConnectTimeout

**Go Function**: `func (c *MmsConnection) SetConnectTimeout(timeoutMs uint32)`  
**C Function**: `MmsConnection_setConnectTimeout()`

**Description**: Sets connection timeout in milliseconds.

---

### SetRequestTimeout

**Go Function**: `func (c *MmsConnection) SetRequestTimeout(timeoutMs uint32)`  
**C Function**: `MmsConnection_setRequestTimeout()`

**Description**: Sets request timeout in milliseconds.

---

### GetRequestTimeout

**Go Function**: `func (c *MmsConnection) GetRequestTimeout() uint32`  
**C Function**: `MmsConnection_getRequestTimeout()`

**Description**: Gets current request timeout.

---

### ConnectAsync

**Go Function**: `func (c *MmsConnection) ConnectAsync(hostname string, port int, callback func(error)) error`  
**C Function**: `MmsConnection_connectAsync()`

**Description**: Initiates asynchronous connection to MMS server.

**Example**:
```go
err := mmsConn.ConnectAsync("192.168.1.10", 102, func(err error) {
    if err != nil {
        log.Println("Failed:", err)
    } else {
        log.Println("Connected successfully")
    }
})
```

---

### Conclude

**Go Function**: `func (c *MmsConnection) Conclude() error`  
**C Function**: `MmsConnection_conclude()`

**Description**: Sends MMS conclude service to gracefully close the association.

---

### ConcludeAsync

**Go Function**: `func (c *MmsConnection) ConcludeAsync(callback func(error)) error`  
**C Function**: `MmsConnection_concludeAsync()`

**Description**: Asynchronous version of Conclude.

---

### AbortAsync

**Go Function**: `func (c *MmsConnection) AbortAsync() error`  
**C Function**: `MmsConnection_abortAsync()`

**Description**: Aborts the MMS connection asynchronously.

---

### Tick

**Go Function**: `func (c *MmsConnection) Tick() bool`  
**C Function**: `MmsConnection_tick()`

**Description**: Processes connection events for non-threaded mode. Returns true if more work pending.

**Example**:
```go
for {
    morePending := mmsConn.Tick()
    if !morePending {
        time.Sleep(10 * time.Millisecond)
    }
}
```

---

### ReadVariableAsync

**Go Function**: `func (c *MmsConnection) ReadVariableAsync(domainID, itemID string, callback func(*MmsValue, error)) error`  
**C Function**: `MmsConnection_readVariableAsync()`

**Description**: Asynchronously reads an MMS variable.

**Example**:
```go
err := mmsConn.ReadVariableAsync("domain", "variable", func(val *MmsValue, err error) {
    if err != nil {
        log.Println("Read failed:", err)
        return
    }
    fmt.Printf("Value: %v\n", val)
})
```

---

### WriteVariableAsync

**Go Function**: `func (c *MmsConnection) WriteVariableAsync(domainID, itemID string, value *MmsValueRef, callback func(error)) error`  
**C Function**: `MmsConnection_writeVariableAsync()`

**Description**: Asynchronously writes an MMS variable.

---

### GetDomainNamesAsync

**Go Function**: `func (c *MmsConnection) GetDomainNamesAsync(callback func([]string, error)) error`  
**C Function**: `MmsConnection_getDomainNamesAsync()`

**Description**: Asynchronously retrieves all domain names from the server.

---

### GetDomainVariableNamesAsync

**Go Function**: `func (c *MmsConnection) GetDomainVariableNamesAsync(domainID string, callback func([]string, error)) error`  
**C Function**: `MmsConnection_getDomainVariableNamesAsync()`

**Description**: Asynchronously retrieves variable names in a domain.

---

### IdentifyAsync

**Go Function**: `func (c *MmsConnection) IdentifyAsync(callback func(vendorName, modelName, revision string, err error)) error`  
**C Function**: `MmsConnection_identifyAsync()`

**Description**: Asynchronously retrieves server identification.

**Example**:
```go
err := mmsConn.IdentifyAsync(func(vendor, model, revision string, err error) {
    if err == nil {
        fmt.Printf("Server: %s %s (Rev %s)\n", vendor, model, revision)
    }
})
```

---

### SetRawMessageHandler

**Go Function**: `func (c *MmsConnection) SetRawMessageHandler(callback func(message []byte, received bool))`  
**C Function**: `MmsConnection_setRawMessageHandler()`

**Description**: Sets handler to intercept raw MMS messages (for debugging/logging).

**Example**:
```go
mmsConn.SetRawMessageHandler(func(message []byte, received bool) {
    direction := "SENT"
    if received {
        direction = "RECV"
    }
    fmt.Printf("[%s] %d bytes: %x\n", direction, len(message), message)
})
```

---

### GetIsoConnectionParameters

**Go Function**: `func (c *MmsConnection) GetIsoConnectionParameters() *IsoConnectionParameters`  
**C Function**: `MmsConnection_getIsoConnectionParameters()`

**Description**: Retrieves ISO connection parameters (selectors, AP titles).

---

### GetMmsConnectionParameters

**Go Function**: `func (c *MmsConnection) GetMmsConnectionParameters() *MmsConnectionParameters`  
**C Function**: `MmsConnection_getMmsConnectionParameters()`

**Description**: Retrieves MMS connection parameters (PDU size, outstanding calls).

**Example**:
```go
params := mmsConn.GetMmsConnectionParameters()
fmt.Printf("Max PDU Size: %d\n", params.MaxPduSize)
```

---

## MMS Value Functions

### NewMmsValue

**Go Function**: `func NewMmsValue(mmsType MmsType, value interface{}) (*MmsValue, error)`  
**C Function**: Various constructors (`MmsValue_newInteger()`, `MmsValue_newBoolean()`, etc.)

**Description**: Creates a new MmsValue of the specified type.

**Example**:
```go
intVal, _ := iec61850.NewMmsValue(iec61850.Integer, int64(42))
boolVal, _ := iec61850.NewMmsValue(iec61850.Boolean, true)
strVal, _ := iec61850.NewMmsValue(iec61850.VisibleString, "Hello")
```

---

### NewMmsValueBitString

**Go Function**: `func NewMmsValueBitString(bitSize int) *MmsValueRef`  
**C Function**: `MmsValue_newBitString()`

**Description**: Creates a new bitstring value.

**Example**:
```go
bitStr := iec61850.NewMmsValueBitString(16)
bitStr.SetBitStringFromInteger(0x0F0F)
```

---

### NewMmsValueVisibleString

**Go Function**: `func NewMmsValueVisibleString(s string) *MmsValueRef`  
**C Function**: `MmsValue_newVisibleString()`

**Description**: Creates a visible string value.

---

### NewMmsValueUtcTimeByMsTime

**Go Function**: `func NewMmsValueUtcTimeByMsTime(ms uint64) *MmsValueRef`  
**C Function**: `MmsValue_newUtcTimeByMsTime()`

**Description**: Creates a UTC time value from milliseconds since epoch.

**Example**:
```go
timeVal := iec61850.NewMmsValueUtcTimeByMsTime(uint64(time.Now().UnixMilli()))
```

---

### MmsValueCreateArray

**Go Function**: `func MmsValueCreateArray(elementType *MmsVariableSpecificationRef, size int) *MmsValueRef`  
**C Function**: `MmsValue_createArray()`

**Description**: Creates an array of MmsValues.

---

### GetType

**Go Function**: `func (r *MmsValueRef) GetType() MmsType`  
**C Function**: `MmsValue_getType()`

**Description**: Returns the MMS type of the value.

---

### ToInt64

**Go Function**: `func (r *MmsValueRef) ToInt64() int64`  
**C Function**: `MmsValue_toInt64()`

**Description**: Converts value to int64.

---

### ToUint32

**Go Function**: `func (r *MmsValueRef) ToUint32() uint32`  
**C Function**: `MmsValue_toUint32()`

**Description**: Converts value to uint32.

---

### ToDouble

**Go Function**: `func (r *MmsValueRef) ToDouble() float64`  
**C Function**: `MmsValue_toDouble()`

**Description**: Converts value to float64.

---

### GetBitStringAsInteger

**Go Function**: `func (r *MmsValueRef) GetBitStringAsInteger() uint32`  
**C Function**: `MmsValue_getBitStringAsInteger()`

**Description**: Gets bitstring value as integer (little-endian).

---

### GetBitStringAsIntegerBigEndian

**Go Function**: `func (r *MmsValueRef) GetBitStringAsIntegerBigEndian() uint32`  
**C Function**: `MmsValue_getBitStringAsIntegerBigEndian()`

**Description**: Gets bitstring value as integer (big-endian).

---

### SetBitStringFromInteger

**Go Function**: `func (r *MmsValueRef) SetBitStringFromInteger(val uint32)`  
**C Function**: `MmsValue_setBitStringFromInteger()`

**Description**: Sets bitstring from integer value (little-endian).

---

### GetElement

**Go Function**: `func (r *MmsValueRef) GetElement(index int) *MmsValueRef`  
**C Function**: `MmsValue_getElement()`

**Description**: Gets an element from array or structure.

**Example**:
```go
arrayVal := iec61850.MmsValueCreateEmptyArray(3)
elem0 := arrayVal.GetElement(0)
```

---

### SetElement

**Go Function**: `func (r *MmsValueRef) SetElement(index int, value *MmsValueRef)`  
**C Function**: `MmsValue_setElement()`

**Description**: Sets an element in array or structure.

---

### GetDataAccessError

**Go Function**: `func (r *MmsValueRef) GetDataAccessError() MmsDataAccessError`  
**C Function**: `MmsValue_getDataAccessError()`

**Description**: Gets data access error code from value.

---

### EncodeMmsData

**Go Function**: `func (r *MmsValueRef) EncodeMmsData(buffer []byte, startPos int, encode bool) int`  
**C Function**: `MmsValue_encodeMmsData()`

**Description**: Encodes MmsValue to binary buffer.

---

### DecodeMmsData

**Go Function**: `func DecodeMmsData(buffer []byte, startPos, length int) (value *MmsValueRef, endPos int)`  
**C Function**: `MmsValue_decodeMmsData()`

**Description**: Decodes MmsValue from binary buffer.

---

## Server Functions

### NewServer

**Go Function**: `func NewServer(iedModel *IedModel) *IedServer`  
**C Function**: `IedServer_create()`

**Description**: Creates a new IED server instance.

**Example**:
```go
model := ied61850.LoadModelFromFile("model.cfg")
server := iec61850.NewServer(model)
defer server.Destroy()
```

---

### NewServerWithConfig

**Go Function**: `func NewServerWithConfig(serverConfig ServerConfig, iedModel *IedModel) *IedServer`  
**C Function**: `IedServer_createWithConfig()`

**Description**: Creates server with custom configuration.

**Example**:
```go
config := iec61850.ServerConfig{
    ReportBufferSize: 100000,
    MaxConnections:   10,
    EnableFileService: true,
}
server := iec61850.NewServerWithConfig(config, model)
```

---

### NewServerWithTlsSupport

**Go Function**: `func NewServerWithTlsSupport(serverConfig ServerConfig, tlsConfig *TLSConfig, iedModel *IedModel) (*IedServer, error)`  
**C Function**: `IedServer_createWithConfig()`

**Description**: Creates server with TLS encryption.

---

### Start

**Go Function**: `func (is *IedServer) Start(port int)`  
**C Function**: `IedServer_start()`

**Description**: Starts the server on specified port (threaded mode).

**Example**:
```go
server.Start(102)
if !server.IsRunning() {
    log.Fatal("Failed to start server")
}
```

---

### StartThreadless

**Go Function**: `func (is *IedServer) StartThreadless(port int)`  
**C Function**: `IedServer_startThreadless()`

**Description**: Starts server in non-threaded mode.

**Example**:
```go
server.StartThreadless(102)
for {
    ready := server.WaitReady(100) // 100ms timeout
    if ready > 0 {
        server.ProcessIncomingData()
    }
    server.PerformPeriodicTasks()
}
```

---

### Stop

**Go Function**: `func (is *IedServer) Stop()`  
**C Function**: `IedServer_stop()`

**Description**: Stops the server.

---

### StopThreadless

**Go Function**: `func (is *IedServer) StopThreadless()`  
**C Function**: `IedServer_stopThreadless()`

**Description**: Stops server in threadless mode.

---

### IsRunning

**Go Function**: `func (is *IedServer) IsRunning() bool`  
**C Function**: `IedServer_isRunning()`

**Description**: Checks if server is running.

---

### WaitReady

**Go Function**: `func (is *IedServer) WaitReady(timeoutMs uint) int`  
**C Function**: `IedServer_waitReady()`

**Description**: Waits for connection data (threadless mode). Returns number of connections ready.

---

### ProcessIncomingData

**Go Function**: `func (is *IedServer) ProcessIncomingData()`  
**C Function**: `IedServer_processIncomingData()`

**Description**: Processes incoming data (threadless mode).

---

### PerformPeriodicTasks

**Go Function**: `func (is *IedServer) PerformPeriodicTasks()`  
**C Function**: `IedServer_performPeriodicTasks()`

**Description**: Runs periodic background tasks (threadless mode).

---

### SetLocalIpAddress

**Go Function**: `func (is *IedServer) SetLocalIpAddress(ipAddress string)`  
**C Function**: `IedServer_setLocalIpAddress()`

**Description**: Sets the local IP address to bind to.

---

### SetFilestoreBasepath

**Go Function**: `func (is *IedServer) SetFilestoreBasepath(basepath string)`  
**C Function**: `IedServer_setFilestoreBasepath()`

**Description**: Sets the base path for file services.

---

### SetConnectionIndicationHandler

**Go Function**: `func (is *IedServer) SetConnectionIndicationHandler(handler ConnectionIndicationHandler)`  
**C Function**: `IedServer_setConnectionIndicationHandler()`

**Description**: Sets handler for client connect/disconnect events.

**Example**:
```go
server.SetConnectionIndicationHandler(func(conn *iec61850.ClientConnection, connected bool) {
    if connected {
        log.Println("Client connected from:", conn.GetPeerAddress())
    } else {
        log.Println("Client disconnected")
    }
})
```

---

### SetClientAuthenticator

**Go Function**: `func (is *IedServer) SetClientAuthenticator(authenticator ClientAuthenticator)`  
**C Function**: `IedServer_setClientAuthenticator()`

**Description**: Sets handler for client authentication.

---

### UpdateAttributeValue

**Go Function**: `func (is *IedServer) UpdateAttributeValue(value *DataAttribute, newValue *MmsValue) error`  
**C Function**: `IedServer_updateAttributeValue()`

**Description**: Updates server-side attribute value.

**Example**:
```go
attr := model.GetModelNodeByObjectReference("Device/LLN0.Beh.stVal", FC_ST)
err := server.UpdateAttributeValue(attr, iec61850.NewMmsValue(iec61850.Integer, int64(1)))
```

---

### LockDataModel

**Go Function**: `func (is *IedServer) LockDataModel()`  
**C Function**: `IedServer_lockDataModel()`

**Description**: Locks data model for thread-safe updates.

---

### UnlockDataModel

**Go Function**: `func (is *IedServer) UnlockDataModel()`  
**C Function**: `IedServer_unlockDataModel()`

**Description**: Unlocks data model.

**Example**:
```go
server.LockDataModel()
defer server.UnlockDataModel()
// Update multiple values safely
```

---

## Control Functions

### Operate

**Go Function**: `func (c *Client) Operate(controlRef string, param ControlObjectParam) error`  
**C Function**: `ControlObjectClient_operate()`

**Description**: Executes a control operation.

**Parameters**:
- `controlRef` - Control object reference
- `param` - Control parameters (value, origin, ctlNum)

**Example**:
```go
param := iec61850.ControlObjectParam{
    CtlVal:  true,
    CtlNum:  1,
    Origin:  iec61850.NewControlOriginator(orCat, orIdent),
    Test:    false,
}
err := client.Operate("Device/XCBR1.Pos", param)
```

---

### Select

**Go Function**: `func (c *Client) Select(controlRef string) error`  
**C Function**: `ControlObjectClient_select()`

**Description**: Selects a control object (for SBO control model).

---

### Cancel

**Go Function**: `func (c *Client) Cancel(controlRef string) error`  
**C Function**: `ControlObjectClient_cancel()`

**Description**: Cancels a selected control operation.

---

## Reporting Functions

### GetRCBValues

**Go Function**: `func (c *Client) GetRCBValues(objectReference string) (*ClientReportControlBlock, error)`  
**C Function**: `ClientReportControlBlock_create()`, getters

**Description**: Retrieves report control block values.

**Example**:
```go
rcb, err := client.GetRCBValues("Device/LLN0.RP.report1")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("RCB enabled: %v\n", rcb.RptEna)
```

---

### SetRCBValues

**Go Function**: `func (c *Client) SetRCBValues(objectReference string, settings ClientReportControlBlock) error`  
**C Function**: Setters + `ClientReportControlBlock_setRCBValues()`

**Description**: Updates report control block settings.

**Example**:
```go
settings := iec61850.ClientReportControlBlock{
    RptEna:  true,
    IntgPd:  5000, // 5 seconds
    OptFlds: iec61850.OptFlds{
        SeqNum:     true,
        TimeStamp:  true,
        DataSet:    true,
    },
}
err := client.SetRCBValues("Device/LLN0.RP.report1", settings)
```

---

### InstallReportHandler

**Go Function**: `func (c *Client) InstallReportHandler(rcbReference string, handler ReportCallbackFunc, handlerParam interface{}) error`  
**C Function**: `IedConnection_installReportHandler()`

**Description**: Installs handler for receiving reports.

**Example**:
```go
handler := func(param interface{}, report *iec61850.ClientReport) {
    fmt.Printf("Report received: %d values\n", len(report.DataSetValues))
    for i, val := range report.DataSetValues {
        fmt.Printf("  [%d]: %v\n", i, val)
    }
}
err := client.InstallReportHandler("Device/LLN0.RP.report1", handler, nil)
```

---

## GOOSE Functions

### NewGooseSubscriber

**Go Function**: `func NewGooseSubscriber(conf SubscriberConf) *GooseSubscriber`  
**C Function**: `GooseSubscriber_create()`

**Description**: Creates a GOOSE subscriber.

**Example**:
```go
conf := iec61850.SubscriberConf{
    InterfaceID: "eth0",
    AppID:       1000,
    GoID:        "GoosePublisher1",
}
subscriber := iec61850.NewGooseSubscriber(conf)
defer subscriber.Destroy()
```

---

### SetGooseReceiver

**Go Function**: `func (subscriber *GooseSubscriber) SetGooseReceiver(receiver func(*GooseSubscriber))`  
**C Function**: `GooseSubscriber_setListener()`

**Description**: Sets callback for received GOOSE messages.

**Example**:
```go
subscriber.SetGooseReceiver(func(sub *iec61850.GooseSubscriber) {
    goID := sub.GetGoID()
    stNum := sub.GetStNum()
    sqNum := sub.GetSqNum()
    fmt.Printf("GOOSE from %s: stNum=%d sqNum=%d\n", goID, stNum, sqNum)
})
```

---

### Subscribe

**Go Function**: `func (subscriber *GooseSubscriber) Subscribe() error`  
**C Function**: `GooseReceiver_addSubscriber()`, `GooseReceiver_start()`

**Description**: Starts receiving GOOSE messages.

---

### NewGoosePublisher

**Go Function**: `func NewGoosePublisher(conf GoosePublisherConf) (*GoosePublisher, error)`  
**C Function**: `GoosePublisher_create()`

**Description**: Creates a GOOSE publisher.

**Example**:
```go
conf := iec61850.GoosePublisherConf{
    InterfaceID: "eth0",
    AppID:       1000,
    GoID:        "MyGoosePublisher",
    GoCbRef:     "Device/LLN0$GO$gcb1",
    DataSetRef:  "Device/LLN0$dataset1",
}
publisher, err := iec61850.NewGoosePublisher(conf)
```

---

### Publish

**Go Function**: `func (publisher *GoosePublisher) Publish(values []*MmsValue) error`  
**C Function**: `GoosePublisher_publish()`

**Description**: Publishes GOOSE message with data values.

---

## Sampled Values (SV) Functions

### NewSVSubscriber

**Go Function**: `func NewSVSubscriber(interfaceID string) *SVSubscriber`  
**C Function**: `SVReceiver_create()`

**Description**: Creates a Sampled Values subscriber.

**Example**:
```go
svSub := iec61850.NewSVSubscriber("eth0")
defer svSub.Destroy()

svSub.SetSVReceiver(func(sub *iec61850.SVSubscriber, appID int, data []byte) {
    fmt.Printf("SV received: AppID=%d, %d bytes\n", appID, len(data))
})

svSub.Subscribe()
```

---

### NewSVPublisher

**Go Function**: `func NewSVPublisher(conf SVPublisherConf) (*SVPublisher, error)`  
**C Function**: `SVPublisher_create()`

**Description**: Creates a Sampled Values publisher.

---

### PublishSV

**Go Function**: `func (publisher *SVPublisher) PublishSV(asdu *SVPublisherASDU) error`  
**C Function**: `SVPublisher_publish()`

**Description**: Publishes sampled values ASDU.

---

## File Services Functions

### GetServerFileDirectory

**Go Function**: `func (c *Client) GetServerFileDirectory(directoryName string) ([]string, error)`  
**C Function**: `IedConnection_getFileDirectory()`

**Description**: Retrieves file directory listing from server.

**Example**:
```go
files, err := client.GetServerFileDirectory("/")
if err != nil {
    log.Fatal(err)
}
for _, file := range files {
    fmt.Println("File:", file)
}
```

---

### GetFileDirectoryEx

**Go Function**: `func (c *Client) GetFileDirectoryEx(directoryName, continueAfter string) ([]FileDirectoryEntry, bool, error)`  
**C Function**: `MmsConnection_fileDirectory()`

**Description**: Gets detailed file directory with metadata.

**Example**:
```go
entries, moreFollows, err := client.GetFileDirectoryEx("/", "")
for _, entry := range entries {
    fmt.Printf("%s: %d bytes (modified: %d)\n",
        entry.FileName, entry.FileSize, entry.LastModified)
}
```

---

### GetFile

**Go Function**: `func (c *Client) GetFile(fileName string) ([]byte, error)`  
**C Function**: `IedConnection_getFile()`

**Description**: Downloads a file from the server.

**Example**:
```go
fileData, err := client.GetFile("/config/settings.cfg")
if err != nil {
    log.Fatal(err)
}
err = os.WriteFile("local_settings.cfg", fileData, 0644)
```

---

### FileOpen

**Go Function**: `func (c *Client) FileOpen(fileName string, openRead bool) (uint32, error)`  
**C Function**: `IedConnection_fileOpen()`

**Description**: Opens a file on the server for reading/writing.

---

### FileRead

**Go Function**: `func (c *Client) FileRead(frsmID uint32, bufferSize int) ([]byte, bool, error)`  
**C Function**: `IedConnection_fileRead()`

**Description**: Reads chunk from open file.

---

### FileClose

**Go Function**: `func (c *Client) FileClose(frsmID uint32) error`  
**C Function**: `IedConnection_fileClose()`

**Description**: Closes an open file.

---

### FileDelete

**Go Function**: `func (c *Client) FileDelete(fileName string) error`  
**C Function**: `IedConnection_fileDelete()`

**Description**: Deletes a file on the server.

---

### ObtainFile

**Go Function**: `func (c *Client) ObtainFile(sourceFile, destFile string) error`  
**C Function**: `MmsConnection_obtainFile()`

**Description**: Requests server to upload a file from client (client→server transfer).

---

### RenameFile

**Go Function**: `func (c *Client) RenameFile(currentName, newName string) error`  
**C Function**: `MmsConnection_fileRename()`

**Description**: Renames a file on the server.

---

## Dataset Functions

### GetDataSetDirectory

**Go Function**: `func (c *Client) GetDataSetDirectory(dataSetRef string) ([]*MmsVariableAccessSpec, error)`  
**C Function**: `IedConnection_getDataSetDirectory()`

**Description**: Retrieves dataset member list.

**Example**:
```go
members, err := client.GetDataSetDirectory("Device/LLN0$dataset1")
for _, member := range members {
    fmt.Printf("Member: %s/%s\n", member.DomainID, member.ItemID)
}
```

---

### CreateDataSet

**Go Function**: `func (c *Client) CreateDataSet(dataSetRef string, dataSetEntries []*MmsVariableAccessSpec) error`  
**C Function**: `IedConnection_createDataSet()`

**Description**: Creates a new dataset on the server.

---

### DeleteDataSet

**Go Function**: `func (c *Client) DeleteDataSet(dataSetRef string) error`  
**C Function**: `IedConnection_deleteDataSet()`

**Description**: Deletes a dataset.

---

## Utility Functions

### NewTimestamp

**Go Function**: `func NewTimestamp(time ...time.Time) *Timestamp`  
**C Function**: `Timestamp_create()`

**Description**: Creates an IEC 61850 timestamp.

**Example**:
```go
ts := iec61850.NewTimestamp(time.Now())
fmt.Printf("Timestamp in ms: %d\n", ts.GetTimeInMs())
```

---

### GetMmsError

**Go Function**: `func GetMmsError(err C.MmsError) error`  
**C Function**: N/A (helper)

**Description**: Converts C MmsError to Go error.

---

### GetIedClientError

**Go Function**: `func GetIedClientError(err C.IedClientError) error`  
**C Function**: N/A (helper)

**Description**: Converts C IedClientError to Go error.

---

### IsBitSet

**Go Function**: `func IsBitSet(val int, pos int) bool`  
**C Function**: N/A (utility)

**Description**: Checks if a bit is set at position.

---

### NewSettings

**Go Function**: `func NewSettings() Settings`  
**C Function**: N/A (helper)

**Description**: Creates default connection settings.

**Example**:
```go
settings := iec61850.NewSettings()
settings.Host = "192.168.1.100"
client, _ := iec61850.NewClient(settings)
```

---

## Type Introspection Functions

### GetVariableAccessAttributes

**Go Function**: `func (c *Client) GetVariableAccessAttributes(domainID, itemID string) (*MmsVariableSpecificationRef, error)`  
**C Function**: `MmsConnection_getVariableAccessAttributes()`

**Description**: Retrieves type specification for an MMS variable.

**Example**:
```go
typeSpec, err := client.GetVariableAccessAttributes("domain", "variable")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Type: %v, Size: %d\n", typeSpec.GetType(), typeSpec.GetSize())
```

---

### GetVariableSpecType

**Go Function**: `func (c *Client) GetVariableSpecType(objectReference string, fc FC) (MmsType, error)`  
**C Function**: `IedConnection_getVariableAccessAttributes()`

**Description**: Gets the MMS type of an IEC 61850 object.

---

*End of Functions Reference*
