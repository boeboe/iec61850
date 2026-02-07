package iec61850

// #include <iec61850_client.h>
import "C"

type MmsType int

type MmsValue struct {
	Type  MmsType
	Value interface{}
}

// UtcTimeValue holds a UTC time from an MMS UTCTime with millisecond precision and time quality.
// Returned when reading a UTCTime attribute via Client.Read() or related APIs.
type UtcTimeValue struct {
	Milliseconds uint64 // Milliseconds since Unix epoch (1970-01-01 00:00:00 UTC)
	TimeQuality  uint8  // IEC 61850 time quality (leapSecondsKnown, clockFailure, clockNotSynchronized, subsecond accuracy)
}

// data types
const (
	Array MmsType = iota
	Structure
	Boolean
	BitString
	Integer
	Unsigned
	Float
	OctetString
	VisibleString
	GeneralizedTime
	BinaryTime
	Bcd
	ObjId
	String
	UTCTime
	DataAccessError
	Int8
	Int16
	Int32
	Int64
	Uint8
	Uint16
	Uint32
)

type MmsDataAccessError int

const (
	DATA_ACCESS_ERROR_SUCCESS_NO_UPDATE             MmsDataAccessError = -3
	DATA_ACCESS_ERROR_NO_RESPONSE                   MmsDataAccessError = -2
	DATA_ACCESS_ERROR_SUCCESS                       MmsDataAccessError = -1
	DATA_ACCESS_ERROR_OBJECT_INVALIDATED            MmsDataAccessError = 0
	DATA_ACCESS_ERROR_HARDWARE_FAULT                MmsDataAccessError = 1
	DATA_ACCESS_ERROR_TEMPORARILY_UNAVAILABLE       MmsDataAccessError = 2
	DATA_ACCESS_ERROR_OBJECT_ACCESS_DENIED          MmsDataAccessError = 3
	DATA_ACCESS_ERROR_OBJECT_UNDEFINED              MmsDataAccessError = 4
	DATA_ACCESS_ERROR_INVALID_ADDRESS               MmsDataAccessError = 5
	DATA_ACCESS_ERROR_TYPE_UNSUPPORTED              MmsDataAccessError = 6
	DATA_ACCESS_ERROR_TYPE_INCONSISTENT             MmsDataAccessError = 7
	DATA_ACCESS_ERROR_OBJECT_ATTRIBUTE_INCONSISTENT MmsDataAccessError = 8
	DATA_ACCESS_ERROR_OBJECT_ACCESS_UNSUPPORTED     MmsDataAccessError = 9
	DATA_ACCESS_ERROR_OBJECT_NONE_EXISTENT          MmsDataAccessError = 10
	DATA_ACCESS_ERROR_OBJECT_VALUE_INVALID          MmsDataAccessError = 11
	DATA_ACCESS_ERROR_UNKNOWN                       MmsDataAccessError = 12
	DATA_ACCESS_ERROR_TYPE_CONFLICT                 MmsDataAccessError = 13
)

type ControlHandlerResult int

const (
	CONTROL_RESULT_FAILED ControlHandlerResult = iota
	CONTROL_RESULT_OK
	CONTROL_RESULT_WAITING
)

type ControlModel int

const (
	// CONTROL_MODEL_STATUS_ONLY No support for control functions. Control object only support status information.
	CONTROL_MODEL_STATUS_ONLY ControlModel = iota
	// CONTROL_MODEL_DIRECT_NORMAL Direct control with normal security: Supports Operate, TimeActivatedOperate (optional), and Cancel (optional).
	CONTROL_MODEL_DIRECT_NORMAL
	// CONTROL_MODEL_SBO_NORMAL Select before operate (SBO) with normal security: Supports Select, Operate, TimeActivatedOperate (optional), and Cancel (optional).
	CONTROL_MODEL_SBO_NORMAL
	// CONTROL_MODEL_DIRECT_ENHANCED Direct control with enhanced security (enhanced security includes the CommandTermination service)
	CONTROL_MODEL_DIRECT_ENHANCED
	// CONTROL_MODEL_SBO_ENHANCED Select before operate (SBO) with enhanced security (enhanced security includes the CommandTermination service)
	CONTROL_MODEL_SBO_ENHANCED
)

type AcseAuthenticationMechanism int

const (
	// ACSE_AUTH_NONE Neither ACSE nor TLS authentication used
	ACSE_AUTH_NONE AcseAuthenticationMechanism = iota

	// ACSE_AUTH_PASSWORD Use ACSE password for client authentication
	ACSE_AUTH_PASSWORD

	// ACSE_AUTH_CERTIFICATE Use ACSE certificate for client authentication
	ACSE_AUTH_CERTIFICATE

	// ACSE_AUTH_TLS Use TLS certificate for client authentication
	ACSE_AUTH_TLS
)

// MmsVariableAccessAttribute describes variable access (read/write) for MMS variable attributes.
type MmsVariableAccessAttribute int32

