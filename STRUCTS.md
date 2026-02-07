# IEC 61850 Go Bindings - Structs Reference

**Version**: 1.6.1  
**Generated**: February 7, 2026

This document provides comprehensive documentation for all exported Go structs/types in the iec61850 package, including their corresponding C structures, field descriptions, and usage examples. It is aligned with [GAPS.md](GAPS.md) (MMS coverage analysis).

---

## Table of Contents

1. [Client Types](#client-types)
2. [Server Types](#server-types)
3. [MMS Types](#mms-types)
4. [Connection Types](#connection-types)
5. [Data Model Types](#data-model-types)
6. [Control Types](#control-types)
7. [Reporting Types](#reporting-types)
8. [GOOSE Types](#goose-types)
9. [Sampled Values Types](#sampled-values-types)
10. [File Service Types](#file-service-types)
11. [Configuration Types](#configuration-types)
12. [Time & Quality Types](#time--quality-types)

---

## Client Types

### Client

**Go Type**: `type Client struct`  
**C Type**: `IedConnection`

**Description**: Represents a client connection to an IEC 61850 server.

**Fields**: (Internal/opaque)

**Example**:
```go
client, err := iec61850.NewClient(settings)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

value, _ := client.Read("Device/GGIO1.AnIn1.mag.f", iec61850.MX)
```

---

### Settings

**Go Type**: 
```go
type Settings struct {
    Host           string
    Port           int
    ConnectTimeout uint // milliseconds
    RequestTimeout uint // milliseconds
}
```

**C Type**: N/A (Go wrapper)

**Description**: Connection settings for IEC 61850 client.

**Example**:
```go
settings := iec61850.Settings{
    Host:           "192.168.1.10",
    Port:           102,
    ConnectTimeout: 10000,
    RequestTimeout: 5000,
}
```

---

### LastApplError

**Go Type**:
```go
type LastApplError struct {
    CtlNum   int // Control number
    Error    int // Error code
    AddCause int // Additional cause
}
```

**C Type**: `LastApplError`

**Description**: Contains error information from control operations.

**Example**:
```go
lastErr := client.GetLastApplError()
if lastErr.Error != 0 {
    fmt.Printf("Control error %d (cause: %d)\n", lastErr.Error, lastErr.AddCause)
}
```

---

### ClientDataSet

**Go Type**: `type ClientDataSet struct` (opaque)  
**C Type**: `ClientDataSet`

**Description**: Holds data set values read from the server via ReadDataSetValues. Call Destroy when done. The underlying MmsValue (MMS_ARRAY) can be passed to NewGooseSubscriberWithDataSet via GooseDataSetValues(); keep the ClientDataSet alive for the lifetime of that subscriber.

**Example**:
```go
dataSet, err := client.ReadDataSetValues("simpleIOGenericIO/LLN0$dataset1")
if err != nil { log.Fatal(err) }
defer dataSet.Destroy()
sub := iec61850.NewGooseSubscriberWithDataSet(conf, &dataSet.GooseDataSetValues())
```

---

### GooseDataSetValues

**Go Type**: `type GooseDataSetValues struct` (opaque; field `p unsafe.Pointer`)  
**C Type**: Wraps `MmsValue*` from `ClientDataSet_getValues()`

**Description**: Opaque handle to the MmsValue array from a ClientDataSet, for use with NewGooseSubscriberWithDataSet. Obtain from ClientDataSet.GooseDataSetValues(); keep the ClientDataSet alive for the lifetime of the subscriber.

**Example**:
```go
gooseVals := dataSet.GooseDataSetValues()
sub := iec61850.NewGooseSubscriberWithDataSet(conf, &gooseVals)
```

---

### IedConnectionState

**Go Type**: `type IedConnectionState int`

**C Type**: `IedConnectionState`

**Description**: Connection state enumeration.

**Values**:
- `IedStateClosed` (0)
- `IedStateConnecting` (1)
- `IedStateConnected` (2)
- `IedStateClosing` (3)

**Example**:
```go
if client.GetState() == iec61850.IedStateConnected {
    fmt.Println("Connected to server")
}
```

---

## Server Types

### IedServer

**Go Type**: `type IedServer struct`  
**C Type**: `IedServer`

**Description**: Represents an IEC 61850 server instance.

**Fields**: (Internal/opaque)

**Example**:
```go
model := iec61850.LoadModel("model.cfg")
server := iec61850.NewServer(model)
defer server.Destroy()

server.Start(102)
if server.IsRunning() {
    fmt.Println("Server started successfully")
}
```

---

### ServerConfig

**Go Type**:
```go
type ServerConfig struct {
    Edition                      uint8
    ReportBufferSize             int
    ReportBufferSizeForURCBs     int
    MaxConnections               int
    SyncIntegrityReportTimes     bool
    EnableFileService            bool
    FileServiceBasePath          string
    EnableDynamicDataSetService  bool
    MaxDomainSpecificDataSets    int
    MaxAssociationSpecificDataSets int
    MaxDataSetEntries            int
    EnableEditSG                 bool
    ReserveEditSGTimeout         int
}
```

**C Type**: `IedServerConfig`

**Description**: Configuration for IED server.

**Example**:
```go
config := iec61850.ServerConfig{
    Edition:                  iec61850.IEC_61850_EDITION_2,
    ReportBufferSize:         100000,
    ReportBufferSizeForURCBs: 50000,
    MaxConnections:           10,
    EnableFileService:        true,
    FileServiceBasePath:      "./vmd-filestore/",
}
server := iec61850.NewServerWithConfig(config, model)
```

---

### ClientConnection

**Go Type**: `type ClientConnection struct`  
**C Type**: `ClientConnection`

**Description**: Represents a client connection on the server side.

**Methods**: `PeerAddress() string`, `LocalAddress() string`, `SecurityToken() unsafe.Pointer`, `Abort() bool`, `ClaimOwnership() *ClientConnection`, `Release()`

**Example**:
```go
server.SetConnectionIndicationHandler(func(conn *iec61850.ClientConnection, connected bool) {
    if connected {
        addr := conn.PeerAddress()
        fmt.Printf("Client connected from %s\n", addr)
    }
})
```

---

### ClientAuthenticator

**Go Type**: `type ClientAuthenticator func(param interface{}, conn *ClientConnection, authnParameter []byte) bool`

**C Type**: `IedServer_ClientAuthenticator` (callback)

**Description**: Callback function for authenticating client connections.

**Example**:
```go
server.SetClientAuthenticator(func(param interface{}, conn *iec61850.ClientConnection, authnParam []byte) bool {
    // Validate client credentials
    username string(authnParam)
    return username == "admin"
})
```

---

### ConnectionIndicationHandler

**Go Type**: `type ConnectionIndicationHandler func(connection *ClientConnection, connected bool)`

**C Type**: `IedConnectionIndicationHandler` (callback)

**Description**: Callback for client connection/disconnection events.

**Example**:
```go
server.SetConnectionIndicationHandler(func(conn *iec61850.ClientConnection, connected bool) {
    if connected {
        fmt.Println("Client connected:", conn.PeerAddress())
    }
})
```

---

### Server MMS handler types

**FileAccessHandler** – `func(service MmsFileServiceType, localFilename, otherFilename string) error`  
Used with SetFileAccessHandler. Return nil to allow, or an error (e.g. AccessDenied) to deny.

**VariableListAccessHandler** – `func(accessType MmsVariableListAccessType, listType MmsVariableListType, domainID, listName string) error`  
Used with InstallVariableListAccessHandler. Return nil to allow.

**ReadJournalHandler** – `func(domainID, logName string) bool`  
Used with InstallReadJournalHandler. Return true to allow, false to deny.

**GetNameListHandler** – `func(nameListType MmsGetNameListType, domainID string) bool`  
Used with InstallGetNameListHandler. Return true to allow.

**ObtainFileHandler** – `func(sourceFilename, destinationFilename string) bool`  
Used with InstallObtainFileHandler (file upload). Return true to allow.

**GetFileCompleteHandler** – `func(destinationFilename string)`  
Used with InstallGetFileCompleteHandler; called when file upload completes.

**Example** (combined):
```go
server.SetFileAccessHandler(func(svc iec61850.MmsFileServiceType, local, other string) error { return nil })
server.InstallVariableListAccessHandler(func(accessType iec61850.MmsVariableListAccessType, listType iec61850.MmsVariableListType, domainID, listName string) error { return nil })
server.InstallReadJournalHandler(func(domainID, logName string) bool { return true })
server.InstallGetNameListHandler(func(nameListType iec61850.MmsGetNameListType, domainID string) bool { return true })
server.InstallObtainFileHandler(func(src, dst string) bool { return true })
server.InstallGetFileCompleteHandler(func(dst string) { log.Println("Upload complete:", dst) })
```

---

### LogStorageRef

**Go Type**: `type LogStorageRef struct` (wraps C LogStorage)  
**C Type**: `LogStorage*`

**Description**: Wraps a C LogStorage pointer (e.g. from SqliteLogStorage_createInstance when built with sqlite). Create with NewLogStorageRef(ptr), call SetMaxLogEntries, Destroy when done.

**Example**:
```go
// ptr from C e.g. SqliteLogStorage_createInstance(...)
ref := iec61850.NewLogStorageRef(ptr)
ref.SetMaxLogEntries(10000)
defer ref.Destroy()
```

---

## MMS Types

### MmsConnection

**Go Type**: `type MmsConnection struct`  
**C Type**: `MmsConnection`

**Description**: Low-level MMS connection object.

**Example**:
```go
mmsConn := iec61850.NewMmsConnection()
defer mmsConn.Destroy()

err := mmsConn.ConnectAsync("192.168.1.10", 102, func(err error) {
    if err == nil {
        fmt.Println("Connected via MMS")
    }
})
```

---

### MmsValue

**Go Type**:
```go
type MmsValue struct {
    Type  MmsType
    Value interface{}
}
```

**C Type**: `MmsValue`

**Description**: High-level MMS value wrapper.

**Example**:
```go
intVal, _ := iec61850.NewMmsValue(iec61850.Integer, int64(42))
boolVal, _ := iec61850.NewMmsValue(iec61850.Boolean, true)
```

---

### MmsValueRef

**Go Type**: `type MmsValueRef struct` (opaque)  
**C Type**: `MmsValue*`

**Description**: Low-level reference to MMS value (direct C binding).

**Example**:
```go
bitStr := iec61850.NewMmsValueBitString(16)
bitStr.SetBitStringFromInteger(0xAAAA)
val := bitStr.GetBitStringAsInteger()
```

---

### MmsType

**Go Type**: `type MmsType int`

**C Type**: `MmsType`

**Description**: MMS data type enumeration.

**Values**: (See ENUMS.md for complete list)
- `Array`, `Structure`, `Boolean`, `BitString`, `Integer`, `Unsigned`, `Float`, `OctetString`, `VisibleString`, `UTCTime`, etc.

**Example**:
```go
if mmsVal.GetType() == iec61850.Integer {
    i := mmsVal.ToInt64()
}
```

---

### MmsDataAccessError

**Go Type**: `type MmsDataAccessError int`

**C Type**: `MmsDataAccessError`

**Description**: MMS data access error codes.

**Values**: 
- `DATA_ACCESS_ERROR_SUCCESS`
- `DATA_ACCESS_ERROR_OBJECT_UNDEFINED`
- `DATA_ACCESS_ERROR_OBJECT_ACCESS_DENIED`
- etc.

**Example**:
```go
var results []iec61850.MmsDataAccessError
_ = mmsConn.WriteNamedVariableList("domain", "list1", values, &results)
for i, code := range results {
    if code != iec61850.DATA_ACCESS_ERROR_SUCCESS {
        fmt.Printf("Item %d: %v\n", i, code)
    }
}
```

---

### MmsConnectionParameters

**Go Type**:
```go
type MmsConnectionParameters struct {
    MaxServOutstandingCalling int32
    MaxServOutstandingCalled  int32
    DataStructureNestingLevel int32
    MaxPduSize                int32
    ServicesSupported         [11]uint8
}
```

**C Type**: `MmsConnectionParameters`

**Description**: MMS connection negotiated parameters.

**Example**:
```go
params := mmsConn.GetMmsConnectionParameters()
fmt.Printf("Max PDU Size: %d bytes\n", params.MaxPduSize)
fmt.Printf("Max Outstanding Calls: %d\n", params.MaxServOutstandingCalling)
```

---

### IsoConnectionParameters

**Go Type**:
```go
type IsoConnectionParameters struct {
    LocalApTitle      []byte
    LocalAeQualifier  int32
    RemoteApTitle     []byte
    RemoteAeQualifier int32
    LocalTSelector    []byte
    LocalSSelector    []byte
    LocalPSelector    []byte
    RemoteTSelector   []byte
    RemoteSSelector   []byte
    RemotePSelector   []byte
}
```

**C Type**: `IsoConnectionParameters`

**Description**: ISO connection layer parameters (AP titles and OSI layer selectors). Get/Set via MmsConnection.GetIsoConnectionParameters / SetIsoConnectionParameters.

**Example**:
```go
isoParams := mmsConn.GetIsoConnectionParameters()
fmt.Printf("Local T-Selector: %x\n", isoParams.LocalTSelector)
```

---

### MmsServerIdentity

**Go Type**:
```go
type MmsServerIdentity struct {
    VendorName string
    ModelName  string
    Revision   string
}
```

**C Type**: `MmsServerIdentity`

**Description**: Server identification information.

**Example**:
```go
identity, err := client.Identify()
if err == nil {
    fmt.Printf("Server: %s %s (Rev %s)\n",
        identity.VendorName, identity.ModelName, identity.Revision)
}
```

---

### MmsServerStatus

**Go Type**:
```go
type MmsServerStatus struct {
    VmdLogicalStatus  int32
    VmdPhysicalStatus int32
    LocalDetail       int32
}
```

**C Type**: N/A (extracted from MmsValue)

**Description**: MMS server status information.

**Example**:
```go
status, err := client.GetServerStatus(false)
if err == nil {
    fmt.Printf("VMD logical=%d physical=%d localDetail=%d\n",
        status.VmdLogicalStatus, status.VmdPhysicalStatus, status.LocalDetail)
}
```

---

### MmsVariableSpecificationRef

**Go Type**: `type MmsVariableSpecificationRef struct` (opaque)  
**C Type**: `MmsVariableSpecification*`

**Description**: MMS type specification for introspection.

**Example**:
```go
typeSpec, err := client.GetVariableAccessAttributes("domain", "variable")
if err == nil {
    fmt.Printf("Type: %v\n", typeSpec.GetType())
    fmt.Printf("Size: %d\n", typeSpec.GetSize())
    if typeSpec.GetType() == iec61850.Structure {
        elements := typeSpec.GetStructureElements()
        fmt.Printf("Structure elements: %v\n", elements)
    }
}
```

---

### MmsFileDirectoryEntryEx

**Go Type**:
```go
type MmsFileDirectoryEntryEx struct {
    Filename         string
    FileSize         uint32
    LastModifiedTime uint64
    FileAttributes   uint32  // May be 0 if server does not provide
}
```

**C Type**: N/A (built from GetFileDirectoryEx result)

**Description**: Extended file directory entry (from GetFileDirectoryExEntries). FileAttributes may be 0 if the server does not provide them.

**Example**:
```go
entries, _, _ := client.GetFileDirectoryExEntries("/", "")
for _, e := range entries {
    fmt.Printf("%s: %d bytes, mtime=%d\n", e.Filename, e.FileSize, e.LastModifiedTime)
}
```

---

### VariableAccessSpec

**Go Type**:
```go
type VariableAccessSpec struct {
    DomainID string
    ItemID   string
}
```

**C Type**: `MmsVariableAccessSpecification`

**Description**: MMS variable access specification (domain + item).

**Example**:
```go
spec := &iec61850.VariableAccessSpec{
    DomainID: "simpleIOGenericIO",
    ItemID:   "GGIO1$ST$Ind1$stVal",
}
```

---

### VariableListEntry

**Go Type**:
```go
type VariableListEntry struct {
    VariableName string
    VariableSpec *VariableAccessSpec
}
```

**C Type**: Part of named variable list

**Description**: Entry in a named variable list.

**Example**:
```go
entry := iec61850.VariableListEntry{
    VariableName: "var1",
    VariableSpec: &iec61850.VariableAccessSpec{DomainID: "domain", ItemID: "item1"},
}
```

---

### JournalVariable

**Go Type**:
```go
type JournalVariable struct {
    Tag   string
    Value *MmsValue
}
```

**C Type**: Part of journal entry

**Description**: Variable in a journal entry.

**Example**:
```go
jvar := iec61850.JournalVariable{
    Tag:   "stVal",
    Value: mmsVal,
}
```

---

### JournalEntry

**Go Type**:
```go
type JournalEntry struct {
    EntryID        *MmsValue       // Octet string
    OccurrenceTime *MmsValue       // Binary time
    Variables      []JournalVariable
}
```

**C Type**: N/A (constructed from journal read response)

**Description**: Complete journal entry with timestamp and data.

**Example**:
```go
entries, moreFollows, err := client.ReadJournalTimeRange(
    "domain", "journal", startTimeMs, endTimeMs)
for _, entry := range entries {
    for _, jvar := range entry.Variables {
        fmt.Printf("%s = %v\n", jvar.Tag, jvar.Value)
    }
}
```

---

### MmsJournalEntry

**Go Type** (from types.go):
```go
type MmsJournalEntry struct {
    EntryID      []byte
    OccurTime    uint64
    EntryContent *MmsValue
}
```

**C Type**: N/A (built from MmsConnection journal read)

**Description**: Low-level journal entry used by MmsConnection.ReadJournalTimeRange / ReadJournalStartAfter. For Client API the higher-level JournalEntry (with Variables) is used.

**Example**:
```go
entries, more, err := mmsConn.ReadJournalTimeRange("domain", "journal", startMs, endMs)
for _, e := range entries {
    fmt.Printf("Entry %x at %d\n", e.EntryID, e.OccurTime)
}
```

---

## Connection Types

### TLSConfig

**Go Type**:
```go
type TLSConfig struct {
    OwnCertificate      []byte
    OwnKey              []byte
    CACerts             [][]byte
    ChainValidation     bool
    AllowOnlyKnownCerts bool
}
```

**C Type**: `TLSConfiguration`

**Description**: TLS configuration for secure connections.

**Example**:
```go
certPEM, _ := os.ReadFile("client_cert.pem")
keyPEM, _ := os.ReadFile("client_key.pem")
caPEM, _ := os.ReadFile("ca_cert.pem")

tlsConfig := &iec61850.TLSConfig{
    OwnCertificate:  certPEM,
    OwnKey:          keyPEM,
    CACerts:         [][]byte{caPEM},
    ChainValidation: true,
}
```

---

### TLSConfiguration

**Go Type**:
```go
type TLSConfiguration struct {
    ChainValidation      bool
    AllowOnlyKnownCerts  bool
    CACertificates       [][]byte
    OwnCertificate       []byte
    OwnKey               []byte
}
```

**C Type**: `TLSConfiguration`

**Description**: Alternative TLS configuration structure (MMS-level).

---

## Data Model Types

### IedModel

**Go Type**:
```go
type IedModel struct {
    Model         C.IedModel
    Name          string
    ModelBuf      []byte
}
```

**C Type**: `IedModel*`

**Description**: IEC 61850 information model.

**Example**:
```go
model := iec61850.LoadModelFromFile("model.cfg")
defer model.Destroy()

server := iec61850.NewServer(model)
```

---

### ModelNode

**Go Type**:
```go
type ModelNode struct {
    ObjectRef    string
    ModelNodeRef C.ModelNode
}
```

**C Type**: `ModelNode`

**Description**: Generic node in the data model tree.

**Example**:
```go
node := model.GetModelNodeByObjectReference("Device/GGIO1.Beh.stVal", iec61850.ST)
fmt.Println("ObjectRef:", node.ObjectRef)
```

---

### LogicalDevice

**Go Type**:
```go
type LogicalDevice struct {
    Parent       *IedModel
    Name         string
    ModelNodeRef C.LogicalDevice
}
```

**C Type**: `LogicalDevice`

**Description**: Logical device in IEC 61850 model.

**Example**:
```go
ld := model.GetLogicalDeviceByName("Device")
lnList := ld.GetLogicalNodes()
```

---

### LogicalNode

**Go Type**:
```go
type LogicalNode struct {
    Parent       *LogicalDevice
    Name         string
    ModelNodeRef C.LogicalNode
}
```

**C Type**: `LogicalNode`

**Description**: Logical node (LN) in the model.

**Example**:
```go
ln := ld.GetLogicalNodeByName("GGIO1")
dataObjects := ln.GetDataObjects()
```

---

### DataObject

**Go Type**:
```go
type DataObject struct {
    Parent       *LogicalNode
    Name         string
    ModelNodeRef C.DataObject
}
```

**C Type**: `DataObject`

**Description**: Data object in logical node.

**Example**:
```go
do := ln.GetDataObjectByName("Beh")
```

---

### DataAttribute

**Go Type**:
```go
type DataAttribute struct {
    Parent       interface{} // *DataObject or *DataAttribute
    Name         string
    Fc           FC
    ModelNodeRef C.DataAttribute
}
```

**C Type**: `DataAttribute`

**Description**: Data attribute (leaf node with actual value).

**Example**:
```go
attr := do.GetChild("stVal", iec61850.ST)
value := server.GetAttributeValue(attr)
```

---

### DataSet

**Go Type**:
```go
type DataSet struct {
    Name         string
    ModelNodeRef C.DataSet
}
```

**C Type**: `DataSet`

**Description**: Dataset collection of data attributes.

**Example**:
```go
ds := ln.GetDataSetByName("dataset1")
members, _ := client.GetDataSetDirectory(ds.GetReference())
```

---

### DataModel

**Go Type**:
```go
type DataModel struct {
    Name string
    LDs  []*LD
}
```

**C Type**: N/A (Go helper)

**Description**: High-level data model representation.

**Example**:
```go
model := iec61850.GetDataModel(iedModel)
for _, ld := range model.LDs {
    fmt.Println("LD:", ld.Inst)
}
```

---

### LD

**Go Type**:
```go
type LD struct {
    Inst string
    LNs  []*LN
}
```

**Description**: Logical device in DataModel.

**Example**:
```go
for _, ld := range model.LDs {
    fmt.Println("LD:", ld.Inst)
    for _, ln := range ld.LNs {
        fmt.Println("  LN:", ln.LnClass, ln.LnInst)
    }
}
```

---

### LN

**Go Type**:
```go
type LN struct {
    LnPrefix string
    LnClass  string
    LnInst   string
    DOs      []*DO
}
```

**Description**: Logical node in DataModel.

**Example**:
```go
for _, ln := range ld.LNs {
    for _, do := range ln.DOs {
        fmt.Println("DO:", do.Name)
    }
}
```

---

### DO

**Go Type**:
```go
type DO struct {
    Name  string
    DAs   []*DA
    SDOs  []*SDO
}
```

**Description**: Data object in LN.

**Example**:
```go
for _, do := range ln.DOs {
    for _, da := range do.DAs {
        fmt.Printf("DA %s FC=%v\n", da.Name, da.Fc)
    }
}
```

---

### DA

**Go Type**:
```go
type DA struct {
    Name  string
    Fc    FC
    Type  MmsType
    Value interface{}
    BDAs  []*BDA
}
```

**Description**: Data attribute in DO.

**Example**:
```go
for _, da := range do.DAs {
    if da.Fc == iec61850.MX {
        fmt.Printf("%s = %v\n", da.Name, da.Value)
    }
}
```

---

## Control Types

### ControlObjectParam

**Go Type**:
```go
type ControlObjectParam struct {
    CtlVal    interface{}
    CtlNum    uint
    Origin    *ControlOriginator
    Test      bool
    Timestamp uint64
}
```

**C Type**: N/A (passed to control operations)

**Description**: Parameters for basic control operations (SPC, DPC).

**Example**:
```go
param := iec61850.ControlObjectParam{
    CtlVal: true,
    CtlNum: 1,
    Origin: iec61850.NewControlOriginator(
        iec61850.CONTROL_ORCAT_AUTOMATIC, 
        "AutomationSystem"),
    Test:   false,
}
err := client.Operate("Device/XCBR1.Pos", param)
```

---

### ControlObjectParamAPC

**Go Type**:
```go
type ControlObjectParamAPC struct {
    CtlVal    float32
    CtlNum    uint
    Origin    *ControlOriginator
    Test      bool
    Timestamp uint64
}
```

**C Type**: N/A

**Description**: Parameters for analog position control (APC).

**Example**:
```go
param := iec61850.ControlObjectParamAPC{
    CtlVal: 42.5,
    CtlNum: 1,
    Origin: iec61850.NewControlOriginator(iec61850.CONTROL_ORCAT_AUTOMATIC, "SCADA"),
    Test:   false,
}
```

---

### ControlObjectParamINC

**Go Type**:
```go
type ControlObjectParamINC struct {
    CtlVal    int32
    CtlNum    uint
    Origin    *ControlOriginator
    Test      bool
    Timestamp uint64
}
```

**C Type**: N/A

**Description**: Parameters for integer number control (INC).

**Example**:
```go
param := iec61850.ControlObjectParamINC{
    CtlVal: 100,
    CtlNum: 1,
    Origin: iec61850.NewControlOriginator(iec61850.CONTROL_ORCAT_AUTOMATIC, "SCADA"),
    Test:   false,
}
```

---

### ControlModel

**Go Type**: `type ControlModel int`

**C Type**: `ControlModel`

**Description**: Control model enumeration.

**Values**:
- `CONTROL_MODEL_STATUS_ONLY`
- `CONTROL_MODEL_DIRECT_NORMAL`
- `CONTROL_MODEL_SBO_NORMAL`
- `CONTROL_MODEL_DIRECT_ENHANCED`
- `CONTROL_MODEL_SBO_ENHANCED`

**Example**:
```go
if ctrlModel == iec61850.CONTROL_MODEL_SBO_NORMAL {
    client.Select("Device/XCBR1.Pos")
    client.Operate("Device/XCBR1.Pos", param)
}
```

---

### ControlHandlerResult

**Go Type**: `type ControlHandlerResult int`

**C Type**: `ControlHandlerResult`

**Description**: Result from control handler callback.

**Values**:
- `CONTROL_RESULT_FAILED`
- `CONTROL_RESULT_OK`
- `CONTROL_RESULT_WAITING`

**Example**:
```go
// In server control handler
return iec61850.CONTROL_RESULT_OK
```

---

## Reporting Types

### ClientReportControlBlock

**Go Type**:
```go
type ClientReportControlBlock struct {
    RptID      string
    RptEna     bool
    Resv       bool
    DatSet     string
    OptFlds    OptFlds
    TrgOps     TrgOps
    IntgPd     uint32
    GI         bool
    Buffered   bool
    ConfRev    uint32
    EntryID    []byte
    TimeOfEntry uint64
    ResvTms    int32
}
```

**C Type**: `ClientReportControlBlock`

**Description**: Report control block (RCB) settings and status.

**Example**:
```go
rcb, _ := client.GetRCBValues("Device/LLN0.BR.brcb01")
fmt.Printf("RCB %s: Enabled=%v, IntgPd=%d ms\n", 
    rcb.RptID, rcb.RptEna, rcb.IntgPd)

rcb.RptEna = true
rcb.IntgPd = 5000
client.SetRCBValues("Device/LLN0.BR.brcb01", *rcb)
```

---

### TrgOps

**Go Type**:
```go
type TrgOps struct {
    DataChange     bool
    QualityChange  bool
    DataUpdate     bool
    Integrity      bool
    GeneralInterrog bool
}
```

**C Type**: Bitfield in `TrgOps`

**Description**: Trigger options for reporting.

**Example**:
```go
rcb.TrgOps = iec61850.TrgOps{
    DataChange:    true,
    QualityChange: true,
    Integrity:     true,
}
```

---

### OptFlds

**Go Type**:
```go
type OptFlds struct {
    SeqNum     bool
    TimeStamp  bool
    ReasonCode bool
    DataSet    bool
    DataRef    bool
    BufOvfl    bool
    EntryID    bool
    ConfRev    bool
}
```

**C Type**: Bitfield in `OptFlds`

**Description**: Optional fields to include in reports.

**Example**:
```go
rcb.OptFlds = iec61850.OptFlds{
    SeqNum:     true,
    TimeStamp:  true,
    ReasonCode: true,
    DataRef:    true,
}
```

---

### ReportCallbackFunc

**Go Type**: `type ReportCallbackFunc func(param interface{}, report *ClientReport)`

**C Type**: `ReportCallbackFunction`

**Description**: Callback for receiving reports.

**Example**:
```go
handler := func(param interface{}, report *iec61850.ClientReport) {
    fmt.Printf("Report: %s\n", report.RptID)
    for i, val := range report.DataSetValues {
        fmt.Printf("  [%d] = %v\n", i, val)
    }
}
client.InstallReportHandler("Device/LLN0.BR.brcb01", handler, nil)
```

---

## GOOSE Types

### CommParameters

**Go Type**:
```go
type CommParameters struct {
    VlanPriority uint8
    VlanID       uint16
    AppID        uint16
    DstAddr      [6]uint8
}
```

**C Type**: `struct sCommParameters` (in libiec61850 `goose_publisher.h`: `vlanPriority`, `vlanId`, `appId`, `dstAddress[6]`).

**Description**: GOOSE/SV communication parameters (VLAN, APPID, destination MAC). This is the explicit Go equivalent of the C struct. It is embedded in **GoosePublisherConf**; when creating a publisher you can set fields either on the embedded struct or as promoted fields on the conf (e.g. `conf.AppID`, `conf.DstAddr`).

**Example**:
```go
params := iec61850.CommParameters{
    VlanPriority: 4,
    VlanID:       0,
    AppID:        0x1000,
    DstAddr:      [6]uint8{0x01, 0x0c, 0xcd, 0x01, 0x00, 0x01},
}
conf := iec61850.GoosePublisherConf{InterfaceID: "eth0", CommParameters: params}
pub, _ := iec61850.NewGoosePublisher(conf)
```

---

### GooseReceiverSocket

**Go Type**: `type GooseReceiverSocket struct` (opaque)  
**C Type**: `EthernetSocket`

**Description**: Opaque handle returned by GooseReceiver.StartThreadless(). Represents the Ethernet socket used for receiving; drive reception by calling HandleMessage with each received frame.

**Example**:
```go
sock := receiver.StartThreadless()
if sock != nil {
    defer receiver.StopThreadless()
    receiver.HandleMessage(ethernetFrame)
}
```

---

### GooseSubscriber

**Go Type**: `type GooseSubscriber struct` (opaque)  
**C Type**: `GooseSubscriber`

**Description**: GOOSE subscriber instance.

**Example**:
```go
conf := iec61850.SubscriberConf{
    InterfaceID: "eth0",
    AppID:       1000,
}
sub := iec61850.NewGooseSubscriber(conf)
defer sub.Destroy()

sub.SetGooseReceiver(func(s *iec61850.GooseSubscriber) {
    goID := s.GetGoID()
    fmt.Printf("GOOSE from %s\n", goID)
})
sub.Subscribe()
```

---

### GoosePublisher

**Go Type**: `type GoosePublisher struct` (opaque)  
**C Type**: `GoosePublisher`

**Description**: GOOSE publisher instance.

**Example**:
```go
pub, err := iec61850.NewGoosePublisher(conf)
if err != nil {
    log.Fatal(err)
}
defer pub.Destroy()
_ = pub.Publish(values)
```

---

### SubscriberConf

**Go Type**:
```go
type SubscriberConf struct {
    InterfaceID   string           // Network interface (e.g. "eth0"); set on GooseReceiver via SetInterfaceID
    DstMacAddr    [6]uint8         // Destination MAC filter
    AppID         uint16           // APPID filter
    Subscriber    string           // GoCB reference in MMS notation (e.g. "simpleIOGenericIO/LLN0$GO$gcbEvents")
    ReportHandler GooseReportCallback  // Callback when GOOSE message received
}
```

**C Type**: N/A (used in NewGooseSubscriber)

**Description**: GOOSE subscriber configuration. InterfaceID is stored in Conf; set the receiver's interface with `receiver.SetInterfaceID(conf.InterfaceID)`.

**Example**:
```go
conf := iec61850.SubscriberConf{
    InterfaceID:   "eth0",
    DstMacAddr:    [6]uint8{0x01, 0x0c, 0xcd, 0x01, 0x00, 0x01},
    AppID:         1000,
    Subscriber:    "simpleIOGenericIO/LLN0$GO$gcbEvents",
    ReportHandler: func(r *iec61850.GooseReport) { fmt.Println(r.GetGoID()) },
}
sub := iec61850.NewGooseSubscriber(conf)
receiver.SetInterfaceID(conf.InterfaceID)
receiver.AddSubscriber(sub)
```

---

### GoosePublisherConf

**Go Type**:
```go
type GoosePublisherConf struct {
    InterfaceID string
    CommParameters  // embedded: VlanPriority, VlanID, AppID, DstAddr are promoted
}
```

**C Type**: Interface name (e.g. `"eth0"`) plus `struct sCommParameters` (see **CommParameters**). The C API takes `CommParameters*` and `const char* interfaceID`; in Go both are in one struct.

**Description**: GOOSE publisher configuration. Embeds **CommParameters**; you can set `InterfaceID` and the promoted fields (e.g. `AppID`, `VlanID`, `DstAddr`, `VlanPriority`) directly.

**Example**:
```go
conf := iec61850.GoosePublisherConf{
    InterfaceID:  "eth0",
    AppID:       0x1000,
    VlanID:       0,
    VlanPriority: 4,
    DstAddr:      [6]uint8{0x01, 0x0c, 0xcd, 0x01, 0x00, 0x01},
}
pub, _ := iec61850.NewGoosePublisher(conf)
// Optional: pub.SetGoID("MyPublisher"); pub.SetGoCbRef("Device/LLN0$GO$gcb1"); ...
```

---

### ClientGooseControlBlock (opaque)

**C Type**: `ClientGooseControlBlock` (opaque handle in `iec61850_client.h`).

**Go**: There is no public Go handle. The C type is used only inside the bindings. To read or write GOOSE control block values from the client, use:

- **GetGoCBValues** / **GetGoCBValuesAsync** – return or pass **ClientGooseControlBlockValues**
- **SetGoCBValues** / **SetGoCBValuesAsync** – accept **ClientGooseControlBlockValues** and parameters mask

See **ClientGooseControlBlockValues** below and [FUNCTIONS.md](FUNCTIONS.md) for GetGoCBValues, SetGoCBValues, and the async variants.

---

### ClientGooseControlBlockValues

**Go Type**:
```go
type ClientGooseControlBlockValues struct {
    GoEna      bool
    GoID       string
    DatSet     string
    ConfRev    uint32
    NdsComm    bool
    MinTime    uint32
    MaxTime    uint32
    FixedOffs  bool
    DstAddress PhyComAddress
}
```

**C Type**: Values read from / written to `ClientGooseControlBlock` via the client API.

**Description**: Holds GOOSE control block attributes. Use with **GetGoCBValues** / **GetGoCBValuesAsync** (to read) and **SetGoCBValues** / **SetGoCBValuesAsync** (to write). **PhyComAddress** holds destination MAC, VLAN priority, VLAN ID, and APPID.

---

### GooseParseError

**Go Type**: `type GooseParseError int`

**Description**: GOOSE parse error codes (see ENUMS.md). Values: `GooseParseErrorNoError`, `GooseParseErrorUnknownTag`, `GooseParseErrorTagDecode`, etc.

**Example**:
```go
errCode := subscriber.GetParseError()
if errCode != iec61850.GooseParseErrorNoError {
    fmt.Println("GOOSE parse error:", errCode)
}
```

---

## Sampled Values Types

### SvSubscriberConf

**Go Type**:
```go
type SvSubscriberConf struct {
    EthAddr [6]uint8       // Ethernet address filter
    AppID   uint16
    Handler SvReportHandler  // Callback for SV reports
}
```

**Description**: Sampled Values subscriber configuration (used with NewSvSubscriber).

**Example**:
```go
conf := iec61850.SvSubscriberConf{
    EthAddr: [6]uint8{0x01, 0x0c, 0xcd, 0x04, 0x00, 0x01},
    AppID:  12401,
    Handler: func(r *iec61850.SvReport) { fmt.Printf("SV: %s\n", r.GetSvID()) },
}
sub := iec61850.NewSvSubscriber(conf)
receiver.AddSubscriber(sub)
```

---

### SvReceiverConf

**Go Type**:
```go
type SvReceiverConf struct {
    InterfaceID string  // Network interface (e.g. "eth0")
}
```

**Description**: SV receiver configuration (used with NewSvReceiver).

**Example**:
```go
receiver := iec61850.NewSvReceiver(iec61850.SvReceiverConf{InterfaceID: "eth0"})
defer receiver.Stop().Destroy()
receiver.AddSubscriber(svSub).Start()
```

---

### SvSubscriber

**Go Type**: `type SvSubscriber struct` (opaque)  
**C Type**: `SVSubscriber`

**Description**: Sampled values subscriber. Create with NewSvSubscriber(conf), then add to an SvReceiver with AddSubscriber.

**Example**:
```go
sub := iec61850.NewSvSubscriber(iec61850.SvSubscriberConf{
    AppID: 12401,
    Handler: func(report *iec61850.SvReport) { /* handle */ },
})
receiver.AddSubscriber(sub)
```

---

### SVPublisher

**Go Type**: `type SVPublisher struct` (opaque)  
**C Type**: `SVPublisher`

**Description**: Sampled values publisher. Create with NewSVPublisher(conf).

**Example**:
```go
pub, err := iec61850.NewSVPublisher(iec61850.SVPublisherConf{
    InterfaceID: "eth0",
    AppID:       4000,
})
defer pub.Destroy()
_ = pub.PublishSV(asdu)
```

---

### SvPublisherConf

**Go Type**:
```go
type SvPublisherConf struct {
    EtherName    string   // Network interface (e.g. "eth0")
    AppID        uint16
    DstAddr      [6]uint8 // Destination MAC
    VlanID       uint16
    VlanPriority uint8
}
```

**Description**: Sampled Values publisher configuration (used with NewSVPublisher).

**Example**:
```go
conf := iec61850.SvPublisherConf{
    EtherName:    "eth0",
    AppID:        12401,
    DstAddr:      [6]uint8{0x01, 0x0c, 0xcd, 0x04, 0x00, 0x01},
    VlanID:       1,
    VlanPriority: 4,
}
pub := iec61850.NewSVPublisher(conf)
defer pub.Destroy()
```

---

### SVPublisherASDU

**Go Type**:
```go
type SVPublisherASDU struct {
    SmpCnt  uint16
    Data    []int32
}
```

**Description**: SV ASDU (Application Service Data Unit).

**Example**:
```go
asdu := &iec61850.SVPublisherASDU{
    SmpCnt: 0,
    Data:   []int32{100, 200, 300},
}
_ = publisher.PublishSV(asdu)
```

---

## File Service Types

### FileDirectoryEntry

**Go Type**:
```go
type FileDirectoryEntry struct {
    FileName     string
    FileSize     uint32
    LastModified uint64
}
```

**C Type**: N/A (extracted from file directory response)

**Description**: File metadata from directory listing.

**Example**:
```go
entries, moreFollows, _ := client.GetFileDirectoryEx("/config", "")
for _, entry := range entries {
    modTime := time.UnixMilli(int64(entry.LastModified))
    fmt.Printf("%s: %d bytes (modified %s)\n",
        entry.FileName, entry.FileSize, modTime)
}
```

---

## Configuration Types

### AcseAuthenticationMechanism

**Go Type**: `type AcseAuthenticationMechanism int`

**C Type**: `AcseAuthenticationMechanism`

**Description**: ACSE authentication method.

**Values**:
- `ACSE_AUTH_NONE`
- `ACSE_AUTH_PASSWORD`
- `ACSE_AUTH_CERTIFICATE`
- `ACSE_AUTH_TLS`

**Example**:
```go
// Server uses ACSE password auth
if mech == iec61850.ACSE_AUTH_PASSWORD {
    return user == "admin" && pass == "secret"
}
```

---

### MmsVariableAccessAttribute

**Go Type**: `type MmsVariableAccessAttribute int32`

**Description**: Variable access permission.

**Values**:
- `MmsVariableReadOnly`
- `MmsVariableWriteOnly`
- `MmsVariableReadWrite`

**Example**:
```go
if attr == iec61850.MmsVariableReadWrite {
    // variable is readable and writable
}
```

---

### MmsFileAccessAttribute

**Go Type**: `type MmsFileAccessAttribute int32`

**Description**: File access permission bitmask.

**Values**:
- `MmsFileAccessNone` (0)
- `MmsFileRead` (1)
- `MmsFileWrite` (2)
- `MmsFileDelete` (4)

**Example**:
```go
if (fileAttrs & iec61850.MmsFileRead) != 0 {
    // file is readable
}
```

---

## Time & Quality Types

### Timestamp

**Go Type**:
```go
type Timestamp struct {
    // Internal C timestamp
}
```

**C Type**: `Timestamp`

**Description**: IEC 61850 timestamp with time quality.

**Example**:
```go
ts := iec61850.NewTimestamp(time.Now())
fmt.Printf("Milliseconds: %d\n", ts.GetTimeInMs())
fmt.Printf("Leap second known: %v\n", ts.IsLeapSecondKnown())
```

---

### UtcTimeValue

**Go Type**:
```go
type UtcTimeValue struct {
    Milliseconds uint64 // since Unix epoch
    TimeQuality  uint8  // time quality flags
}
```

**C Type**: N/A (extracted from UTCTime)

**Description**: UTC time with quality from IEC 61850.

**Example**:
```go
value, _ := client.Read("Device/LLN0.Beh.t", iec61850.ST)
utcTime := value.(iec61850.UtcTimeValue)
t := time.UnixMilli(int64(utcTime.Milliseconds))
fmt.Printf("Time: %s (quality: 0x%02x)\n", t, utcTime.TimeQuality)
```

---

### Quality

**Go Type**: `type Quality uint16`

**C Type**: `Quality`

**Description**: IEC 61850 quality flags.

**Constants**:
- `QUALITY_VALIDITY_GOOD`
- `QUALITY_VALIDITY_INVALID`
- `QUALITY_DETAIL_OVERFLOW`
- `QUALITY_DETAIL_OLD_DATA`
- `QUALITY_SOURCE_SUBSTITUTED`
- `QUALITY_TEST`
- `QUALITY_OPERATOR_BLOCKED`
- etc.

**Example**:
```go
quality := iec61850.QUALITY_VALIDITY_GOOD | iec61850.QUALITY_TEST
validity := quality.GetValidity()
```

---

### Validity

**Go Type**: `type Validity uint16`

**Description**: Quality validity enumeration.

**Values**:
- `VALIDITY_GOOD`
- `VALIDITY_INVALID`
- `VALIDITY_RESERVED`
- `VALIDITY_QUESTIONABLE`

**Example**:
```go
if quality.Validity == iec61850.VALIDITY_GOOD {
    fmt.Println("Value is good")
}
```

---

### FC (Functional Constraint)

**Go Type**: `type FC int`

**C Type**: `FunctionalConstraint`

**Description**: IEC 61850 functional constraint.

**Values**: (See ENUMS.md for complete list)
- `ST` - Status information
- `MX` - Measurands
- `SP` - Setpoint
- `DC` - Description
- `CF` - Configuration
- etc.

**Example**:
```go
// Read status value
stVal, _ := client.ReadBool("Device/XCBR1.Pos.stVal", iec61850.ST)

// Read measured value
magF, _ := client.ReadFloat32("Device/MMXU1.A.phsA.cVal.mag.f", iec61850.MX)

// Read description
desc, _ := client.ReadString("Device/LLN0.NamPlt.vendor", iec61850.DC)
```

---

*End of Structs Reference*
