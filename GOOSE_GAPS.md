# IEC 61850 GOOSE – Coverage Analysis and Plan

**Analysis Date**: February 7, 2026  
**Scope**: libiec61850 GOOSE C API vs iec61850 Go CGO bindings (functions, structs, enums, types).

This document provides a gap analysis of GOOSE coverage and a plan to increase it. It complements [GAPS.md](GAPS.md) (MMS coverage).

**Platform support:** GOOSE (and SV) bindings are built when the build tag matches: `(linux && (amd64 || arm64 || arm)) || (darwin && (amd64 || arm64)) || (windows && amd64)`. libiec61850 provides an Ethernet HAL for Linux (`hal/ethernet/linux`), Windows (`hal/ethernet/win32`, L2 GOOSE may need winpcap), and BSD/macOS (`hal/ethernet/bsd`). Prebuilt libs for darwin and windows must be built with GOOSE/Ethernet support (CONFIG_INCLUDE_GOOSE_SUPPORT and CONFIG_IEC61850_L2_GOOSE in CMake).

---

## Executive Summary

| Category | C API items | Go implemented | Coverage | Remaining gaps |
|----------|-------------|----------------|----------|----------------|
| **GoosePublisher** | 17 (Ethernet) | 16 | **~94%** | createRemote (R-GOOSE; out of scope) |
| **GooseReceiver** | 14 (Ethernet) | 14 | **100%** | — |
| **GooseSubscriber** | 19 | 19 | **100%** | — |
| **Client GoCB** | 4 service APIs | 4 | **100%** | — |
| **Server GOOSE** | 6 | 4 | **~67%** | SetGooseInterfaceIdEx, UseGooseVlanTag (per-GCB) |
| **R-GOOSE** | 2 (createRemote) | 0 | **0%** | Out of scope (RSession not bound) |
| **Structs / enums / types** | 8 | 8 | **100%** | — |

**Overall GOOSE coverage (Ethernet GOOSE, by item count): ~95%.** Structs/enums/types: CommParameters (named Go type), PhyComAddress, ClientGooseControlBlock (opaque; use ClientGooseControlBlockValues), publisher/receiver/subscriber opaques, GooseParseError, GOCB constants, GooseListener, GetGoCBValuesHandler—all covered; STRUCTS.md documents CommParameters, GoosePublisherConf, and ClientGooseControlBlock. Production ready for publisher, receiver, subscriber, and client GoCB. Server GOOSE: enable/disable/interface implemented; per-GCB interface/VLAN optional.

**API docs**: GOOSE is documented in [FUNCTIONS.md](FUNCTIONS.md), [STRUCTS.md](STRUCTS.md), [ENUMS.md](ENUMS.md).

---

## Part 1: Function Coverage

### 1.1 GoosePublisher (goose_publisher.h)

| C function / type | Go binding | Status |
|-------------------|------------|--------|
| `CommParameters` | `CommParameters` (named type; embedded in `GoosePublisherConf`) | ✅ |
| `GoosePublisher_create(parameters, interfaceID)` | `NewGoosePublisher(conf)` | ✅ |
| `GoosePublisher_createEx(parameters, interfaceID, useVlanTag)` | `NewGoosePublisherEx(conf, useVlanTag)` | ✅ |
| `GoosePublisher_createRemote(RSession, appId)` | — | ❌ Out of scope (R-GOOSE) |
| `GoosePublisher_destroy` | `Close()` | ✅ |
| `GoosePublisher_publish(self, LinkedList)` | `Publish(dataSet *LinkedListValue)` | ✅ |
| `GoosePublisher_publishAndDump(...)` | `PublishAndDump(dataSet, msgBuf)` | ✅ |
| `GoosePublisher_setGoID` | `SetGoID(goID)` | ✅ |
| `GoosePublisher_setGoCbRef` | `SetGoCbRef` | ✅ |
| `GoosePublisher_setTimeAllowedToLive` | `SetTimeAllowedToLive` | ✅ |
| `GoosePublisher_setDataSetRef` | `SetDataSetRef` | ✅ |
| `GoosePublisher_setConfRev` | `SetConfRev` | ✅ |
| `GoosePublisher_setSimulation` | `SetSimulation` | ✅ |
| `GoosePublisher_setStNum` / `setSqNum` | `SetStNum` / `SetSqNum` | ✅ |
| `GoosePublisher_setNeedsCommission` | `SetNeedsCommission` | ✅ |
| `GoosePublisher_increaseStNum` | `IncreaseStNum` | ✅ |
| `GoosePublisher_reset` | `Reset` | ✅ |