const (
	MmsVariableReadOnly  MmsVariableAccessAttribute = 0
	MmsVariableWriteOnly MmsVariableAccessAttribute = 1
	MmsVariableReadWrite MmsVariableAccessAttribute = 2
	// Aliases with MMS_ prefix for compatibility.
	MMS_VARIABLE_READ_ONLY  = MmsVariableReadOnly
	MMS_VARIABLE_WRITE_ONLY = MmsVariableWriteOnly
	MMS_VARIABLE_READ_WRITE = MmsVariableReadWrite
	MMS_ACCESS_READ_WRITE   = MmsVariableReadWrite
	MMS_ACCESS_READ_ONLY    = MmsVariableReadOnly
	MMS_ACCESS_WRITE_ONLY   = MmsVariableWriteOnly
)

// MmsFileAccessAttribute describes file access permissions (bitmask).
type MmsFileAccessAttribute int32

const (
	MmsFileAccessNone MmsFileAccessAttribute = 0
	MmsFileRead       MmsFileAccessAttribute = 1
	MmsFileWrite      MmsFileAccessAttribute = 2
	MmsFileDelete     MmsFileAccessAttribute = 4
	// Aliases with MMS_ prefix for compatibility.
	MMS_FILE_ACCESS_NONE = MmsFileAccessNone
	MMS_FILE_READ        = MmsFileRead
	MMS_FILE_WRITE       = MmsFileWrite
	MMS_FILE_DELETE      = MmsFileDelete
)

// MmsJournalVariable identifies a journal variable type (tag or entry ID).
type MmsJournalVariable int32

const (
	MmsJournalVariableTag     MmsJournalVariable = 0
	MmsJournalVariableEntryID MmsJournalVariable = 1
	// Aliases with MMS_ prefix for compatibility.
	MMS_JOURNAL_VARIABLE_TAG      = MmsJournalVariableTag
	MMS_JOURNAL_VARIABLE_ENTRY_ID = MmsJournalVariableEntryID
)

// MmsDeletableType indicates whether a named variable list (or similar) is deletable and by whom.
type MmsDeletableType int32

const (
	MmsDeletableNone           MmsDeletableType = 0
	MmsDeletableAASpecific     MmsDeletableType = 1
	MmsDeletableDomainSpecific MmsDeletableType = 2
	MmsDeletableVMDSpecific    MmsDeletableType = 3
	// Aliases with MMS_ prefix for compatibility.
	MMS_DELETABLE_NONE            = MmsDeletableNone
	MMS_DELETABLE_AA_SPECIFIC     = MmsDeletableAASpecific
	MMS_DELETABLE_DOMAIN_SPECIFIC = MmsDeletableDomainSpecific
	MMS_DELETABLE_VMD_SPECIFIC    = MmsDeletableVMDSpecific
)

// MmsDeletable is an alias for MmsDeletableType for compatibility.
type MmsDeletable = MmsDeletableType

const (
	MMS_DELETABLE_NOT    = MmsDeletableNone
	MMS_DELETABLE_AA     = MmsDeletableAASpecific
	MMS_DELETABLE_DOMAIN = MmsDeletableDomainSpecific
	MMS_DELETABLE_VMD    = MmsDeletableVMDSpecific
)

// MmsNamedVariableListType indicates the scope of a named variable list.
type MmsNamedVariableListType int32

const (
	NamedVariableListTypeVMDSpecific              MmsNamedVariableListType = 0
	NamedVariableListTypeDomainSpecific           MmsNamedVariableListType = 1
	NamedVariableListTypeAssociationSpecific      MmsNamedVariableListType = 2
	NAMED_VARIABLE_LIST_TYPE_VMD_SPECIFIC                                  = NamedVariableListTypeVMDSpecific
	NAMED_VARIABLE_LIST_TYPE_DOMAIN_SPECIFIC                               = NamedVariableListTypeDomainSpecific
	NAMED_VARIABLE_LIST_TYPE_ASSOCIATION_SPECIFIC                          = NamedVariableListTypeAssociationSpecific
)

// MmsServerState represents the MMS server state.
type MmsServerState int32

const (
	MmsServerStateIdle       MmsServerState = 0
	MmsServerStateLoading    MmsServerState = 1
	MmsServerStateRunning    MmsServerState = 2
	MMS_SERVER_STATE_IDLE                   = MmsServerStateIdle
	MMS_SERVER_STATE_LOADING                = MmsServerStateLoading
	MMS_SERVER_STATE_RUNNING                = MmsServerStateRunning
)

// MmsVariableAccessSpec describes a variable reference in a named variable list (domain + item).
type MmsVariableAccessSpec struct {
	DomainID string
	ItemID   string
}

// MmsNamedVariableListAttributes holds attributes of a named variable list (deletable, list type, variable specs).
type MmsNamedVariableListAttributes struct {
	IsDeletable   bool
	DeletableType MmsDeletable
	ListType      MmsNamedVariableListType
	Variables     []MmsVariableAccessSpec
}

// MmsServerConnectionState represents the state of an MMS server connection.
type MmsServerConnectionState int32

