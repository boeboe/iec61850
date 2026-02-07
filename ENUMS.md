# IEC 61850 Go Bindings - Enums & Constants Reference

**Version**: 1.6.1  
**Generated**: February 7, 2026

This document provides comprehensive documentation for all exported constants and enumerations in the iec61850 package, including their corresponding C definitions and usage examples. It is aligned with [GAPS.md](GAPS.md) (MMS coverage analysis).

---

## Table of Contents

1. [MMS Data Types](#mms-data-types)
2. [MMS Error Codes](#mms-error-codes)
3. [IED Error Codes](#ied-error-codes)
4. [Functional Constraints](#functional-constraints)
5. [Connection States](#connection-states)
6. [Control Models](#control-models)
7. [Quality Flags](#quality-flags)
8. [Server Types](#server-types)
9. [GOOSE Types](#goose-types)
10. [Authentication & Security](#authentication--security)

---

## MMS Data Types

### MmsType

**Go Type**: `type MmsType int`  
**C Type**: `MmsType`

**Description**: MMS basic data types.

**Values**:

| Constant | Value | C Equivalent | Description |
|----------|-------|-------------|-------------|
| `Array` | 0 | `MMS_ARRAY` | Array of elements |
| `Structure` | 1 | `MMS_STRUCTURE` | Structured type |
| `Boolean` | 2 | `MMS_BOOLEAN` | Boolean value |
| `BitString` | 3 | `MMS_BIT_STRING` | Bit string |
| `Integer` | 4 | `MMS_INTEGER` | Signed integer |
| `Unsigned` | 5 | `MMS_UNSIGNED` | Unsigned integer |
| `Float` | 6 | `MMS_FLOAT` | Floating point |
| `OctetString` | 7 | `MMS_OCTET_STRING` | Octet string |
| `VisibleString` | 8 | `MMS_VISIBLE_STRING` | ASCII string |
| `GeneralizedTime` | 9 | `MMS_GENERALIZED_TIME` | Generalized time |
| `BinaryTime` | 10 | `MMS_BINARY_TIME` | Binary time |
| `Bcd` | 11 | `MMS_BCD` | Binary coded decimal |
| `ObjId` | 12 | `MMS_OBJ_ID` | Object identifier |
| `String` | 13 | `MMS_STRING` | MMS string |
| `UTCTime` | 14 | `MMS_UTC_TIME` | UTC time |
| `DataAccessError` | 15 | `MMS_DATA_ACCESS_ERROR` | Access error |
| `Int8` | 16 | - | 8-bit signed int |
| `Int16` | 17 | - | 16-bit signed int |
| `Int32` | 18 | - | 32-bit signed int |
| `Int64` | 19 | - | 64-bit signed int |
| `Uint8` | 20 | - | 8-bit unsigned int |
| `Uint16` | 21 | - | 16-bit unsigned int |
| `Uint32` | 22 | - | 32-bit unsigned int |

**Example**:
```go
// Create different MMS types
intVal, _ := iec61850.NewMmsValue(iec61850.Integer, int64(42))
boolVal, _ := iec61850.NewMmsValue(iec61850.Boolean, true)
strVal, _ := iec61850.NewMmsValue(iec61850.VisibleString, "Hello")

// Check type
if intVal.Type == iec61850.Integer {
    fmt.Println("It's an integer")
}
```

---

## MMS Error Codes

### MmsDataAccessError

**Go Type**: `type MmsDataAccessError int`  
**C Type**: `MmsDataAccessError`

**Description**: MMS data access error codes.

**Values**:

| Constant | Value | C Equivalent | Description |
|----------|-------|-------------|-------------|
| `DATA_ACCESS_ERROR_SUCCESS_NO_UPDATE` | -3 | - | Success but no update |
| `DATA_ACCESS_ERROR_NO_RESPONSE` | -2 | - | No response received |
| `DATA_ACCESS_ERROR_SUCCESS` | -1 | - | Success |
| `DATA_ACCESS_ERROR_OBJECT_INVALIDATED` | 0 | `DATA_ACCESS_ERROR_OBJECT_INVALIDATED` | Object invalidated |
| `DATA_ACCESS_ERROR_HARDWARE_FAULT` | 1 | `DATA_ACCESS_ERROR_HARDWARE_FAULT` | Hardware fault |
| `DATA_ACCESS_ERROR_TEMPORARILY_UNAVAILABLE` | 2 | `DATA_ACCESS_ERROR_TEMPORARILY_UNAVAILABLE` | Temporarily unavailable |
| `DATA_ACCESS_ERROR_OBJECT_ACCESS_DENIED` | 3 | `DATA_ACCESS_ERROR_OBJECT_ACCESS_DENIED` | Access denied |
| `DATA_ACCESS_ERROR_OBJECT_UNDEFINED` | 4 | `DATA_ACCESS_ERROR_OBJECT_UNDEFINED` | Object not found |
| `DATA_ACCESS_ERROR_INVALID_ADDRESS` | 5 | `DATA_ACCESS_ERROR_INVALID_ADDRESS` | Invalid address |
| `DATA_ACCESS_ERROR_TYPE_UNSUPPORTED` | 6 | `DATA_ACCESS_ERROR_TYPE_UNSUPPORTED` | Type not supported |
| `DATA_ACCESS_ERROR_TYPE_INCONSISTENT` | 7 | `DATA_ACCESS_ERROR_TYPE_INCONSISTENT` | Type mismatch |
| `DATA_ACCESS_ERROR_OBJECT_ATTRIBUTE_INCONSISTENT` | 8 | `DATA_ACCESS_ERROR_OBJECT_ATTRIBUTE_INCONSISTENT` | Attribute inconsistent |
| `DATA_ACCESS_ERROR_OBJECT_ACCESS_UNSUPPORTED` | 9 | `DATA_ACCESS_ERROR_OBJECT_ACCESS_UNSUPPORTED` | Access not supported |
| `DATA_ACCESS_ERROR_OBJECT_NONE_EXISTENT` | 10 | `DATA_ACCESS_ERROR_OBJECT_NONE_EXISTENT` | Object doesn't exist |
| `DATA_ACCESS_ERROR_OBJECT_VALUE_INVALID` | 11 | `DATA_ACCESS_ERROR_OBJECT_VALUE_INVALID` | Invalid value |
| `DATA_ACCESS_ERROR_UNKNOWN` | 12 | `DATA_ACCESS_ERROR_UNKNOWN` | Unknown error |
| `DATA_ACCESS_ERROR_TYPE_CONFLICT` | 13 | `DATA_ACCESS_ERROR_TYPE_CONFLICT` | Type conflict |

**Example**:
```go
_, err := client.Write("Device/GGIO1.AnIn1.mag.f", iec61850.MX, "invalid")
if err != nil {
    // Check specific error
    if err == iec61850.TypeMismatch {
        fmt.Println("Wrong data type!")
    }
}
```

---

## IED Error Codes

### IedClientError

**C Type**: `IedClientError`

**Common Go Error Variables**:

| Go Error | C Constant | Description |
|----------|-----------|-------------|
| `NoError` | `IED_ERROR_OK` | No error |
| `NotConnected` | `IED_ERROR_NOT_CONNECTED` | Not connected |
| `AlreadyConnected` | `IED_ERROR_ALREADY_CONNECTED` | Already connected |
| `ConnectionLost` | `IED_ERROR_CONNECTION_LOST` | Connection lost |
| `ServiceNotSupported` | `IED_ERROR_SERVICE_NOT_SUPPORTED` | Service not supported |
| `ConnectionRejected` | `IED_ERROR_CONNECTION_REJECTED` | Connection rejected |
| `ObjectDoesNotExist` | `IED_ERROR_OBJECT_DOES_NOT_EXIST` | Object not found |
| `ObjectExists` | `IED_ERROR_OBJECT_EXISTS` | Object already exists |
| `ObjectAccessDenied` | `IED_ERROR_OBJECT_ACCESS_DENIED` | Access denied |
| `TypeInconsistent` | `IED_ERROR_TYPE_INCONSISTENT` | Type mismatch |
| `TemporarilyUnavailable` | `IED_ERROR_TEMPORARILY_UNAVAILABLE` | Temporarily unavailable |
| `ObjectUndefined` | `IED_ERROR_OBJECT_UNDEFINED` | Object undefined |
| `InvalidAddress` | `IED_ERROR_INVALID_ADDRESS` | Invalid address |
| `HardwareFault` | `IED_ERROR_HARDWARE_FAULT` | Hardware fault |
| `TypeUnsupported` | `IED_ERROR_TYPE_UNSUPPORTED` | Type not supported |
| `ObjectAttributeInconsistent` | `IED_ERROR_OBJECT_ATTRIBUTE_INCONSISTENT` | Attribute inconsistent |
| `ObjectValueInvalid` | `IED_ERROR_OBJECT_VALUE_INVALID` | Invalid value |
| `ObjectInvalidated` | `IED_ERROR_OBJECT_INVALIDATED` | Object invalidated |
| `Timeout` | `IED_ERROR_TIMEOUT` | Request timeout |
| `FileError` | `IED_ERROR_FILE_ERROR` | File service error |

**Example**:
```go
client, err := iec61850.NewClient(settings)
if err == iec61850.ConnectionRejected {
    log.Fatal("Server rejected connection")
} else if err == iec61850.Timeout {
    log.Fatal("Connection timeout")
}
```

---

## Functional Constraints

### FC

**Go Type**: `type FC int`  
**C Type**: `FunctionalConstraint`

**Description**: IEC 61850 functional constraint classification.

**Values**:

| Constant | Value | Name | Description | Typical Use |
|----------|-------|------|-------------|-------------|
| `ST` | 0 | Status | Status information | stVal, q, t  |
| `MX` | 1 | Measurands | Measured values | mag, ang, q, t |
| `SP` | 2 | Setpoint | Setpoint values | setMag, setVal |
| `SV` | 3 | Substitution | Substitution values | subVal, subQ |
| `CF` | 4 | Configuration | Configuration data | Parameters, settings |
| `DC` | 5 | Description | Device description | d, dU (vendor, model) |
| `SG` | 6 | Setting group | Setting group values | Group parameters |
| `SE` | 7 | Setting group edit | Editable settings | Edit mode data |
| `SR` | 8 | Service response | Service tracking | Last command, error |
| `OR` | 9 | Operate received | Operate tracking | Operation tracking |
| `BL` | 10 | Blocking | Blocking information | Block conditions |
| `EX` | 11 | Extended definition | Extended attributes | Vendor extensions |
| `CO` | 12 | Control | Control (deprecated) | Legacy control |
| `US` | 13 | Unicast SV | Unicast SV | Unicast sampled values |
| `MS` | 14 | Multicast SV | Multicast SV | Multicast sampled values |
| `RP` | 15 | Unbuffered report | Unbuffered reporting | URCB parameters |
| `BR` | 16 | Buffered report | Buffered reporting | BRCB parameters |
| `LG` | 17 | Log control | Log control blocks | Log settings |
| `GO` | 18 | GOOSE control | GOOSE control blocks | GOOSE configuration |
| `ALL` | 99 | All FCs | Wildcard | Read all FCs |
| `NONE` | -1 | None | No FC | Special cases |

**Example**:
```go
// Read different functional constraints
status, _ := client.ReadBool("Device/XCBR1.Pos.stVal", iec61850.ST)
measured, _ := client.ReadFloat32("Device/MMXU1.A.phsA.cVal.mag.f", iec61850.MX)
setpoint, _ := client.ReadFloat32("Device/GAPC1.AnOut.setMag.f", iec61850.SP)
description, _ := client.ReadString("Device/LLN0.NamPlt.vendor", iec61850.DC)
config, _ := client.ReadInt32("Device/MMXU1.ARtg.setMag.i", iec61850.CF)
```

---

## Connection States

### IedConnectionState

**Go Type**: `type IedConnectionState int`  
**C Type**: `IedConnectionState`

**Description**: IED connection state.

**Values**:

| Constant | Value | C Equivalent | Description |
|----------|-------|-------------|-------------|
| `IedStateClosed` | 0 | `IED_STATE_CLOSED` | Connection closed |
| `IedStateConnecting` | 1 | `IED_STATE_CONNECTING` | Connecting in progress |
| `IedStateConnected` | 2 | `IED_STATE_CONNECTED` | Connected |
| `IedStateClosing` | 3 | `IED_STATE_CLOSING` | Closing in progress |

**Example**:
```go
state := client.GetState()
switch state {
case iec61850.IedStateClosed:
    fmt.Println("Not connected")
case iec61850.IedStateConnecting:
    fmt.Println("Connecting...")
case iec61850.IedStateConnected:
    fmt.Println("Connected!")
case iec61850.IedStateClosing:
    fmt.Println("Closing...")
}
```

---

### MmsConnectionState

**C Type**: `MmsConnectionState`

**Values**:
- `MMS_CONNECTION_STATE_IDLE`
- `MMS_CONNECTION_STATE_CONNECTING`
- `MMS_CONNECTION_STATE_CONNECTED`
- `MMS_CONNECTION_STATE_CLOSING`
- `MMS_CONNECTION_STATE_CLOSED`

**Example**:
```go
// State is used internally; after ConnectAsync callback, connection is MMS_CONNECTION_STATE_CONNECTED
if state == iec61850.MMS_CONNECTION_STATE_CONNECTED {
    fmt.Println("MMS association established")
}
```

---

## Control Models

### ControlModel

**Go Type**: `type ControlModel int`  
**C Type**: `ControlModel`

**Description**: IEC 61850 control model.

**Values**:

| Constant | Value | C Equivalent | Description |
|----------|-------|-------------|-------------|
| `CONTROL_MODEL_STATUS_ONLY` | 0 | `CONTROL_MODEL_STATUS_ONLY` | Status only, no control |
| `CONTROL_MODEL_DIRECT_NORMAL` | 1 | `CONTROL_MODEL_DIRECT_NORMAL` | Direct control with normal security |
| `CONTROL_MODEL_SBO_NORMAL` | 2 | `CONTROL_MODEL_SBO_NORMAL` | Select-before-operate with normal security |
| `CONTROL_MODEL_DIRECT_ENHANCED` | 3 | `CONTROL_MODEL_DIRECT_ENHANCED` | Direct control with enhanced security |
| `CONTROL_MODEL_SBO_ENHANCED` | 4 | `CONTROL_MODEL_SBO_ENHANCED` | Select-before-operate with enhanced security |

**Example**:
```go
ctlModel := iec61850.CONTROL_MODEL_SBO_NORMAL
// For SBO model:
err := client.Select("Device/XCBR1.Pos")
if err == nil {
    param := iec61850.ControlObjectParam{CtlVal: true, CtlNum: 1}
    err = client.Operate("Device/XCBR1.Pos", param)
}
```

---

### ControlHandlerResult

**Go Type**: `type ControlHandlerResult int`  
**C Type**: `ControlHandlerResult`

**Description**: Result from server-side control handler.

**Values**:

| Constant | Value | C Equivalent | Description |
|----------|-------|-------------|-------------|
| `CONTROL_RESULT_FAILED` | 0 | `CONTROL_RESULT_FAILED` | Control operation failed |
| `CONTROL_RESULT_OK` | 1 | `CONTROL_RESULT_OK` | Control operation successful |
| `CONTROL_RESULT_WAITING` | 2 | `CONTROL_RESULT_WAITING` | Waiting for completion |

**Example**:
```go
// In server control handler
return iec61850.CONTROL_RESULT_OK
```

---

## Quality Flags

### Quality

**Go Type**: `type Quality uint16`  
**C Type**: `Quality`

**Description**: IEC 61850 quality flags (bitfield).

**Values**:

| Constant | Value | Bits | Description |
|----------|-------|------|-------------|
| `QUALITY_VALIDITY_GOOD` | 0x0000 | 0-1 | Valid data |
| `QUALITY_VALIDITY_INVALID` | 0x0002 | 0-1 | Invalid data |
| `QUALITY_VALIDITY_RESERVED` | 0x0001 | 0-1 | Reserved |
| `QUALITY_VALIDITY_QUESTIONABLE` | 0x0003 | 0-1 | Questionable data |
| `QUALITY_DETAIL_OVERFLOW` | 0x0004 | 2 | Overflow detected |
| `QUALITY_DETAIL_OUT_OF_RANGE` | 0x0008 | 3 | Value out of range |
| `QUALITY_DETAIL_BAD_REFERENCE` | 0x0010 | 4 | Bad reference |
| `QUALITY_DETAIL_OSCILLATORY` | 0x0020 | 5 | Oscillatory condition |
| `QUALITY_DETAIL_FAILURE` | 0x0040 | 6 | Failure detected |
| `QUALITY_DETAIL_OLD_DATA` | 0x0080 | 7 | Old data (not updated) |
| `QUALITY_DETAIL_INCONSISTENT` | 0x0100 | 8 | Inconsistent data |
| `QUALITY_DETAIL_INACCURATE` | 0x0200 | 9 | Inaccurate measurement |
| `QUALITY_SOURCE_SUBSTITUTED` | 0x0400 | 10 | Substituted value |
| `QUALITY_TEST` | 0x0800 | 11 | Test mode |
| `QUALITY_OPERATOR_BLOCKED` | 0x1000 | 12 | Operator blocked |
| `QUALITY_DERIVED` | 0x2000 | 13 | Derived value |

**Example**:
```go
quality := iec61850.QUALITY_VALIDITY_GOOD | iec61850.QUALITY_SOURCE_SUBSTITUTED
validity := quality.GetValidity()

if validity == iec61850.VALIDITY_GOOD {
    fmt.Println("Data is valid")
}

if quality & iec61850.QUALITY_SOURCE_SUBSTITUTED != 0 {
    fmt.Println("Value is substituted")
}

if quality & iec61850.QUALITY_TEST != 0 {
    fmt.Println("Device is in test mode")
}
```

---

### Validity

**Go Type**: `type Validity uint16`  
**C Type**: `Validity`

**Description**: Quality validity enumeration.

**Values**:

| Constant | Value | C Equivalent | Description |
|----------|-------|-------------|-------------|
| `VALIDITY_GOOD` | 0 | `VALIDITY_GOOD` | Data is valid |
| `VALIDITY_INVALID` | 1 | `VALIDITY_INVALID` | Data is invalid |
| `VALIDITY_RESERVED` | 2 | `VALIDITY_RESERVED` | Reserved |
| `VALIDITY_QUESTIONABLE` | 3 | `VALIDITY_QUESTIONABLE` | Data quality is questionable |

**Example**:
```go
if quality.GetValidity() == iec61850.VALIDITY_GOOD {
    fmt.Println("Value is valid")
}
```

---

## Server Types

### Edition

**Values for ServerConfig**:
```go
const (
    IEC_61850_EDITION_1   = 0
    IEC_61850_EDITION_2   = 1
    IEC_61850_EDITION_2_1 = 2
)
```

**Example**:
```go
config := iec61850.ServerConfig{
    Edition: iec61850.IEC_61850_EDITION_2,
}
```

---

### MmsServerState

**Go Type**: `type MmsServerState int32`  
**C Type**: `MmsServerState`

**Values**:

| Constant | Value | Description |
|----------|-------|-------------|
| `MmsServerStateIdle` | 0 | Server idle |
| `MmsServerStateLoading` | 1 | Server loading |
| `MmsServerStateRunning` | 2 | Server running |

Aliases:
- `MMS_SERVER_STATE_IDLE` = `MmsServerStateIdle`
- `MMS_SERVER_STATE_LOADING` = `MmsServerStateLoading`
- `MMS_SERVER_STATE_RUNNING` = `MmsServerStateRunning`

**Example**:
```go
if serverState == iec61850.MmsServerStateRunning {
    fmt.Println("Server is running")
}
```

---

### MmsServerConnectionState

**Go Type**: `type MmsServerConnectionState int32`  
**C Type**: `MmsConnectionState`

**Description**: State of an MMS server connection.

**Values**:

| Constant | Value | Description |
|----------|-------|-------------|
| `MmsServerConnectionStateIdle` | 0 | Idle |
| `MmsServerConnectionStateAssociation` | 1 | Association established |
| `MmsServerConnectionStateConcluded` | 2 | Concluded |

Aliases: `MMS_CON_STATE_IDLE`, `MMS_CON_STATE_ASSOCIATION`, `MMS_CON_STATE_CONCLUDED`

**Example**:
```go
// Check connection state after MMS connect
params := mmsConn.GetMmsConnectionParameters()
// State is typically reflected via connection lifecycle; use Idle/Association/Concluded as needed
var state iec61850.MmsServerConnectionState = iec61850.MmsServerConnectionStateConcluded
_ = state
```

---

### MmsServiceSupportedBitmap

**Go Type**: `type MmsServiceSupportedBitmap uint32`  
**C Type**: Bitmask

**Description**: Bitmap of MMS services supported by the server (from GetConnectionParameters / MmsConnectionParameters).

**Values** (selected): `MmsServiceStatus`, `MmsServiceGetNameList`, `MmsServiceIdentify`, `MmsServiceRead`, `MmsServiceWrite`, `MmsServiceGetVariableAccess`, `MmsServiceDefineNamedVariableList`, `MmsServiceGetNamedVariableListAttrs`, `MmsServiceDeleteNamedVariableList`, `MmsServiceFileOpen`, `MmsServiceFileRead`, `MmsServiceFileClose`, `MmsServiceFileDelete`, `MmsServiceFileDirectory`, `MmsServiceJournalRead` (and `MMS_SERVICE_*` aliases).

**Example**:
```go
params := mmsConn.GetMmsConnectionParameters()
// Check if server supports file directory (bit in ServicesSupported or via GetConnectionParameters)
if params != nil {
    // Use params.ServicesSupported or server status to determine supported services
    _ = params.MaxPduSize
}
```

---

### MMS File Service Types

**Go Type**: `type MmsFileServiceType int`  
**C Type**: `MmsFileServiceType` (MMS_FILE_ACCESS_TYPE_*)

**Description**: Type of MMS file service requested in a file access handler.

**Values**:

| Constant | Value | C Equivalent | Description |
|----------|-------|--------------|-------------|
| `MmsFileAccessReadDirectory` | 0 | `MMS_FILE_ACCESS_TYPE_READ_DIRECTORY` | Read directory |
| `MmsFileAccessOpen` | 1 | `MMS_FILE_ACCESS_TYPE_OPEN` | Open file |
| `MmsFileAccessObtain` | 2 | `MMS_FILE_ACCESS_TYPE_OBTAIN` | Obtain (upload) file |
| `MmsFileAccessDelete` | 3 | `MMS_FILE_ACCESS_TYPE_DELETE` | Delete file |
| `MmsFileAccessRename` | 4 | `MMS_FILE_ACCESS_TYPE_RENAME` | Rename file |

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

### MMS Get Name List Types

**Go Type**: `type MmsGetNameListType int`  
**C Type**: `MmsGetNameListType` (MMS_GETNAMELIST_*)

**Description**: Type of GetNameList request in InstallGetNameListHandler.

**Values**:

| Constant | Value | Description |
|----------|-------|-------------|
| `MmsGetNameListDomains` | 0 | Domain names |
| `MmsGetNameListJournals` | 1 | Journal names |
| `MmsGetNameListDataSets` | 2 | Data set names |
| `MmsGetNameListData` | 3 | Data (variable) names |

**Example**:
```go
server.InstallGetNameListHandler(func(nameListType iec61850.MmsGetNameListType, domainID string) bool {
    if nameListType == iec61850.MmsGetNameListData {
        return true // allow variable name list
    }
    return true
})
```

---

### Variable List Access Types

**Go Type**: `type MmsVariableListAccessType int`

**Description**: Kind of named variable list access in VariableListAccessHandler.

**Values**:

| Constant | Value | Description |
|----------|-------|-------------|
| `MmsVarlistCreate` | 0 | Create list |
| `MmsVarlistDelete` | 1 | Delete list |
| `MmsVarlistRead` | 2 | Read list |
| `MmsVarlistWrite` | 3 | Write list |
| `MmsVarlistGetDirectory` | 4 | Get directory |

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

### Variable List Types

**Go Type**: `type MmsVariableListType int`

**Description**: Scope of the named variable list (domain-, association-, or VMD-specific).

**Values**:

| Constant | Value | Description |
|----------|-------|-------------|
| `MmsVarlistTypeDomainSpecific` | 0 | Domain-specific |
| `MmsVarlistTypeAssociationSpecific` | 1 | Association-specific |
| `MmsVarlistTypeVmdSpecific` | 2 | VMD-specific |

**Example**:
```go
server.InstallVariableListAccessHandler(func(accessType iec61850.MmsVariableListAccessType, listType iec61850.MmsVariableListType, domainID, listName string) error {
    switch listType {
    case iec61850.MmsVarlistTypeDomainSpecific:
        return nil
    case iec61850.MmsVarlistTypeAssociationSpecific:
        return nil
    case iec61850.MmsVarlistTypeVmdSpecific:
        return nil
    }
    return nil
})
```

---

## GOOSE Types

### GooseParseError

**Go Type**: `type GooseParseError int`  
**C Type**: `GooseParseError`

**Description**: Parse error code when GOOSE message parsing fails (GooseSubscriber_isValid false).

**Values**:

| Constant | Value | Description |
|----------|-------|-------------|
| `GooseParseErrorNoError` | 0 | No error |
| `GooseParseErrorUnknownTag` | 1 | Unknown tag |
| `GooseParseErrorTagDecode` | 2 | Tag decode error |
| `GooseParseErrorSublevel` | 3 | Sublevel error |
| `GooseParseErrorOverflow` | 4 | Overflow |
| `GooseParseErrorUnderflow` | 5 | Underflow |
| `GooseParseErrorTypeMismatch` | 6 | Type mismatch |
| `GooseParseErrorLengthMismatch` | 7 | Length mismatch |
| `GooseParseErrorInvalidPadding` | 8 | Invalid padding |

**Example**:
```go
if !sub.IsValid() {
    parseErr := sub.GetParseError()
    if parseErr != iec61850.GooseParseErrorNoError {
        fmt.Printf("Parse error: %d\n", parseErr)
    }
}
```

---

### GoCB element constants (parameters mask)

**Go Type**: `const` (e.g. `GoCBElementGoEna`, `GoCBElementAll`)  
**C Type**: `GOCB_ELEMENT_*` macros in `iec61850_client.h`

**Description**: Bit masks for GetGoCBValues/SetGoCBValues `parametersMask` to select which GoCB elements to read or write.

**Values**:

| Constant | Value | C Equivalent | Description |
|----------|-------|-------------|-------------|
| `GoCBElementGoEna` | 1 | `GOCB_ELEMENT_GO_ENA` | Go enable |
| `GoCBElementGoID` | 2 | `GOCB_ELEMENT_GO_ID` | GoID |
| `GoCBElementDatSet` | 4 | `GOCB_ELEMENT_DATSET` | Data set reference |
| `GoCBElementConfRev` | 8 | `GOCB_ELEMENT_CONF_REV` | Configuration revision |
| `GoCBElementNdsComm` | 16 | `GOCB_ELEMENT_NDS_COMM` | Needs commission |
| `GoCBElementDstAddress` | 32 | `GOCB_ELEMENT_DST_ADDRESS` | Destination address |
| `GoCBElementMinTime` | 64 | `GOCB_ELEMENT_MIN_TIME` | Min time |
| `GoCBElementMaxTime` | 128 | `GOCB_ELEMENT_MAX_TIME` | Max time |
| `GoCBElementFixedOffs` | 256 | `GOCB_ELEMENT_FIXED_OFFS` | Fixed offset |
| `GoCBElementAll` | 511 | `GOCB_ELEMENT_ALL` | All elements |

**Example**:
```go
vals, err := client.GetGoCBValues("Device/LLN0$GO$gcb1")
// Set only destination address
err = client.SetGoCBValues("Device/LLN0$GO$gcb1", newVals, iec61850.GoCBElementDstAddress, false)
```

---

## Authentication & Security

### AcseAuthenticationMechanism

**Go Type**: `type AcseAuthenticationMechanism int`  
**C Type**: `AcseAuthenticationMechanism`

**Description**: ACSE authentication mechanism.

**Values**:

| Constant | Value | C Equivalent | Description |
|----------|-------|-------------|-------------|
| `ACSE_AUTH_NONE` | 0 | `ACSE_AUTH_NONE` | No authentication |
| `ACSE_AUTH_PASSWORD` | 1 | `ACSE_AUTH_PASSWORD` | Password authentication |
| `ACSE_AUTH_CERTIFICATE` | 2 | `ACSE_AUTH_CERTIFICATE` | Certificate authentication |
| `ACSE_AUTH_TLS` | 3 | `ACSE_AUTH_TLS` | TLS certificate authentication |

**Example**:
```go
if mech == iec61850.ACSE_AUTH_PASSWORD {
    return checkPassword(user, pass)
}
```

---

### MmsVariableAccessAttribute

**Go Type**: `type MmsVariableAccessAttribute int32`  
**C Type**: `MmsVariableAccessAttribute`

**Description**: Variable access permissions.

**Values**:

| Constant | Value | Description |
|----------|-------|-------------|
| `MmsVariableReadOnly` | 0 | Read-only access |
| `MmsVariableWriteOnly` | 1 | Write-only access |
| `MmsVariableReadWrite` | 2 | Read-write access |

Aliases:
- `MMS_VARIABLE_READ_ONLY` = `MmsVariableReadOnly`
- `MMS_VARIABLE_WRITE_ONLY` = `MmsVariableWriteOnly`
- `MMS_VARIABLE_READ_WRITE` = `MmsVariableReadWrite`
- `MMS_ACCESS_READ_WRITE` = `MmsVariableReadWrite`
- `MMS_ACCESS_READ_ONLY` = `MmsVariableReadOnly`
- `MMS_ACCESS_WRITE_ONLY` = `MmsVariableWriteOnly`

**Example**:
```go
if attr == iec61850.MmsVariableReadWrite {
    // variable supports read and write
}
```

---

### MmsFileAccessAttribute

**Go Type**: `type MmsFileAccessAttribute int32`  
**C Type**: N/A (bitfield)

**Description**: File access permission flags (bitfield).

**Values**:

| Constant | Value | Description |
|----------|-------|-------------|
| `MmsFileAccessNone` | 0 | No access |
| `MmsFileRead` | 1 | Read permission |
| `MmsFileWrite` | 2 | Write permission |
| `MmsFileDelete` | 4 | Delete permission |

Aliases:
- `MMS_FILE_ACCESS_NONE` = `MmsFileAccessNone`
- `MMS_FILE_READ` = `MmsFileRead`
- `MMS_FILE_WRITE` = `MmsFileWrite`
- `MMS_FILE_DELETE` = `MmsFileDelete`

**Example**:
```go
// Grant read and write permissions
perms := iec61850.MmsFileRead | iec61850.MmsFileWrite

// Check permissions
if perms & iec61850.MmsFileDelete != 0 {
    fmt.Println("Delete permitted")
}
```

---

## Named Variable List Types

### MmsDeletableType

**Go Type**: `type MmsDeletableType int32`  
**C Type**: `MmsDeletable`

**Description**: Indicates deletability of named variable lists.

**Values**:

| Constant | Value | Description |
|----------|-------|-------------|
| `MmsDeletableNone` | 0 | Not deletable |
| `MmsDeletableAASpecific` | 1 | AA-specific deletable |
| `MmsDeletableDomainSpecific` | 2 | Domain-specific deletable |
| `MmsDeletableVMDSpecific` | 3 | VMD-specific deletable |

Aliases:
- `MMS_DELETABLE_NONE` = `MmsDeletableNone`
- `MMS_DELETABLE_AA_SPECIFIC` = `MmsDeletableAASpecific`
- `MMS_DELETABLE_DOMAIN_SPECIFIC` = `MmsDeletableDomainSpecific`
- `MMS_DELETABLE_VMD_SPECIFIC` = `MmsDeletableVMDSpecific`
- `MMS_DELETABLE_NOT` = `MmsDeletableNone`
- `MMS_DELETABLE_AA` = `MmsDeletableAASpecific`
- `MMS_DELETABLE_DOMAIN` = `MmsDeletableDomainSpecific`
- `MMS_DELETABLE_VMD` = `MmsDeletableVMDSpecific`

**Example**:
```go
attrs, _ := mmsConn.GetNamedVariableListAttributes("domain", "list1")
if attrs.Deletable == iec61850.MmsDeletableDomainSpecific {
    // list can be deleted at domain scope
}
```

---

### MmsNamedVariableListType

**Go Type**: `type MmsNamedVariableListType int32`  
**C Type**: `MmsNamedVariableListType`

**Description**: Named variable list scope.

**Values**:

| Constant | Value | Description |
|----------|-------|-------------|
| `NamedVariableListTypeVMDSpecific` | 0 | VMD-specific |
| `NamedVariableListTypeDomainSpecific` | 1 | Domain-specific |
| `NamedVariableListTypeAssociationSpecific` | 2 | Association-specific |

Aliases:
- `NAMED_VARIABLE_LIST_TYPE_VMD_SPECIFIC`
- `NAMED_VARIABLE_LIST_TYPE_DOMAIN_SPECIFIC`
- `NAMED_VARIABLE_LIST_TYPE_ASSOCIATION_SPECIFIC`

**Example**:
```go
if listType == iec61850.NamedVariableListTypeAssociationSpecific {
    // list is association-specific
}
```

---

## Journal Types

### MmsJournalVariable

**Go Type**: `type MmsJournalVariable int32`  
**C Type**: `MmsJournalVariable`

**Description**: Journal variable identifier type.

**Values**:

| Constant | Value | Description |
|----------|-------|-------------|
| `MmsJournalVariableTag` | 0 | Variable tag |
| `MmsJournalVariableEntryID` | 1 | Entry ID |

Aliases:
- `MMS_JOURNAL_VARIABLE_TAG` = `MmsJournalVariableTag`
- `MMS_JOURNAL_VARIABLE_ENTRY_ID` = `MmsJournalVariableEntryID`

**Example**:
```go
// Used when reading journal structure
if jvarType == iec61850.MmsJournalVariableTag {
    // variable identified by tag
}
```

---

## Reporting Constants

### Reason Codes

**C Type**: `ReasonForInclusion`

**Common values** (used in report callbacks):
- `IEC61850_REASON_DATA_CHANGE` (1)
- `IEC61850_REASON_QUALITY_CHANGE` (2)
- `IEC61850_REASON_DATA_UPDATE` (4)
- `IEC61850_REASON_INTEGRITY` (8)
- `IEC61850_REASON_GI` (16)
- `IEC61850_REASON_APPLICATION_TRIGGER` (32)

**Example**:
```go
// In report callback, check reason for inclusion
for _, reason := range report.ReasonCodes {
    if reason == iec61850.IEC61850_REASON_DATA_CHANGE {
        fmt.Println("Report due to data change")
    }
}
```

---

## Setting Group Constants

### Edit State

**Go Type**: `type EditState int`

**Values**:
```go
const (
    EditStateOff EditState = iota
    EditStateOn
    EditStateReserved
)
```

**Example**:
```go
editState, err := client.GetEditSGValue("Device/LLN0$SG$1", "editState")
if editState == iec61850.EditStateOn {
    fmt.Println("Edit mode is active")
}
```

---

## Client Directory Object Types

**Go Type**: `type IecDirectoryCategory int`  
**C Type**: `IecDirectoryCategory`

**Values**:
```go
const (
    IEC61850_OBJECT_TYPE_OTHER IecDirectoryCategory = iota
    IEC61850_OBJECT_TYPE_LD
    IEC61850_OBJECT_TYPE_LN
    IEC61850_OBJECT_TYPE_DO
    IEC61850_OBJECT_TYPE_DA
)
```

---

## Usage Patterns

### Combining Quality Flags

```go
// Build quality value
quality := iec61850.QUALITY_VALIDITY_GOOD | 
            iec61850.QUALITY_SOURCE_SUBSTITUTED |
            iec61850.QUALITY_TEST

// Check individual flags
if quality & iec61850.QUALITY_TEST != 0 {
    fmt.Println("Test mode active")
}

// Extract validity
validity := quality.GetValidity()
if validity == iec61850.VALIDITY_QUESTIONABLE {
    fmt.Println("Data quality is questionable")
}
```

---

### FC-Based Operations

```go
// Read different functional constraints from same object
stVal, _ := client.ReadBool("Device/XCBR1.Pos.stVal", iec61850.ST)      // Status
q, _ := client.ReadInt32("Device/XCBR1.Pos.q", iec61850.ST)             // Quality
t, _ := client.Read("Device/XCBR1.Pos.t", iec61850.ST)                 // Timestamp
vendor, _ := client.ReadString("Device/XCBR1.NamPlt.vendor", iec61850.DC) // Description
pulseConfig, _ := client.ReadInt32("Device/XCBR1.pulseConfig", iec61850.CF) // Configuration
```

---

### Error Handling Patterns

```go
value, err := client.Read("NonExistent/Object", iec61850.ST)
if err != nil {
    switch err {
    case iec61850.ObjectDoesNotExist:
        log.Println("Object not found")
    case iec61850.ObjectAccessDenied:
        log.Println("Access denied")
    case iec61850.TypeInconsistent:
        log.Println("Type mismatch")
    case iec61850.Timeout:
        log.Println("Request timeout")
    default:
        log.Printf("Other error: %v", err)
    }
}
```

---

*End of Enums & Constants Reference*