**Coverage: 16/17 Ethernet (createRemote out of scope).**

---

### 1.2 GooseReceiver (goose_receiver.h)

| C function | Go binding | Status |
|------------|------------|--------|
| `GooseReceiver_create()` | `NewGooseReceiver()` | ✅ |
| `GooseReceiver_createEx(buffer)` | `NewGooseReceiverEx(buffer)` | ✅ |
| `GooseReceiver_createRemote(RSession)` | — | ❌ Out of scope (R-GOOSE) |
| `GooseReceiver_setInterfaceId` | `SetInterfaceID` | ✅ |
| `GooseReceiver_getInterfaceId` | `GetInterfaceID` | ✅ |
| `GooseReceiver_addSubscriber` | `AddSubscriber` | ✅ |
| `GooseReceiver_removeSubscriber` | `RemoveSubscriber` | ✅ |
| `GooseReceiver_start` | `Start` | ✅ |
| `GooseReceiver_stop` | `Stop` | ✅ |
| `GooseReceiver_isRunning` | `IsRunning` | ✅ |
| `GooseReceiver_destroy` | `Destroy` | ✅ |
| `GooseReceiver_startThreadless` | `StartThreadless()` | ✅ |
| `GooseReceiver_stopThreadless` | `StopThreadless()` | ✅ |
| `GooseReceiver_tick` | `Tick` | ✅ |
| `GooseReceiver_handleMessage(buffer, size)` | `HandleMessage(buffer)` | ✅ |

**Coverage: 14/14 Ethernet (createRemote out of scope).**

---

### 1.3 GooseSubscriber (goose_subscriber.h)

| C function | Go binding | Status |
|------------|------------|--------|
| `GooseSubscriber_create(goCbRef, dataSetValues)` | `NewGooseSubscriber(conf)` or `NewGooseSubscriberWithDataSet(conf, dataSetValues)` | ✅ |
| `GooseSubscriber_setListener` | Via receiver callback bridge | ✅ |
| `GooseSubscriber_setDstMac` / `setAppId` | In `NewGooseSubscriber` from `SubscriberConf` | ✅ |
| All getters (getGoId … getVlanPrio) | GetGoID, GetStNum, GetTimestamp, etc. | ✅ |
| `GooseSubscriber_getDataSetValues` | `GetDataSetValues` | ✅ |
| `GooseSubscriber_destroy` | `Destroy` | ✅ |
| `GooseSubscriber_setObserver` | `SetObserver()` | ✅ |

**Coverage: 19/19 (100%).**

---

### 1.4 Client GoCB (iec61850_client.h)

| C function | Go binding | Status |
|------------|------------|--------|
| `IedConnection_getGoCBValues` | `GetGoCBValues(goCBReference)` | ✅ |
| `IedConnection_getGoCBValuesAsync` | `GetGoCBValuesAsync(goCBReference, callback)` | ✅ |
| `IedConnection_setGoCBValues` | `SetGoCBValues(..., parametersMask, singleRequest)` | ✅ |
| `IedConnection_setGoCBValuesAsync` | `SetGoCBValuesAsync(goCBReference, values, ...)` | ✅ |

**Coverage: 4/4 (100%).**

---

### 1.5 Server GOOSE (iec61850_server.h)

| C function / config | Go binding | Status |
|---------------------|------------|--------|
| `IedServerConfig_useIntegratedGoosePublisher` | `ServerConfig.UseIntegratedGoosePublisher` | ✅ |
| `IedServer_enableGoosePublishing` | `EnableGoosePublishing()` | ✅ |
| `IedServer_disableGoosePublishing` | `DisableGoosePublishing()` | ✅ |
| `IedServer_setGooseInterfaceId` | `SetGooseInterfaceId(interfaceId)` | ✅ |
| `IedServer_setGooseInterfaceIdEx` | — | ❌ Missing (per-GCB; needs LogicalNode) |
| `IedServer_useGooseVlanTag` | — | ❌ Missing (per-GCB; needs LogicalNode) |

**Coverage: 4/6 (≈67%).**

---

## Part 2: Structs, Enums, and Types

The following eight items are the GOOSE-related structs, enums, and callback types. All are covered in the bindings and in STRUCTS.md/ENUMS.md where applicable.

### 2.1 Structs

