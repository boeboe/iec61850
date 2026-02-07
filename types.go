package iec61850

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
	MmsVariableReadOnly   MmsVariableAccessAttribute = 0
	MmsVariableWriteOnly  MmsVariableAccessAttribute = 1
	MmsVariableReadWrite  MmsVariableAccessAttribute = 2
	// Aliases with MMS_ prefix for compatibility.
	MMS_VARIABLE_READ_ONLY  = MmsVariableReadOnly
	MMS_VARIABLE_WRITE_ONLY = MmsVariableWriteOnly
	MMS_VARIABLE_READ_WRITE = MmsVariableReadWrite
)

// MmsFileAccessAttribute describes file access permissions (bitmask).
type MmsFileAccessAttribute int32

const (
	MmsFileRead   MmsFileAccessAttribute = 1
	MmsFileWrite  MmsFileAccessAttribute = 2
	MmsFileDelete MmsFileAccessAttribute = 4
	// Aliases with MMS_ prefix for compatibility.
	MMS_FILE_READ   = MmsFileRead
	MMS_FILE_WRITE  = MmsFileWrite
	MMS_FILE_DELETE = MmsFileDelete
)

// MmsJournalVariable identifies a journal variable type (tag or entry ID).
type MmsJournalVariable int32

const (
	MmsJournalVariableTag     MmsJournalVariable = 0
	MmsJournalVariableEntryID MmsJournalVariable = 1
	// Aliases with MMS_ prefix for compatibility.
	MMS_JOURNAL_VARIABLE_TAG     = MmsJournalVariableTag
	MMS_JOURNAL_VARIABLE_ENTRY_ID = MmsJournalVariableEntryID
)

// MmsDeletableType indicates whether a named variable list (or similar) is deletable and by whom.
type MmsDeletableType int32

const (
	MmsDeletableNone            MmsDeletableType = 0
	MmsDeletableAASpecific      MmsDeletableType = 1
	MmsDeletableDomainSpecific MmsDeletableType = 2
	// Aliases with MMS_ prefix for compatibility.
	MMS_DELETABLE_NONE             = MmsDeletableNone
	MMS_DELETABLE_AA_SPECIFIC      = MmsDeletableAASpecific
	MMS_DELETABLE_DOMAIN_SPECIFIC  = MmsDeletableDomainSpecific
)

// MmsJournalEntry holds a single journal entry (simplified form with numeric IDs and a single value).
// For full journal entries with multiple variables use client_mms.JournalEntry from ReadJournalTimeRange/ReadJournalStartAfter.
type MmsJournalEntry struct {
	EntryID     uint64
	OccurTime   uint64
	EntryValues *MmsValue
}