const (
	MmsServerConnectionStateIdle        MmsServerConnectionState = 0
	MmsServerConnectionStateAssociation MmsServerConnectionState = 1
	MmsServerConnectionStateConcluded   MmsServerConnectionState = 2
	MMS_CON_STATE_IDLE                                           = MmsServerConnectionStateIdle
	MMS_CON_STATE_ASSOCIATION                                    = MmsServerConnectionStateAssociation
	MMS_CON_STATE_CONCLUDED                                      = MmsServerConnectionStateConcluded
)

// MmsServiceSupportedBitmap is a bitmap of supported MMS services.
type MmsServiceSupportedBitmap uint32

const (
	MmsServiceStatus                               MmsServiceSupportedBitmap = 0x0001
	MmsServiceGetNameList                          MmsServiceSupportedBitmap = 0x0002
	MmsServiceIdentify                             MmsServiceSupportedBitmap = 0x0004
	MmsServiceRead                                 MmsServiceSupportedBitmap = 0x0008
	MmsServiceWrite                                MmsServiceSupportedBitmap = 0x0010
	MmsServiceGetVariableAccess                    MmsServiceSupportedBitmap = 0x0020
	MmsServiceDefineNamedVariableList              MmsServiceSupportedBitmap = 0x0040
	MmsServiceGetNamedVariableListAttrs            MmsServiceSupportedBitmap = 0x0080
	MmsServiceDeleteNamedVariableList              MmsServiceSupportedBitmap = 0x0100
	MmsServiceFileOpen                             MmsServiceSupportedBitmap = 0x0200
	MmsServiceFileRead                             MmsServiceSupportedBitmap = 0x0400
	MmsServiceFileClose                            MmsServiceSupportedBitmap = 0x0800
	MmsServiceFileDelete                           MmsServiceSupportedBitmap = 0x1000
	MmsServiceFileDirectory                        MmsServiceSupportedBitmap = 0x2000
	MmsServiceJournalRead                          MmsServiceSupportedBitmap = 0x4000
	MMS_SERVICE_STATUS                                                       = MmsServiceStatus
	MMS_SERVICE_GET_NAME_LIST                                                = MmsServiceGetNameList
	MMS_SERVICE_IDENTIFY                                                     = MmsServiceIdentify
	MMS_SERVICE_READ                                                         = MmsServiceRead
	MMS_SERVICE_WRITE                                                        = MmsServiceWrite
	MMS_SERVICE_GET_VARIABLE_ACCESS                                          = MmsServiceGetVariableAccess
	MMS_SERVICE_DEFINE_NAMED_VARIABLE_LIST                                   = MmsServiceDefineNamedVariableList
	MMS_SERVICE_GET_NAMED_VARIABLE_LIST_ATTRIBUTES                           = MmsServiceGetNamedVariableListAttrs
	MMS_SERVICE_DELETE_NAMED_VARIABLE_LIST                                   = MmsServiceDeleteNamedVariableList
	MMS_SERVICE_FILE_OPEN                                                    = MmsServiceFileOpen
	MMS_SERVICE_FILE_READ                                                    = MmsServiceFileRead
	MMS_SERVICE_FILE_CLOSE                                                   = MmsServiceFileClose
	MMS_SERVICE_FILE_DELETE                                                  = MmsServiceFileDelete
	MMS_SERVICE_FILE_DIRECTORY                                               = MmsServiceFileDirectory
	MMS_SERVICE_JOURNAL_READ                                                 = MmsServiceJournalRead
)

// IsoConnectionParameters holds ISO layer connection parameters (AP title, selectors).
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

// MmsJournalEntry holds a single journal entry (entry ID bytes, occurrence time, content).
type MmsJournalEntry struct {
	EntryID      []byte
	OccurTime    uint64
	EntryContent *MmsValue
}

// MmsJournalVariableSpec specifies a journal variable (tag and optional value spec).
type MmsJournalVariableSpec struct {
	Tag       string
	ValueSpec *MmsVariableSpecificationRef
}

// MmsConnectionParameters holds MMS layer connection parameters (max outstanding calls, PDU size, etc.).
// Returned by Client.GetConnectionParameters after connection is established.
type MmsConnectionParameters struct {
	MaxServOutstandingCalling int32
	MaxServOutstandingCalled  int32
	DataStructureNestingLevel int32
	MaxPduSize                int32
	ServicesSupported         [11]uint8
}

// MmsServerStatus holds the result of GetServerStatus (VMD logical/physical status).
type MmsServerStatus struct {
	VmdLogicalStatus  int32
	VmdPhysicalStatus int32
	LocalDetail       int32
}

// MmsVariableListAttributes holds attributes of a named variable list (deletable, variable names).
type MmsVariableListAttributes struct {
	IsDeletable   bool
	DeletableType MmsDeletableType
	NumberOfItems int32
	VariableNames []string
}

// MmsFileDirectoryEntryEx holds extended file directory entry (name, size, last modified, attributes).
type MmsFileDirectoryEntryEx struct {
	Filename         string
	FileSize         uint32
	LastModifiedTime uint64
	FileAttributes   uint32
}