| C type | Go type | Status | Notes |
|--------|---------|--------|-------|
| `CommParameters` (vlanPriority, vlanId, appId, dstAddress[6]) | `CommParameters` (VlanPriority, VlanID, AppID, DstAddr) | ✅ | Named type; embedded in `GoosePublisherConf`. C layout equivalent. |
| `PhyComAddress` (vlanPriority, vlanId, appId, dstAddress[6]) | `PhyComAddress` (Addr, VlanPriority, VlanId, AppId) | ✅ | Used in client_gocb and SCL. |
| `ClientGooseControlBlock` (opaque) | No public handle; values via `ClientGooseControlBlockValues` | ✅ | Documented in STRUCTS.md: use GetGoCBValues/SetGoCBValues and ClientGooseControlBlockValues. |
| `sGoosePublisher` / `sGooseReceiver` / `sGooseSubscriber` | `GoosePublisher` / `GooseReceiver` / `GooseSubscriber` (opaque) | ✅ | Internal C pointers, no struct field exposure. |

**Note:** `GoosePublisherConf` is the Go config struct; it embeds `CommParameters` and adds `InterfaceID`. STRUCTS.md documents the C↔Go mapping.

---

### 2.2 Enums and constants

| C enum / constant | Go binding | Status |
|-------------------|------------|--------|
| `GooseParseError` (NO_ERROR, UNKNOWN_TAG, …) | `GooseParseError` + `GooseParseErrorNoError`, etc. | ✅ |
| `GOCB_ELEMENT_GO_ENA` … `GOCB_ELEMENT_ALL` | `GoCBElementGoEna` … `GoCBElementAll` | ✅ |

**Coverage: complete for public GOOSE enums/constants used in the bindings.**

---

### 2.3 Callback and handler types

| C type | Go binding | Status |
|--------|------------|--------|
| `GooseListener(GooseSubscriber, void* parameter)` | `GooseReportCallback func(report *GooseReport)` via bridge | ✅ |
| `IedConnection_GetGoCBValuesHandler` | Bridged via `GetGoCBValuesAsync(..., callback)` | ✅ |

---

## Part 3: Plan to Increase Coverage

### Phase 1 – Low-effort, high-value (recommended first)

1. **GoosePublisher**
   - Add **`SetGoID(goID string)`** – single C call, needed for correct GoID in messages when different from GoCbRef.

2. **Server GOOSE**
   - Add **`EnableGoosePublishing()`** and **`DisableGoosePublishing()`** on `IedServer` – wrap `IedServer_enableGoosePublishing` / `IedServer_disableGoosePublishing`.
   - Add **`SetGooseInterfaceId(interfaceId string)`** – wrap `IedServer_setGooseInterfaceId`.
   - Optionally add **`SetGooseInterfaceIdEx(ln, gcbName, interfaceId)`** and **`UseGooseVlanTag(ln, gcbName, useVlanTag)`** if per-GCB configuration is required.

3. **Documentation**
   - Add a **GOOSE** section to STRUCTS.md for `CommParameters` (or note that `GoosePublisherConf` is the Go equivalent).
   - Ensure ENUMS.md lists all `GooseParseError` and GoCB element constants (cross-check with FUNCTIONS.md/STRUCTS.md).

**Deliverables:** SetGoID; server enable/disable + SetGooseInterfaceId; doc updates.  
**Estimate:** Small (on the order of a few hours).

---

### Phase 2 – Threadless and observer ✅ Implemented

4. **GooseReceiver** ✅
   - **`StartThreadless()`** returns **`*GooseReceiverSocket`** (opaque handle); **`StopThreadless()`** – wrap `GooseReceiver_startThreadless` / `GooseReceiver_stopThreadless`.
   - **`HandleMessage(buffer []byte)`** – wrap `GooseReceiver_handleMessage` for custom Ethernet input paths.

5. **GooseSubscriber** ✅
   - **`SetObserver()`** – wrap `GooseSubscriber_setObserver` for “listen to any GOOSE” mode.

6. **Optional: Subscriber with pre-allocated data set** ✅
   - **`Client.ReadDataSetValues(dataSetReference)`** returns **`*ClientDataSet`** (wrap `IedConnection_readDataSetValues`).
   - **`ClientDataSet.GooseDataSetValues()`** returns **`GooseDataSetValues`** (wrap `ClientDataSet_getValues`); keep ClientDataSet alive while the subscriber uses it.
   - **`NewGooseSubscriberWithDataSet(conf, dataSetValues *GooseDataSetValues)`** – wrap `GooseSubscriber_create(goCbRef, dataSetValues)`; pass nil for same behaviour as NewGooseSubscriber.

**Deliverables:** Threadless receiver API; HandleMessage; SetObserver; create-with-dataSet (ReadDataSetValues + NewGooseSubscriberWithDataSet).  
**Estimate:** Small to medium.

---

### Phase 3 – Async and VLAN option ✅ Implemented

