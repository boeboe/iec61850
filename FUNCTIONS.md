# IEC 61850 Go Bindings - Functions Reference

**Version**: 1.6.1  
**Generated**: February 7, 2026

This document provides comprehensive documentation for all exported Go functions in the iec61850 package, including their corresponding C functions, descriptions, and usage examples. It is aligned with [GAPS.md](GAPS.md) (MMS coverage analysis).

---

## Table of Contents

1. [Client Connection Functions](#client-connection-functions)
2. [Client Read/Write Functions](#client-readwrite-functions)
3. [Client MMS & Discovery Functions](#client-mms--discovery-functions)
4. [MMS Connection Functions](#mms-connection-functions)
5. [MMS Value Functions](#mms-value-functions)
6. [Server Functions](#server-functions)
7. [Control Functions](#control-functions)
8. [Reporting Functions](#reporting-functions)
9. [GOOSE Functions](#goose-functions)
10. [Sampled Values (SV) Functions](#sampled-values-sv-functions)
11. [File Services Functions](#file-services-functions)
12. [Dataset Functions](#dataset-functions)
13. [Utility Functions](#utility-functions)
14. [Type System (MmsVariableSpecificationRef)](#type-system-mmsvariablespecificationref)

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

### NewClientWithoutConnect

**Go Function**: `func NewClientWithoutConnect(settings Settings) (*Client, error)`  
**C Function**: `IedConnection_create()`, timeout setters (no connect)

**Description**: Creates a client without connecting. Use ConnectAsync or ConnectWithAuth to connect later. Close() must be called when done.

**Example**:
```go
client, _ := iec61850.NewClientWithoutConnect(settings)
defer client.Close()
client.ConnectAsync("192.168.1.10", 102, func(err error) {
    if err == nil { fmt.Println("Connected") }
})
```

---

### NewClientWithoutConnectWithTls

**Go Function**: `func NewClientWithoutConnectWithTls(settings Settings, tlsConfig *TLSConfig) (*Client, error)`  
**C Function**: `IedConnection_createWithTlsSupport()`

**Description**: Like NewClientWithoutConnect but with TLS configuration.

**Example**:
```go
tlsCfg := &iec61850.TLSConfig{ChainValidation: true, OwnCertificate: cert, OwnKey: key, CACerts: [][]byte{ca}}
client, _ := iec61850.NewClientWithoutConnectWithTls(settings, tlsCfg)
defer client.Close()
err := client.ConnectWithAuth("192.168.1.10", 102, "user", "secret")
```

---

### ConnectWithAuth

**Go Function**: `func (c *Client) ConnectWithAuth(hostname string, port int, username, password string) error`  
**C Function**: `AcseAuthenticationParameter_*`, `IedConnection_connect()`

**Description**: Connects using ACSE password authentication. Client must have been created with NewClientWithoutConnect or NewClientWithoutConnectWithTls. Username is not used by ACSE password; only password is sent.

**Example**:
```go
client, _ := iec61850.NewClientWithoutConnect(settings)
defer client.Close()
err := client.ConnectWithAuth("192.168.1.10", 102, "", "mypassword")
if err != nil {
    log.Fatal(err)
}
```

---

### ConnectAsync

**Go Function**: `func (c *Client) ConnectAsync(hostname string, port int, callback func(error))`  
**C Function**: `IedConnection_connectAsync()`

**Description**: Starts non-blocking connection. Callback is invoked with nil when connected or with an error on failure/close. Client must have been created with NewClientWithoutConnect or NewClientWithoutConnectWithTls.

**Example**:
```go
client, _ := iec61850.NewClientWithoutConnect(settings)
defer client.Close()
client.ConnectAsync("192.168.1.10", 102, func(err error) {
    if err != nil {
        log.Println("Connection failed:", err)
        return
    }
    log.Println("Connected")
})
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

### Disconnect

**Go Function**: `func (c *Client) Disconnect()`  
**C Function**: `MmsConnection_close()`

**Description**: Closes the MMS association without destroying the connection object. Use Close() to release the client fully.

**Example**:
```go
client.Disconnect()
```

---

### Abort

**Go Function**: `func (c *Client) Abort() error`  
**C Function**: `MmsConnection_abort()`

**Description**: Aborts the MMS connection immediately.

**Example**:
```go
if err := client.Abort(); err != nil {
    log.Println("Abort failed:", err)
}
```

---

### SetConnectionLostHandler

**Go Function**: `func (c *Client) SetConnectionLostHandler(handler func(error))`  
**C Function**: `MmsConnection_setConnectionLostHandler()`

**Description**: Sets callback invoked when the connection is lost. Call before Connect.

**Example**:
```go
client.SetConnectionLostHandler(func(err error) {
    log.Println("Connection lost:", err)
})
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

**Example**:
```go
err := client.GetLastApplError()
if err.Code != 0 {
    fmt.Printf("Last control error: %s\n", err.Description)
}
```

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

**Example**:
```go
value, err := client.Read("Device/GGIO1.AnIn1.mag.f", iec61850.MX)
if err == nil {
    fmt.Printf("Value: %v\n", value)
}
```

---

### ReadMultiple

**Go Function**: `func (c *Client) ReadMultiple(objectRefs []string, fc FC) ([]interface{}, error)`  
**C Function**: `MmsConnection_readMultipleVariables()`

**Description**: Reads multiple data attributes in one request. Returns values in same order as objectRefs.

**Example**:
```go
refs := []string{"Device/GGIO1.AnIn1.mag.f", "Device/GGIO1.AnIn2.mag.f"}
vals, err := client.ReadMultiple(refs, iec61850.MX)
```

**Parameters**:
- `objectRef` - Object reference (e.g., "simpleIOGenericIO/GGIO1.AnIn1.mag.f")
- `fc` - Functional constraint (ST, MX, SP, etc.)

**Returns**: Value as interface{} and error

**Example**:
```go
value, err := client.Read("Device/GGIO1.AnIn1.mag.f", iec61850.MX)
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

**Example**:
```go
val, err := client.ReadInt64("Device/LLN0.Beh.t", iec61850.ST)
```

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

**Example**:
```go
val, err := client.ReadFloat64("Device/GGIO1.AnIn1.mag.f", iec61850.MX)
```

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

**Example**:
```go
bits, err := client.ReadBitString("Device/GGIO1.Ind1.stVal", iec61850.ST)
```

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

## Client MMS & Discovery Functions

Functions in this section use the MMS layer (via IedConnection). See [GAPS.md](GAPS.md) for full C↔Go mapping.

### GetConnectionParameters

**Go Function**: `func (c *Client) GetConnectionParameters() (*MmsConnectionParameters, error)`  
**C Function**: `MmsConnection_getMmsConnectionParameters()`

**Description**: Returns negotiated MMS connection parameters (max PDU size, outstanding calls, etc.). Call after connection is established.

**Example**:
```go
params, err := client.GetConnectionParameters()
if err == nil {
    fmt.Printf("Max PDU size: %d\n", params.MaxPduSize)
}
```

---

### Identify

**Go Function**: `func (c *Client) Identify() (*MmsServerIdentity, error)`  
**C Function**: `MmsConnection_identify()`

**Description**: Returns server identification (vendor, model, revision).

**Example**:
```go
id, err := client.Identify()
if err == nil {
    fmt.Printf("%s %s (Rev %s)\n", id.VendorName, id.ModelName, id.Revision)
}
```

---

### GetServerStatus

**Go Function**: `func (c *Client) GetServerStatus(extendedDerivation bool) (*MmsServerStatus, error)`  
**C Function**: `MmsConnection_getServerStatus()`

**Description**: Returns MMS server status (VMD logical/physical status, local detail).

**Example**:
```go
status, err := client.GetServerStatus(false)
if err == nil {
    fmt.Printf("VMD logical: %d, physical: %d\n", status.VmdLogicalStatus, status.VmdPhysicalStatus)
}
```

---

### GetDomainNames

**Go Function**: `func (c *Client) GetDomainNames() ([]string, error)`  
**C Function**: `MmsConnection_getDomainNames()`

**Description**: Returns list of MMS domain names on the server.

**Example**:
```go
domains, err := client.GetDomainNames()
for _, d := range domains {
    fmt.Println("Domain:", d)
}
```

---

### GetDomainVariableNames

**Go Function**: `func (c *Client) GetDomainVariableNames(domainID string) ([]string, error)`  
**C Function**: `MmsConnection_getDomainVariableNames()`

**Description**: Returns variable names in the given domain.

**Example**:
```go
names, err := client.GetDomainVariableNames("MYDOMAIN")
```

---

### GetDomainVariableListNames

**Go Function**: `func (c *Client) GetDomainVariableListNames(domainID string) ([]string, error)`  
**C Function**: `MmsConnection_getDomainVariableListNames()`

**Description**: Returns named variable list names in the domain.

**Example**:
```go
listNames, err := client.GetDomainVariableListNames("MYDOMAIN")
```

---

### GetVariableListNamesAssociationSpecific

**Go Function**: `func (c *Client) GetVariableListNamesAssociationSpecific() ([]string, error)`  
**C Function**: `MmsConnection_getVariableListNamesAssociationSpecific()`

**Description**: Returns association-specific variable list names.

**Example**:
```go
names, err := client.GetVariableListNamesAssociationSpecific()
```

---

### WriteMultipleVariables

**Go Function**: `func (c *Client) WriteMultipleVariables(domainID string, itemIDs []string, values []*MmsValueRef) ([]MmsDataAccessError, error)`  
**C Function**: `MmsConnection_writeMultipleVariables()`

**Description**: Writes multiple MMS variables in one request. Returns per-item access results.

**Example**:
```go
items := []string{"var1", "var2"}
vals := []*iec61850.MmsValueRef{ref1, ref2}
results, err := client.WriteMultipleVariables("domain", items, vals)
```

---

### WriteNamedVariableList

**Go Function**: `func (c *Client) WriteNamedVariableList(domainID, listName string, values []*MmsValueRef) ([]MmsDataAccessError, error)`  
**C Function**: `MmsConnection_writeNamedVariableList()`

**Description**: Writes values to a named variable list. Returns per-element access results.

**Example**:
```go
results, err := client.WriteNamedVariableList("domain", "list1", values)
```

---

### ReadNamedVariableListValues

**Go Function**: `func (c *Client) ReadNamedVariableListValues(domainID, listName string, specWithResult bool) ([]*MmsValue, error)`  
**C Function**: `MmsConnection_readNamedVariableListValues()`

**Description**: Reads all values from a named variable list.

**Example**:
```go
vals, err := client.ReadNamedVariableListValues("domain", "list1", true)
```

---

### ReadNamedVariableListDirectory

**Go Function**: `func (c *Client) ReadNamedVariableListDirectory(domainID, listName string) (entries []VariableListEntry, deletable bool, err error)`  
**C Function**: `MmsConnection_readNamedVariableListDirectory()`

**Description**: Returns directory (variable specs) for a named variable list.

**Example**:
```go
entries, deletable, err := client.ReadNamedVariableListDirectory("domain", "list1")
```

---

### ReadDataSetValues

**Go Function**: `func (c *Client) ReadDataSetValues(dataSetReference string) (*ClientDataSet, error)`  
**C Function**: `IedConnection_readDataSetValues()`

**Description**: Reads data set values from the server. dataSetReference is the object reference (e.g. "LD/LN.dsName" or "@asName"). The returned ClientDataSet can be used with NewGooseSubscriberWithDataSet via GooseDataSetValues(); keep the ClientDataSet alive for the lifetime of that subscriber.

**Example**:
```go
dataSet, err := client.ReadDataSetValues("simpleIOGenericIO/LLN0$dataset1")
if err != nil { log.Fatal(err) }
defer dataSet.Destroy()
gooseVals := dataSet.GooseDataSetValues()
sub := iec61850.NewGooseSubscriberWithDataSet(conf, &gooseVals)
```

---

### ClientDataSet.Destroy

**Go Function**: `func (d *ClientDataSet) Destroy()`  
**C Function**: `ClientDataSet_destroy()`

**Description**: Frees the ClientDataSet. Do not use it or any GooseDataSetValues derived from it after Destroy.

---

### ClientDataSet.GooseDataSetValues

**Go Function**: `func (d *ClientDataSet) GooseDataSetValues() GooseDataSetValues`  
**C Function**: `ClientDataSet_getValues()`

**Description**: Returns a handle to the underlying MmsValue (MMS_ARRAY) for use with NewGooseSubscriberWithDataSet. The ClientDataSet must remain alive while the subscriber uses it.

---

### DefineNamedVariableList

**Go Function**: `func (c *Client) DefineNamedVariableList(domainID, listName string, variableSpecs []VariableAccessSpec) error`  
**C Function**: `MmsConnection_defineNamedVariableList()`

**Description**: Creates a domain-specific named variable list.

**Example**:
```go
specs := []iec61850.VariableAccessSpec{{DomainID: "domain", ItemID: "var1"}}
err := client.DefineNamedVariableList("domain", "mylist", specs)
```

---

### DefineNamedVariableListAssociationSpecific

**Go Function**: `func (c *Client) DefineNamedVariableListAssociationSpecific(listName string, variableSpecs []VariableAccessSpec) error`  
**C Function**: `MmsConnection_defineNamedVariableListAssociationSpecific()`

**Description**: Creates an association-specific named variable list.

**Example**:
```go
err := client.DefineNamedVariableListAssociationSpecific("mylist", specs)
```

---

### DeleteNamedVariableList

**Go Function**: `func (c *Client) DeleteNamedVariableList(domainID, listName string) (bool, error)`  
**C Function**: `MmsConnection_deleteNamedVariableList()`

**Description**: Deletes a named variable list. Returns true if deleted.

**Example**:
```go
deleted, err := client.DeleteNamedVariableList("domain", "list1")
```

---

### DeleteAssociationSpecificNamedVariableList

**Go Function**: `func (c *Client) DeleteAssociationSpecificNamedVariableList(listName string) (bool, error)`  
**C Function**: `MmsConnection_deleteAssociationSpecificNamedVariableList()`

**Description**: Deletes an association-specific named variable list.

**Example**:
```go
deleted, err := client.DeleteAssociationSpecificNamedVariableList("list1")
```

---

### ReadJournalTimeRange

**Go Function**: `func (c *Client) ReadJournalTimeRange(domainID, itemID string, startTimeMs, endTimeMs uint64) (entries []JournalEntry, moreFollows bool, err error)`  
**C Function**: `MmsConnection_readJournalTimeRange()`

**Description**: Reads journal entries within a time range.

**Example**:
```go
entries, more, err := client.ReadJournalTimeRange("domain", "journal", startMs, endMs)
```

---

### ReadJournalStartAfter

**Go Function**: `func (c *Client) ReadJournalStartAfter(domainID, itemID string, timeSpecificationMs uint64, entrySpecification []byte) (entries []JournalEntry, moreFollows bool, err error)`  
**C Function**: `MmsConnection_readJournalStartAfter()`

**Description**: Reads journal entries after a given time/entry.

**Example**:
```go
entries, more, err := client.ReadJournalStartAfter("domain", "journal", timeMs, entryID)
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

**Example**:
```go
mmsConn.Destroy()
```

---

### SetConnectTimeout

**Go Function**: `func (c *MmsConnection) SetConnectTimeout(timeoutMs uint32)`  
**C Function**: `MmsConnection_setConnectTimeout()`

**Description**: Sets connection timeout in milliseconds.

**Example**:
```go
mmsConn.SetConnectTimeout(5000)
```

---

### SetRequestTimeout

**Go Function**: `func (c *MmsConnection) SetRequestTimeout(timeoutMs uint32)`  
**C Function**: `MmsConnection_setRequestTimeout()`

**Description**: Sets request timeout in milliseconds.

**Example**:
```go
mmsConn.SetRequestTimeout(3000)
```

---

### GetRequestTimeout

**Go Function**: `func (c *MmsConnection) GetRequestTimeout() uint32`  
**C Function**: `MmsConnection_getRequestTimeout()`

**Description**: Gets current request timeout.

**Example**:
```go
timeout := mmsConn.GetRequestTimeout()
```

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

**Example**:
```go
_ = mmsConn.Conclude()
```

---

### ConcludeAsync

**Go Function**: `func (c *MmsConnection) ConcludeAsync(callback func(error)) error`  
**C Function**: `MmsConnection_concludeAsync()`

**Description**: Asynchronous version of Conclude.

**Example**:
```go
_ = mmsConn.ConcludeAsync(func(err error) { log.Println("Concluded:", err) })
```

---

### AbortAsync

**Go Function**: `func (c *MmsConnection) AbortAsync() error`  
**C Function**: `MmsConnection_abortAsync()`

**Description**: Aborts the MMS connection asynchronously.

**Example**:
```go
_ = mmsConn.AbortAsync()
```

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

**Example**:
```go
_ = mmsConn.WriteVariableAsync("domain", "var", valueRef, func(err error) {})
```

---

### GetDomainNamesAsync

**Go Function**: `func (c *MmsConnection) GetDomainNamesAsync(callback func([]string, error)) error`  
**C Function**: `MmsConnection_getDomainNamesAsync()`

**Description**: Asynchronously retrieves all domain names from the server.

**Example**:
```go
_ = mmsConn.GetDomainNamesAsync("", func(names []string, more bool, err error) {})
```

---

### GetDomainVariableNamesAsync

**Go Function**: `func (c *MmsConnection) GetDomainVariableNamesAsync(domainID string, callback func([]string, error)) error`  
**C Function**: `MmsConnection_getDomainVariableNamesAsync()`

**Description**: Asynchronously retrieves variable names in a domain.

**Example**:
```go
_ = mmsConn.GetDomainVariableNamesAsync("domain", "", func(names []string, more bool, err error) {})
```

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

**Example**:
```go
iso := mmsConn.GetIsoConnectionParameters()
fmt.Printf("Local AP title: %x\n", iso.LocalApTitle)
```

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

### SetLocalDetail / GetLocalDetail

**Go Function**: `func (c *MmsConnection) SetLocalDetail(localDetail int32)` / `func (c *MmsConnection) GetLocalDetail() int32`  
**C Function**: `MmsConnection_setLocalDetail()` / `MmsConnection_getLocalDetail()`

**Description**: Sets or gets the local detail (max MMS PDU size negotiation).

**Example**:
```go
mmsConn.SetLocalDetail(8192)
detail := mmsConn.GetLocalDetail()
```

---

### SetMaxOutstandingCalls

**Go Function**: `func (c *MmsConnection) SetMaxOutstandingCalls(calling, called int)`  
**C Function**: `MmsConnection_setMaxOutstandingCalls()`

**Description**: Sets maximum outstanding calling/called requests.

**Example**:
```go
mmsConn.SetMaxOutstandingCalls(5, 5)
```

---

### GetConnectTimeout

**Go Function**: `func (c *MmsConnection) GetConnectTimeout() uint32`  
**C Function**: `MmsConnection_getConnectTimeout()`

**Description**: Returns connection timeout in milliseconds.

**Example**:
```go
ms := mmsConn.GetConnectTimeout()
```

---

### SetIsoConnectionParameters

**Go Function**: `func (c *MmsConnection) SetIsoConnectionParameters(params *IsoConnectionParameters)`  
**C Function**: `MmsConnection_setIsoConnectionParameters()`

**Description**: Sets ISO layer parameters (AP title, selectors). Call before Connect.

**Example**:
```go
mmsConn.SetIsoConnectionParameters(&iec61850.IsoConnectionParameters{})
```

---

### SetFilestoreBasepath

**Go Function**: `func (c *MmsConnection) SetFilestoreBasepath(basepath string)`  
**C Function**: `MmsConnection_setFilestoreBasepath()`

**Description**: Sets client-side filestore base path for file operations.

**Example**:
```go
mmsConn.SetFilestoreBasepath("/tmp/mms")
```

---

### SetInformationReportHandler

**Go Function**: `func (c *MmsConnection) SetInformationReportHandler(callback func(domainName, variableListName string, value *MmsValue, isVariableListName bool))`  
**C Function**: `MmsConnection_setInformationReportHandler()`

**Description**: Sets callback for unsolicited information reports (e.g. MMS data set updates).

**Example**:
```go
mmsConn.SetInformationReportHandler(func(domain, list string, val *iec61850.MmsValue, isList bool) {})
```

---

### ReadArrayElements

**Go Function**: `func (c *MmsConnection) ReadArrayElements(domainID, itemID string, startIndex, numberOfElements uint32) (*MmsValue, error)`  
**C Function**: `MmsConnection_readArrayElements()`

**Description**: Reads a range of array elements from an MMS variable.

**Example**:
```go
val, err := mmsConn.ReadArrayElements("domain", "arrayVar", 0, 10)
```

---

### WriteArrayElements

**Go Function**: `func (c *MmsConnection) WriteArrayElements(domainID, itemID string, index, numberOfElements int, value *MmsValueRef) (MmsDataAccessError, error)`  
**C Function**: `MmsConnection_writeArrayElements()`

**Description**: Writes a value to array elements.

**Example**:
```go
errCode, err := mmsConn.WriteArrayElements("domain", "arr", 0, 5, valueRef)
```

---

### ReadNamedVariableListValues

**Go Function**: `func (c *MmsConnection) ReadNamedVariableListValues(domainID, listName string, specification bool) (*MmsValue, error)`  
**C Function**: `MmsConnection_readNamedVariableListValues()`

**Description**: Reads all values from a named variable list (returns structure/array MmsValue).

**Example**:
```go
val, err := mmsConn.ReadNamedVariableListValues("domain", "list1", true)
```

---

### WriteMultipleVariables / WriteNamedVariableList

**Go Function**: `func (c *MmsConnection) WriteMultipleVariables(domainID string, items []string, values []*MmsValueRef, accessResults *[]MmsDataAccessError) error`  
**C Function**: `MmsConnection_writeMultipleVariables()`

**Go Function**: `func (c *MmsConnection) WriteNamedVariableList(domainID, listName string, values []*MmsValueRef, accessResults *[]MmsDataAccessError) error`  
**C Function**: `MmsConnection_writeNamedVariableList()`

**Description**: Batch write variables or a named variable list; accessResults receives per-item results.

**Example**:
```go
var results []iec61850.MmsDataAccessError
_ = mmsConn.WriteNamedVariableList("domain", "list1", values, &results)
```

---

### GetNamedVariableListAttributes / GetNamedVariableListAttributesAsync

**Go Function**: `func (c *MmsConnection) GetNamedVariableListAttributes(domainID, listName string) (*MmsNamedVariableListAttributes, error)`  
**C Function**: `MmsConnection_getNamedVariableListAttributes()`

**Description**: Returns list attributes (deletable, list type, variable specs). Async variant takes a callback.

**Example**:
```go
attrs, err := mmsConn.GetNamedVariableListAttributes("domain", "list1")
```

---

### GetDomainVariableListNames / GetDomainJournals / GetVMDVariableNames

**Go Function**: `func (c *MmsConnection) GetDomainVariableListNames(domainID string) ([]string, error)`  
**Go Function**: `func (c *MmsConnection) GetDomainJournals(domainID string) ([]string, error)`  
**Go Function**: `func (c *MmsConnection) GetVMDVariableNames() ([]string, error)`

**C Function**: `MmsConnection_getDomainVariableListNames()` / `getDomainJournals()` / `getVMDVariableNames()`

**Description**: Discovery: variable list names in domain, journal names in domain, VMD-level variable names.

**Example**:
```go
lists, _ := mmsConn.GetDomainVariableListNames("domain")
journals, _ := mmsConn.GetDomainJournals("domain")
vmdVars, _ := mmsConn.GetVMDVariableNames()
```

---

### GetServerStatus

**Go Function**: `func (c *MmsConnection) GetServerStatus(extendedDerivation bool) (*MmsServerStatus, error)`  
**C Function**: `MmsConnection_getServerStatus()`

**Description**: Returns MMS server status (VMD logical/physical status).

**Example**:
```go
status, err := mmsConn.GetServerStatus(false)
```

---

### SendRawData

**Go Function**: `func (c *MmsConnection) SendRawData(buffer []byte) error`  
**C Function**: `MmsConnection_sendRawData()`

**Description**: Sends raw MMS data (for custom or test use).

**Example**:
```go
_ = mmsConn.SendRawData([]byte{0x00, 0x01})
```

---

### FileDirectoryAsync

**Go Function**: `func (c *MmsConnection) FileDirectoryAsync(fileSpecification, continueAfter string, callback func(entries []MmsFileDirectoryEntryEx, moreFollows bool, err error)) error`  
**C Function**: `MmsConnection_fileDirectoryAsync()`

**Description**: Asynchronous file directory; callback receives entries and moreFollows.

**Example**:
```go
_ = mmsConn.FileDirectoryAsync("/", "", func(entries []iec61850.MmsFileDirectoryEntryEx, more bool, err error) {})
```

---

### ReadJournalTimeRange / ReadJournalStartAfter

**Go Function**: `func (c *MmsConnection) ReadJournalTimeRange(domainID, journalName string, startTime, endTime uint64) ([]*MmsJournalEntry, bool, error)`  
**Go Function**: `func (c *MmsConnection) ReadJournalStartAfter(domainID, journalName string, entryID []byte, timeSpec *uint64) ([]*MmsJournalEntry, bool, error)`

**C Function**: `MmsConnection_readJournalTimeRange()` / `readJournalStartAfter()`

**Description**: Read journal entries by time range or after a given entry/time. Async variants available.

**Example**:
```go
entries, more, err := mmsConn.ReadJournalTimeRange("domain", "journal", startMs, endMs)
// or: entries, more, err := mmsConn.ReadJournalStartAfter("domain", "journal", entryID, &timeMs)
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

**Example**:
```go
strVal := iec61850.NewMmsValueVisibleString("hello")
```

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

**Example**:
```go
typeSpec, _ := client.GetVariableAccessAttributes("domain", "arrayVar")
arr := iec61850.MmsValueCreateArray(typeSpec, 5)
```

---

### GetType

**Go Function**: `func (r *MmsValueRef) GetType() MmsType`  
**C Function**: `MmsValue_getType()`

**Description**: Returns the MMS type of the value.

**Example**:
```go
t := mmsVal.GetType()
if t == iec61850.Integer { ... }
```

---

### ToInt64

**Go Function**: `func (r *MmsValueRef) ToInt64() int64`  
**C Function**: `MmsValue_toInt64()`

**Description**: Converts value to int64.

**Example**:
```go
i := mmsVal.ToInt64()
```

---

### ToUint32

**Go Function**: `func (r *MmsValueRef) ToUint32() uint32`  
**C Function**: `MmsValue_toUint32()`

**Description**: Converts value to uint32.

**Example**:
```go
u := mmsVal.ToUint32()
```

---

### ToDouble

**Go Function**: `func (r *MmsValueRef) ToDouble() float64`  
**C Function**: `MmsValue_toDouble()`

**Description**: Converts value to float64.

**Example**:
```go
f := mmsVal.ToDouble()
```

---

### GetBitStringAsInteger

**Go Function**: `func (r *MmsValueRef) GetBitStringAsInteger() uint32`  
**C Function**: `MmsValue_getBitStringAsInteger()`

**Description**: Gets bitstring value as integer (little-endian).

**Example**:
```go
u := mmsVal.GetBitStringAsInteger()
```

---

### GetBitStringAsIntegerBigEndian

**Go Function**: `func (r *MmsValueRef) GetBitStringAsIntegerBigEndian() uint32`  
**C Function**: `MmsValue_getBitStringAsIntegerBigEndian()`

**Description**: Gets bitstring value as integer (big-endian).

**Example**:
```go
u := mmsVal.GetBitStringAsIntegerBigEndian()
```

---

### SetBitStringFromInteger

**Go Function**: `func (r *MmsValueRef) SetBitStringFromInteger(val uint32)`  
**C Function**: `MmsValue_setBitStringFromInteger()`

**Description**: Sets bitstring from integer value (little-endian).

**Example**:
```go
mmsVal.SetBitStringFromInteger(0xFF)
```

---

### SetBitStringFromIntegerBigEndian

**Go Function**: `func (r *MmsValueRef) SetBitStringFromIntegerBigEndian(val uint32)`  
**C Function**: `MmsValue_setBitStringFromIntegerBigEndian()`

**Description**: Sets bitstring from integer value (big-endian).

**Example**:
```go
mmsVal.SetBitStringFromIntegerBigEndian(0xFF)
```

---

### GetBitStringSize / GetNumberOfSetBits

**Go Function**: `func (r *MmsValueRef) GetBitStringSize() int`  
**Go Function**: `func (r *MmsValueRef) GetNumberOfSetBits() int`  
**C Function**: `MmsValue_getBitStringSize()` / bit count

**Description**: Returns bitstring size in bits; GetNumberOfSetBits returns count of set bits.

**Example**:
```go
size := mmsVal.GetBitStringSize()
setBits := mmsVal.GetNumberOfSetBits()
```

---

### GetArraySize

**Go Function**: `func (r *MmsValueRef) GetArraySize() int`  
**C Function**: `MmsValue_getArraySize()`

**Description**: Returns number of elements in array or structure.

**Example**:
```go
n := mmsVal.GetArraySize()
```

---

### GetSizeInMemory / Free

**Go Function**: `func (r *MmsValueRef) GetSizeInMemory() int`  
**Go Function**: `func (r *MmsValueRef) Free()`  
**C Function**: `MmsValue_getSizeInMemory()` / `MmsValue_delete()`

**Description**: Returns size in memory; Free releases the C MmsValue (call when you own it).

**Example**:
```go
bytes := mmsVal.GetSizeInMemory()
mmsVal.Free() // when you own the value
```

---

### SetVisibleString / SetMmsString / SetBinaryTime

**Go Function**: `func (r *MmsValueRef) SetVisibleString(s string)`  
**Go Function**: `func (r *MmsValueRef) SetMmsString(s string)`  
**Go Function**: `func (r *MmsValueRef) SetBinaryTime(ms uint64)`  
**C Function**: `MmsValue_setVisibleString()` / `setMmsString()` / `setBinaryTime()`

**Description**: Set string or binary time value in place.

**Example**:
```go
mmsVal.SetVisibleString("hello")
mmsVal.SetBinaryTime(uint64(time.Now().UnixMilli()))
```

---

### GetBinaryTimeAsUtcMs

**Go Function**: `func (r *MmsValueRef) GetBinaryTimeAsUtcMs() uint64`  
**C Function**: `MmsValue_getBinaryTimeAsUtcMs()`

**Description**: Returns binary time as milliseconds since epoch.

**Example**:
```go
ms := mmsVal.GetBinaryTimeAsUtcMs()
```

---

### MmsValueCreateEmptyArray / MmsValueNewDefaultValue

**Go Function**: `func MmsValueCreateEmptyArray(size int) *MmsValueRef`  
**Go Function**: `func MmsValueNewDefaultValue(typeSpec *MmsVariableSpecificationRef) *MmsValueRef`  
**C Function**: `MmsValue_createEmptyArray()` / `MmsValue_newDefaultValue()`

**Description**: Creates empty array or default value from type specification.

**Example**:
```go
arr := iec61850.MmsValueCreateEmptyArray(10)
defVal := iec61850.MmsValueNewDefaultValue(typeSpec)
```

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

**Example**:
```go
arr.SetElement(0, elemVal)
```

---

### GetDataAccessError

**Go Function**: `func (r *MmsValueRef) GetDataAccessError() MmsDataAccessError`  
**C Function**: `MmsValue_getDataAccessError()`

**Description**: Gets data access error code from value.

**Example**:
```go
errCode := mmsVal.GetDataAccessError()
```

---

### EncodeMmsData

**Go Function**: `func (r *MmsValueRef) EncodeMmsData(buffer []byte, startPos int, encode bool) int`  
**C Function**: `MmsValue_encodeMmsData()`

**Description**: Encodes MmsValue to binary buffer.

**Example**:
```go
buf := make([]byte, 256)
n := mmsVal.EncodeMmsData(buf, 0, true)
```

---

### DecodeMmsData

**Go Function**: `func DecodeMmsData(buffer []byte, startPos, length int) (value *MmsValueRef, endPos int)`  
**C Function**: `MmsValue_decodeMmsData()`

**Description**: Decodes MmsValue from binary buffer.

**Example**:
```go
val, endPos := iec61850.DecodeMmsData(buf, 0, len(buf))
```

---

### CMmsValueToMmsValue

**Go Function**: `func CMmsValueToMmsValue(cVal *C.MmsValue) *MmsValue`  
**C Function**: N/A (Go helper)

**Description**: Converts a C MmsValue (MmsValueRef) to the high-level Go MmsValue. Caller retains ownership of the C value (e.g. must call MmsValue_delete if the C layer transferred it).

**Example**:
```go
// cVal is *C.MmsValue from e.g. MmsConnection ReadVariable
goVal := iec61850.CMmsValueToMmsValue(cVal)
if goVal != nil {
    fmt.Printf("Type: %v Value: %v\n", goVal.Type, goVal.Value)
}
```

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

**Example**:
```go
server, err := iec61850.NewServerWithTlsSupport(config, tlsConfig, model)
```

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

**Example**:
```go
server.Stop()
```

---

### StopThreadless

**Go Function**: `func (is *IedServer) StopThreadless()`  
**C Function**: `IedServer_stopThreadless()`

**Description**: Stops server in threadless mode.

**Example**:
```go
server.StopThreadless()
```

---

### IsRunning

**Go Function**: `func (is *IedServer) IsRunning() bool`  
**C Function**: `IedServer_isRunning()`

**Description**: Checks if server is running.

**Example**:
```go
if server.IsRunning() {
    fmt.Println("Server is up")
}
```

---

### WaitReady

**Go Function**: `func (is *IedServer) WaitReady(timeoutMs uint) int`  
**C Function**: `IedServer_waitReady()`

**Description**: Waits for connection data (threadless mode). Returns number of connections ready.

**Example**:
```go
n := server.WaitReady(100)
```

---

### ProcessIncomingData

**Go Function**: `func (is *IedServer) ProcessIncomingData()`  
**C Function**: `IedServer_processIncomingData()`

**Description**: Processes incoming data (threadless mode).

**Example**:
```go
server.ProcessIncomingData()
```

---

### PerformPeriodicTasks

**Go Function**: `func (is *IedServer) PerformPeriodicTasks()`  
**C Function**: `IedServer_performPeriodicTasks()`

**Description**: Runs periodic background tasks (threadless mode).

**Example**:
```go
server.PerformPeriodicTasks()
```

---

### SetLocalIpAddress

**Go Function**: `func (is *IedServer) SetLocalIpAddress(ipAddress string)`  
**C Function**: `IedServer_setLocalIpAddress()`

**Description**: Sets the local IP address to bind to.

**Example**:
```go
server.SetLocalIpAddress("0.0.0.0")
```

---

### EnableGoosePublishing

**Go Function**: `func (is *IedServer) EnableGoosePublishing()`  
**C Function**: `IedServer_enableGoosePublishing()`

**Description**: Enables GOOSE publishing on the server when using the integrated GOOSE publisher (see ServerConfig.UseIntegratedGoosePublisher).

**Example**:
```go
server.SetGooseInterfaceId("eth0")
server.EnableGoosePublishing()
server.Start(102)
```

---

### DisableGoosePublishing

**Go Function**: `func (is *IedServer) DisableGoosePublishing()`  
**C Function**: `IedServer_disableGoosePublishing()`

**Description**: Disables GOOSE publishing on the server.

**Example**:
```go
server.DisableGoosePublishing()
```

---

### SetGooseInterfaceId

**Go Function**: `func (is *IedServer) SetGooseInterfaceId(interfaceId string)`  
**C Function**: `IedServer_setGooseInterfaceId()`

**Description**: Sets the Ethernet interface used for GOOSE (e.g. "eth0"). May be called before or after Start.

**Example**:
```go
server.SetGooseInterfaceId("eth0")
```

---

### SetFilestoreBasepath

**Go Function**: `func (is *IedServer) SetFilestoreBasepath(basepath string)`  
**C Function**: `IedServer_setFilestoreBasepath()`

**Description**: Sets the (virtual) filestore base path for MMS file services. Call before Start.

**Example**:
```go
server.SetFilestoreBasepath("/var/lib/ied/filestore")
server.Start(102)
```

---

### GetFilestoreBasepath

**Go Function**: `func (is *IedServer) GetFilestoreBasepath() string`  
**C Function**: N/A (value stored by Go binding)

**Description**: Returns the filestore base path last set with SetFilestoreBasepath, or empty string.

**Example**:
```go
server.SetFilestoreBasepath("/var/lib/ied/filestore")
server.Start(102)
base := server.GetFilestoreBasepath()
fmt.Println("Filestore:", base)
```

---

### SetFileAccessHandler

**Go Function**: `func (is *IedServer) SetFileAccessHandler(handler FileAccessHandler)`  
**C Function**: `MmsServer_installFileAccessHandler()`

**Description**: Installs a callback invoked when a client requests an MMS file service. Return nil to allow, or an error (e.g. AccessDenied) to deny.

**Example**:
```go
server.SetFileAccessHandler(func(service iec61850.MmsFileServiceType, local, other string) error {
    if service == iec61850.MmsFileAccessDelete {
        return iec61850.AccessDenied
    }
    return nil
})
```

---

### InstallVariableListAccessHandler

**Go Function**: `func (is *IedServer) InstallVariableListAccessHandler(handler VariableListAccessHandler)`  
**C Function**: `MmsServer_installVariableListAccessHandler()`

**Description**: Installs a callback when a client accesses a named variable list (create, delete, read, write, get directory). Return nil to allow.

**Example**:
```go
server.InstallVariableListAccessHandler(func(accessType iec61850.MmsVariableListAccessType, listType iec61850.MmsVariableListType, domainID, listName string) error {
    if accessType == iec61850.MmsVarlistDelete {
        return iec61850.AccessDenied
    }
    return nil
})
```

---

### InstallReadJournalHandler

**Go Function**: `func (is *IedServer) InstallReadJournalHandler(handler ReadJournalHandler)`  
**C Function**: `MmsServer_installReadJournalHandler()` (via shim)

**Description**: Installs a callback when a client accesses a journal. Return true to allow, false to deny.

**Example**:
```go
server.InstallReadJournalHandler(func(domainID, logName string) bool {
    return logName != "restricted-journal"
})
```

---

### InstallGetNameListHandler

**Go Function**: `func (is *IedServer) InstallGetNameListHandler(handler GetNameListHandler)`  
**C Function**: `MmsServer_installGetNameListHandler()` (via shim)

**Description**: Installs a callback when a client requests a name list (domains, journals, data sets, or data). Return true to allow.

**Example**:
```go
server.InstallGetNameListHandler(func(nameListType iec61850.MmsGetNameListType, domainID string) bool {
    return true
})
```

---

### InstallObtainFileHandler

**Go Function**: `func (is *IedServer) InstallObtainFileHandler(handler ObtainFileHandler)`  
**C Function**: `MmsServer_installObtainFileHandler()` (via shim)

**Description**: Installs a callback when a client uploads a file (obtainFile). Return true to allow.

**Example**:
```go
server.InstallObtainFileHandler(func(sourceFilename, destinationFilename string) bool {
    return strings.HasPrefix(destinationFilename, "/uploads/")
})
```

---

### InstallGetFileCompleteHandler

**Go Function**: `func (is *IedServer) InstallGetFileCompleteHandler(handler GetFileCompleteHandler)`  
**C Function**: `MmsServer_installGetFileCompleteHandler()` (via shim)

**Description**: Installs a callback invoked when a file upload (obtainFile) has completed.

**Example**:
```go
server.InstallGetFileCompleteHandler(func(destinationFilename string) {
    log.Printf("File upload complete: %s", destinationFilename)
})
```

---

### SetMaxDataSetEntries

**Go Function**: `func (is *IedServer) SetMaxDataSetEntries(maxDataSetEntries int)`  
**C Function**: `MmsServer_setMaxDataSetEntries()`

**Description**: Sets the maximum number of data set entries for dynamic data sets.

**Example**:
```go
server.SetMaxDataSetEntries(64)
server.Start(102)
```

---

### EnableJournalService

**Go Function**: `func (is *IedServer) EnableJournalService(enable bool)`  
**C Function**: `MmsServer_enableJournalService()`

**Description**: Enables or disables the MMS journal service at runtime. Requires CONFIG_MMS_SERVER_CONFIG_SERVICES_AT_RUNTIME in the C library.

**Example**:
```go
server.EnableJournalService(true)
server.Start(102)
```

---

### SetMaxMmsConnections / SetMaxConnections

**Go Function**: `func (is *IedServer) SetMaxMmsConnections(maxConnections int)`  
**C Function**: `MmsServer_setMaxConnections()`

**Description**: Sets the maximum number of MMS client connections. SetMaxConnections is an alias.

**Example**:
```go
server.SetMaxMmsConnections(20)
```

---

### SetMaxAssociationSpecificDataSets / SetMaxDomainSpecificDataSets

**Go Function**: `func (is *IedServer) SetMaxAssociationSpecificDataSets(maxDataSets int)`  
**Go Function**: `func (is *IedServer) SetMaxDomainSpecificDataSets(maxDataSets int)`  
**C Function**: `MmsServer_setMaxAssociationSpecificDataSets()` / `setMaxDomainSpecificDataSets()`

**Description**: Sets maximum number of association-specific or domain-specific data sets.

**Example**:
```go
server.SetMaxAssociationSpecificDataSets(16)
server.SetMaxDomainSpecificDataSets(32)
```

---

### EnableMmsFileService / EnableDynamicNamedVariableLists

**Go Function**: `func (is *IedServer) EnableMmsFileService(enable bool)`  
**Go Function**: `func (is *IedServer) EnableDynamicNamedVariableLists(enable bool)`  
**C Function**: `MmsServer_enableFileService()` / `enableDynamicNamedVariableListService()`

**Description**: Enables or disables MMS file service or dynamic named variable list service at runtime.

**Example**:
```go
server.EnableMmsFileService(true)
server.EnableDynamicNamedVariableLists(true)
```

---

### SetConnectionIndicationHandler

**Go Function**: `func (is *IedServer) SetConnectionIndicationHandler(handler ConnectionIndicationHandler)`  
**C Function**: `IedServer_setConnectionIndicationHandler()`

**Description**: Sets handler for client connect/disconnect events.

**Example**:
```go
server.SetConnectionIndicationHandler(func(conn *iec61850.ClientConnection, connected bool) {
    if connected {
        log.Println("Client connected from:", conn.PeerAddress())
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

**Example**:
```go
server.SetClientAuthenticator(func(conn *iec61850.ClientConnection, user, pass string) bool {
    return user == "admin" && pass == "secret"
})
```

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

**Example**:
```go
server.LockDataModel()
// ... update model ...
server.UnlockDataModel()
```

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

**Example**:
```go
err := client.Select("Device/XCBR1.Pos")
```

---

### Cancel

**Go Function**: `func (c *Client) Cancel(controlRef string) error`  
**C Function**: `ControlObjectClient_cancel()`

**Description**: Cancels a selected control operation.

**Example**:
```go
err := client.Cancel("Device/XCBR1.Pos")
```

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

### NewGooseSubscriberWithDataSet

**Go Function**: `func NewGooseSubscriberWithDataSet(conf SubscriberConf, dataSetValues *GooseDataSetValues) *GooseSubscriber`  
**C Function**: `GooseSubscriber_create(goCbRef, dataSetValues)`

**Description**: Creates a GOOSE subscriber that writes received data set values into the pre-allocated MmsValue from dataSetValues. Obtain dataSetValues from Client.ReadDataSetValues and ClientDataSet.GooseDataSetValues(); the ClientDataSet must remain alive for the lifetime of the subscriber. Pass nil for dataSetValues to use auto-allocated values (same as NewGooseSubscriber).

**Example**:
```go
dataSet, _ := client.ReadDataSetValues("simpleIOGenericIO/LLN0$dataset1")
defer dataSet.Destroy()
sub := iec61850.NewGooseSubscriberWithDataSet(conf, &dataSet.GooseDataSetValues())
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

### SetObserver

**Go Function**: `func (subscriber *GooseSubscriber) SetObserver()`  
**C Function**: `GooseSubscriber_setObserver()`

**Description**: Configures the subscriber to listen to any received GOOSE message (observer mode). When set, the subscriber still has access to goCbRef, goId, and datSet of the received message.

**Example**:
```go
subscriber.SetObserver()
receiver.AddSubscriber(subscriber)
```

---

### Subscribe

**Go Function**: `func (subscriber *GooseSubscriber) Subscribe() error`  
**C Function**: `GooseReceiver_addSubscriber()`, `GooseReceiver_start()`

**Description**: Starts receiving GOOSE messages.

**Example**:
```go
err := subscriber.Subscribe()
```

---

### StartThreadless (GooseReceiver)

**Go Function**: `func (receiver *GooseReceiver) StartThreadless() *GooseReceiverSocket`  
**C Function**: `GooseReceiver_startThreadless()`

**Description**: Starts the GOOSE receiver in non-threaded mode. Returns an opaque socket handle; drive reception by calling HandleMessage with each received Ethernet frame. Call StopThreadless to stop.

**Example**:
```go
sock := receiver.StartThreadless()
if sock != nil {
    defer receiver.StopThreadless()
    // Read frames (e.g. from raw socket) and call receiver.HandleMessage(frame)
}
```

---

### StopThreadless (GooseReceiver)

**Go Function**: `func (receiver *GooseReceiver) StopThreadless()`  
**C Function**: `GooseReceiver_stopThreadless()`

**Description**: Stops the receiver when running in threadless mode (after StartThreadless).

**Example**:
```go
receiver.StopThreadless()
```

---

### HandleMessage (GooseReceiver)

**Go Function**: `func (receiver *GooseReceiver) HandleMessage(buffer []byte)`  
**C Function**: `GooseReceiver_handleMessage()`

**Description**: Parses a GOOSE message from a raw Ethernet frame. Use when driving reception yourself (e.g. with StartThreadless or custom socket reads). buffer must contain the complete Ethernet frame.

**Example**:
```go
// After reading frame from socket or elsewhere:
receiver.HandleMessage(ethernetFrame)
```

---

### NewGooseReceiverEx (GooseReceiver)

**Go Function**: `func NewGooseReceiverEx(buffer []byte) *GooseReceiver`  
**C Function**: `GooseReceiver_createEx()`

**Description**: Creates a GOOSE receiver that uses the given buffer for message handling instead of allocating its own. Pass nil or an empty slice for the default buffer. When buffer is non-nil, the receiver keeps a reference for its lifetime.

**Example**:
```go
buf := make([]byte, 1500)
receiver := iec61850.NewGooseReceiverEx(buf)
defer receiver.Destroy()
```

---

### NewGoosePublisher

**Go Function**: `func NewGoosePublisher(conf GoosePublisherConf) (*GoosePublisher, error)`  
**C Function**: `GoosePublisher_create()` (via `NewGoosePublisherEx(conf, true)`)

**Description**: Creates a GOOSE publisher with VLAN tags enabled.

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

### NewGoosePublisherEx

**Go Function**: `func NewGoosePublisherEx(conf GoosePublisherConf, useVlanTag bool) (*GoosePublisher, error)`  
**C Function**: `GoosePublisher_createEx()`

**Description**: Creates a GOOSE publisher with optional VLAN tag. Set useVlanTag to false to disable VLAN tags in sent frames when not needed.

**Example**:
```go
publisher, err := iec61850.NewGoosePublisherEx(conf, false)
```

---

### Publish

**Go Function**: `func (publisher *GoosePublisher) Publish(values []*MmsValue) error`  
**C Function**: `GoosePublisher_publish()`

**Description**: Publishes GOOSE message with data values.

**Example**:
```go
vals := []*iec61850.MmsValue{intVal, boolVal}
err := publisher.Publish(vals)
```

---

### SetGoID

**Go Function**: `func (publisher *GoosePublisher) SetGoID(goID string)`  
**C Function**: `GoosePublisher_setGoID()`

**Description**: Sets the GOOSE identifier string sent in GOOSE messages (e.g. when it differs from GoCbRef).

**Example**:
```go
publisher.SetGoID("MyDevice/LLN0$GO$gcb1")
```

---

### PublishAndDump (GoosePublisher)

**Go Function**: `func (publisher *GoosePublisher) PublishAndDump(dataSet *LinkedListValue, msgBuf []byte) (msgLen int, err error)`  
**C Function**: `GoosePublisher_publishAndDump()`

**Description**: Publishes a GOOSE message and copies the raw encoded payload into msgBuf. Returns the number of bytes written (use msgBuf[:msgLen]). msgBuf must be non-nil and have positive length (e.g. 1500+ bytes).

**Example**:
```go
msgBuf := make([]byte, 1500)
n, err := publisher.PublishAndDump(dataSet, msgBuf)
if err == nil {
    rawPayload := msgBuf[:n]
    // log or inspect rawPayload
}
```

---

## GOOSE Control Block (GoCB) Client Functions

### GetGoCBValuesAsync

**Go Function**: `func (c *Client) GetGoCBValuesAsync(goCBReference string, callback func(*ClientGooseControlBlockValues, error)) (uint32, error)`  
**C Function**: `IedConnection_getGoCBValuesAsync()`

**Description**: Reads GOOSE control block values from the server asynchronously. The callback is invoked with the values and nil error on success, or nil and error on failure. Returns the invoke ID and an error if the request could not be started.

**Example**:
```go
invokeID, err := client.GetGoCBValuesAsync("Device/LLN0.gcb1", func(v *iec61850.ClientGooseControlBlockValues, err error) {
    if err != nil { log.Println(err); return }
    fmt.Printf("GoEna: %v, GoID: %s\n", v.GoEna, v.GoID)
})
```

---

### SetGoCBValuesAsync

**Go Function**: `func (c *Client) SetGoCBValuesAsync(goCBReference string, values *ClientGooseControlBlockValues, parametersMask uint32, singleRequest bool, callback func(error)) (uint32, error)`  
**C Function**: `IedConnection_setGoCBValuesAsync()`

**Description**: Writes GOOSE control block values to the server asynchronously. The callback is invoked with nil on success or an error on failure. Returns the invoke ID and an error if the request could not be started.

**Example**:
```go
invokeID, err := client.SetGoCBValuesAsync("Device/LLN0.gcb1", values, iec61850.GoCBElementDstAddress, false, func(err error) {
    if err != nil { log.Println(err) }
})
```

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

**Example**:
```go
conf := iec61850.SVPublisherConf{
    InterfaceID: "eth0",
    AppID:       4000,
    SvID:        "SVPublisher1",
}
pub, err := iec61850.NewSVPublisher(conf)
```

---

### PublishSV

**Go Function**: `func (publisher *SVPublisher) PublishSV(asdu *SVPublisherASDU) error`  
**C Function**: `SVPublisher_publish()`

**Description**: Publishes sampled values ASDU.

**Example**:
```go
asdu := &iec61850.SVPublisherASDU{SvID: "SV1", SmpCnt: 0, Samples: samples}
err := publisher.PublishSV(asdu)
```

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

### GetFileDirectoryExEntries

**Go Function**: `func (c *Client) GetFileDirectoryExEntries(directoryName, continueAfter string) ([]MmsFileDirectoryEntryEx, bool, error)`  
**C Function**: Uses GetFileDirectoryEx internally

**Description**: Returns the file directory as a slice of MmsFileDirectoryEntryEx (Filename, FileSize, LastModifiedTime, FileAttributes). FileAttributes may be 0 if the server does not provide them.

**Example**:
```go
entries, more, err := client.GetFileDirectoryExEntries("/", "")
if err != nil {
    log.Fatal(err)
}
for _, e := range entries {
    fmt.Printf("%s %d bytes\n", e.Filename, e.FileSize)
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

**Example**:
```go
frsmID, err := client.FileOpen("/config/settings.cfg", true)
```

---

### FileRead

**Go Function**: `func (c *Client) FileRead(frsmID uint32, bufferSize int) ([]byte, bool, error)`  
**C Function**: `IedConnection_fileRead()`

**Description**: Reads chunk from open file.

**Example**:
```go
data, more, err := client.FileRead(frsmID, 1024)
```

---

### FileClose

**Go Function**: `func (c *Client) FileClose(frsmID uint32) error`  
**C Function**: `IedConnection_fileClose()`

**Description**: Closes an open file.

**Example**:
```go
_ = client.FileClose(frsmID)
```

---

### FileDelete

**Go Function**: `func (c *Client) FileDelete(fileName string) error`  
**C Function**: `IedConnection_fileDelete()`

**Description**: Deletes a file on the server.

**Example**:
```go
err := client.FileDelete("/temp/old.cfg")
```

---

### ObtainFile

**Go Function**: `func (c *Client) ObtainFile(sourceFile, destFile string) error`  
**C Function**: `MmsConnection_obtainFile()`

**Description**: Requests server to upload a file from client (client→server transfer).

**Example**:
```go
err := client.ObtainFile("/local/config.cfg", "/server/config.cfg")
```

---

### RenameFile

**Go Function**: `func (c *Client) RenameFile(currentName, newName string) error`  
**C Function**: `MmsConnection_fileRename()`

**Description**: Renames a file on the server.

**Example**:
```go
err := client.RenameFile("/old.txt", "/new.txt")
```

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

**Example**:
```go
entries := []*iec61850.MmsVariableAccessSpec{
    {DomainID: "domain", ItemID: "var1"},
    {DomainID: "domain", ItemID: "var2"},
}
err := client.CreateDataSet("Device/LLN0$mydataset", entries)
```

---

### DeleteDataSet

**Go Function**: `func (c *Client) DeleteDataSet(dataSetRef string) error`  
**C Function**: `IedConnection_deleteDataSet()`

**Description**: Deletes a dataset.

**Example**:
```go
err := client.DeleteDataSet("Device/LLN0$mydataset")
```

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

**Example**:
```go
// When calling C API that returns MmsError
var cErr C.MmsError
// ... C call sets cErr ...
if err := iec61850.GetMmsError(cErr); err != nil {
    log.Println("MMS error:", err)
}
```

---

### GetIedClientError

**Go Function**: `func GetIedClientError(err C.IedClientError) error`  
**C Function**: N/A (helper)

**Description**: Converts C IedClientError to Go error.

**Example**:
```go
// When C returns IedClientError
if err := iec61850.GetIedClientError(cErr); err != nil {
    log.Println("Client error:", err)
}
```

---

### IsBitSet

**Go Function**: `func IsBitSet(val int, pos int) bool`  
**C Function**: N/A (utility)

**Description**: Checks if a bit is set at position.

**Example**:
```go
optFlds := 0x07
if iec61850.IsBitSet(optFlds, 0) {
    fmt.Println("Bit 0 is set")
}
```

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

**Example**:
```go
mmsType, err := client.GetVariableSpecType("Device/GGIO1.AnIn1.mag.f", iec61850.MX)
if err == nil {
    fmt.Println("Type:", mmsType)
}
```

---

## Type System (MmsVariableSpecificationRef)

Type introspection for MMS variables. See [GAPS.md](GAPS.md) Part 4. All 12 C functions are implemented.

### Methods on MmsVariableSpecificationRef

| Go Method | C Function | Description |
|-----------|------------|-------------|
| `GetType()` | `MmsVariableSpecification_getType()` | MMS type code |
| `GetName()` | `MmsVariableSpecification_getName()` | Type name |
| `GetSize()` | `MmsVariableSpecification_getSize()` | Size (e.g. array count) |
| `GetChildSpecificationByIndex(i)` | `getChildSpecificationByIndex()` | Child type by index |
| `GetChildSpecificationByName(name)` | `getChildSpecificationByName()` | Child type by name |
| `GetArrayElementSpecification()` | `getArrayElementSpecification()` | Element type of array |
| `IsValueOfType(value)` | `isValueOfType()` | Checks value matches this spec |
| `GetChildValue(value, childId)` | `getChildValue()` | Extracts child from structure value |
| `GetNamedVariableRecursive(nameId)` | `getNamedVariableRecursive()` | Nested spec by name |
| `GetExponentWidth()` | `getExponentWidth()` | Float exponent width |
| `GetStructureElements()` | `getStructureElements()` | List of structure element names |
| `Free()` | `MmsVariableSpecification_destroy()` | Releases C spec (when owned) |

### Constructors

**NewMmsVariableSpecification** – `func NewMmsVariableSpecification(typ MmsType, name string, size int) *MmsVariableSpecificationRef`  
Creates a type spec for a primitive or array.

**CreateStructure** – `func CreateStructure(name string, elements []*MmsVariableSpecificationRef) *MmsVariableSpecificationRef`  
Creates a structure type from element specs.

**CreateArray** – `func CreateArray(name string, elementType *MmsVariableSpecificationRef, elementCount int) *MmsVariableSpecificationRef`  
Creates an array type.

**Example**:
```go
spec, _ := client.GetVariableAccessAttributes("domain", "var")
defer spec.Free()
if spec.GetType() == iec61850.Structure {
    for _, name := range spec.GetStructureElements() {
        child := spec.GetChildSpecificationByName(name)
        fmt.Printf("  %s: %v\n", name, child.GetType())
    }
}
```

---

*End of Functions Reference*
