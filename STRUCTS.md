# IEC 61850 Go Bindings - Structs Reference

**Version**: 1.6.1  
**Generated**: February 7, 2026

This document provides comprehensive documentation for all exported Go structs/types in the iec61850 package, including their corresponding C structures, field descriptions, and usage examples.

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

**Example**:
```go
server.SetConnectionIndicationHandler(func(conn *iec61850.ClientConnection, connected bool) {
    if connected {
        addr := conn.GetPeerAddress()
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
    LocalTSelector    []byte
    LocalSSelector    []byte
    LocalPSelector    []byte
    RemoteTSelector   []byte
    RemoteSSelector   []byte
    RemotePSelector   []byte
    LocalAeQualifier  int32
    RemoteAeQualifier int32
    LocalApTitle      []byte
    RemoteApTitle     []byte
}
```

**C Type**: `IsoConnectionParameters`

**Description**: ISO connection layer parameters (OSI layer selectors).

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
}
```

**C Type**: N/A (extracted from callbacks)

**Description**: Extended file directory entry with metadata.

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

---

### ControlHandlerResult

**Go Type**: `type ControlHandlerResult int`

**C Type**: `ControlHandlerResult`

**Description**: Result from control handler callback.

**Values**:
- `CONTROL_RESULT_FAILED`
- `CONTROL_RESULT_OK`
- `CONTROL_RESULT_WAITING`

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

---

### SubscriberConf

**Go Type**:
```go
type SubscriberConf struct {
    InterfaceID string
    AppID       int32
    GoID        string
}
```

**C Type**: N/A (used in NewGooseSubscriber)

**Description**: GOOSE subscriber configuration.

---

### GoosePublisherConf

**Go Type**:
```go
type GoosePublisherConf struct {
    InterfaceID string
    AppID       int32
    GoID        string
    GoCbRef     string
    DataSetRef  string
    ConfRev     uint32
}
```

**C Type**: N/A

**Description**: GOOSE publisher configuration.

---

### GooseParseError

**Go Type**: `type GooseParseError int`

**Description**: GOOSE parsing error codes.

**Values**:
- `GOOSE_PARSE_ERROR_NO_ERROR`
- `GOOSE_PARSE_ERROR_INVALID_LENGTH`
- `GOOSE_PARSE_ERROR_APDU_ERROR`
- etc.

---

## Sampled Values Types

### SVSubscriber

**Go Type**: `type SVSubscriber struct` (opaque)  
**C Type**: `SVReceiver`

**Description**: Sampled values subscriber.

**Example**:
```go
svSub := iec61850.NewSVSubscriber("eth0")
defer svSub.Destroy()

svSub.SetSVReceiver(func(sub *iec61850.SVSubscriber, appID int, data []byte) {
    fmt.Printf("SV AppID %d: %d bytes\n", appID, len(data))
})
svSub.Subscribe()
```

---

### SVPublisher

**Go Type**: `type SVPublisher struct` (opaque)  
**C Type**: `SVPublisher`

**Description**: Sampled values publisher.

---

### SVPublisherConf

**Go Type**:
```go
type SVPublisherConf struct {
    InterfaceID string
    AppID       int32
    SvID        string
    DstAddress  [6]byte
    VLANID      uint16
    VLANPriority uint8
}
```

**Description**: SV publisher configuration.

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

---

### MmsVariableAccessAttribute

**Go Type**: `type MmsVariableAccessAttribute int32`

**Description**: Variable access permission.

**Values**:
- `MmsVariableReadOnly`
- `MmsVariableWriteOnly`
- `MmsVariableReadWrite`

---

### MmsFileAccessAttribute

**Go Type**: `type MmsFileAccessAttribute int32`

**Description**: File access permission bitmask.

**Values**:
- `MmsFileAccessNone` (0)
- `MmsFileRead` (1)
- `MmsFileWrite` (2)
- `MmsFileDelete` (4)

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