7. **Client GoCB async** ✅
   - **`GetGoCBValuesAsync(goCBReference, callback func(*ClientGooseControlBlockValues, error)) (uint32, error)`** – wrap `IedConnection_getGoCBValuesAsync`; returns invoke ID.
   - **`SetGoCBValuesAsync(goCBReference, values, parametersMask, singleRequest, callback func(error)) (uint32, error)`** – wrap `IedConnection_setGoCBValuesAsync`; returns invoke ID.
   - Raw **`IedConnection_GetGoCBValuesHandler`** is not exposed; async is bridged to Go callbacks only.

8. **GoosePublisher createEx** ✅
   - **`NewGoosePublisherEx(conf GoosePublisherConf, useVlanTag bool)`** – wrap `GoosePublisher_createEx`; `NewGoosePublisher(conf)` now calls `NewGoosePublisherEx(conf, true)`.

**Deliverables:** GetGoCBValuesAsync, SetGoCBValuesAsync; NewGoosePublisherEx.  
**Estimate:** Medium (async callback bridging).

---

### Phase 4 – R-GOOSE and advanced (optional) ✅ Partially implemented

9. **R-GOOSE** – **Out of scope**
   - **RSession** is not bound in the Go bindings (no `RSession` type or create API). R-GOOSE is **documented as out-of-scope** for the current bindings. To add it later: bind RSession (e.g. from `r_session.h`), then add **`NewGoosePublisherRemote(session RSession, appId uint16)`** and **`NewGooseReceiverRemote(session RSession)`**.

10. **GooseReceiver createEx** ✅
    - **`NewGooseReceiverEx(buffer []byte)`** – wrap `GooseReceiver_createEx(buffer)`. Pass nil or empty slice for default buffer; non-nil buffer is kept by the receiver for its lifetime.

11. **GoosePublisher publishAndDump** ✅
    - **`PublishAndDump(dataSet *LinkedListValue, msgBuf []byte) (msgLen int, err error)`** – wrap `GoosePublisher_publishAndDump`; returns bytes written into msgBuf for debugging or logging of raw GOOSE payloads.

**Deliverables:** createEx (NewGooseReceiverEx); PublishAndDump. R-GOOSE out-of-scope until RSession is bound.  
**Estimate:** Medium to large (RSession dependency and testing).

---

## Part 4: Summary Checklist

| Item | Phase | C API | Action |
|------|--------|--------|--------|
| SetGoID | 1 | `GoosePublisher_setGoID` | Add to goose_publisher.go |
| EnableGoosePublishing | 1 | `IedServer_enableGoosePublishing` | Add to server.go |
| DisableGoosePublishing | 1 | `IedServer_disableGoosePublishing` | Add to server.go |
| SetGooseInterfaceId | 1 | `IedServer_setGooseInterfaceId` | Add to server.go |
| SetGooseInterfaceIdEx | 1 | `IedServer_setGooseInterfaceIdEx` | Optional; needs LogicalNode |
| UseGooseVlanTag | 1 | `IedServer_useGooseVlanTag` | Optional; needs LogicalNode |
| StartThreadless / StopThreadless | 2 | GooseReceiver_* | ✅ goose_receiver.go |
| HandleMessage | 2 | GooseReceiver_handleMessage | ✅ goose_receiver.go |
| SetObserver | 2 | GooseSubscriber_setObserver | ✅ goose_subscriber.go |
| ReadDataSetValues + NewGooseSubscriberWithDataSet | 2 | IedConnection_readDataSetValues, ClientDataSet_getValues, GooseSubscriber_create | ✅ client_mms.go, goose_subscriber.go |
| GetGoCBValuesAsync | 3 | IedConnection_getGoCBValuesAsync | ✅ client_gocb.go |
| SetGoCBValuesAsync | 3 | IedConnection_setGoCBValuesAsync | ✅ client_gocb.go |
| NewGoosePublisherEx | 3 | GoosePublisher_createEx | ✅ goose_publisher.go |
| NewGoosePublisherRemote | 4 | GoosePublisher_createRemote | Out-of-scope (RSession not bound) |
| NewGooseReceiverRemote | 4 | GooseReceiver_createRemote | Out-of-scope (RSession not bound) |
| NewGooseReceiverEx | 4 | GooseReceiver_createEx | ✅ goose_receiver.go |
| PublishAndDump | 4 | GoosePublisher_publishAndDump | ✅ goose_publisher.go |
| CommParameters type | Doc | struct sCommParameters | STRUCTS.md or type alias |
| ClientGooseControlBlock (opaque) | Doc | — | Note in STRUCTS.md |

---

*End of GOOSE Gaps Analysis*
