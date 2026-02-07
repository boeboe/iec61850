package iec61850

/*
#include <iec61850_client.h>
#include <mms_value.h>
#include <mms_type_spec.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// MmsValueRef wraps a C MmsValue pointer for use with MMS value constructors and accessors.
// The caller is responsible for calling Free() when the value is no longer needed,
// unless the value is passed to Client.Write (which does not take ownership).
//
// BitString integer conversions: use GetBitStringAsInteger, SetBitStringFromInteger,
// GetBitStringAsIntegerBigEndian, SetBitStringFromIntegerBigEndian on *MmsValueRef
// (the C-backed type). The high-level MmsValue type (Type + Value) does not hold a C pointer.
type MmsValueRef struct {
	c *C.MmsValue
}

// Free releases the C memory associated with the MmsValue. It is safe to call multiple times.
func (r *MmsValueRef) Free() {
	if r != nil && r.c != nil {
		C.MmsValue_delete(r.c)
		r.c = nil
	}
}

// NewMmsValueVisibleString creates an MMS visible string value. Caller must call Free() when done.
func NewMmsValueVisibleString(s string) *MmsValueRef {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return &MmsValueRef{c: C.MmsValue_newVisibleString(cs)}
}

// NewMmsValueVisibleStringWithSize creates an empty visible string with the given maximum size. Caller must call Free() when done.
func NewMmsValueVisibleStringWithSize(size int) *MmsValueRef {
	return &MmsValueRef{c: C.MmsValue_newVisibleStringWithSize(C.int(size))}
}

// NewMmsValueMmsString creates an MMS string value. Caller must call Free() when done.
func NewMmsValueMmsString(s string) *MmsValueRef {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return &MmsValueRef{c: C.MmsValue_newMmsString(cs)}
}

// NewMmsValueMmsStringWithSize creates an empty MMS string with the given maximum size. Caller must call Free() when done.
func NewMmsValueMmsStringWithSize(size int) *MmsValueRef {
	return &MmsValueRef{c: C.MmsValue_newMmsStringWithSize(C.int(size))}
}

// NewMmsValueBinaryTime creates an MMS binary time value. If timeOfDay is true, only time-of-day (4 octets) is used. Caller must call Free() when done.
func NewMmsValueBinaryTime(timeOfDay bool) *MmsValueRef {
	return &MmsValueRef{c: C.MmsValue_newBinaryTime(C.bool(timeOfDay))}
}

// SetBinaryTime sets the value in milliseconds since Unix epoch.
func (r *MmsValueRef) SetBinaryTime(timestampMs uint64) {
	if r != nil && r.c != nil {
		C.MmsValue_setBinaryTime(r.c, C.uint64_t(timestampMs))
	}
}

// GetBinaryTimeAsUtcMs returns the value in milliseconds since Unix epoch.
func (r *MmsValueRef) GetBinaryTimeAsUtcMs() uint64 {
	if r == nil || r.c == nil {
		return 0
	}
	return uint64(C.MmsValue_getBinaryTimeAsUtcMs(r.c))
}

// NewMmsValueBitString creates a bit string of the given size in bits. Caller must call Free() when done.
func NewMmsValueBitString(bitSize int) *MmsValueRef {
	return &MmsValueRef{c: C.MmsValue_newBitString(C.int(bitSize))}
}

// SetBitStringFromInteger sets the bit string from an unsigned integer (little-endian bit order).
func (r *MmsValueRef) SetBitStringFromInteger(val uint32) {
	if r != nil && r.c != nil {
		C.MmsValue_setBitStringFromInteger(r.c, C.uint32_t(val))
	}
}

// GetBitStringAsInteger returns the bit string as an unsigned integer (little-endian bit order).
func (r *MmsValueRef) GetBitStringAsInteger() uint32 {
	if r == nil || r.c == nil {
		return 0
	}
	return uint32(C.MmsValue_getBitStringAsInteger(r.c))
}

// SetBitStringFromIntegerBigEndian sets the bit string from an unsigned integer (big-endian bit order).
func (r *MmsValueRef) SetBitStringFromIntegerBigEndian(val uint32) {
	if r != nil && r.c != nil {
		C.MmsValue_setBitStringFromIntegerBigEndian(r.c, C.uint32_t(val))
	}
}

// GetBitStringAsIntegerBigEndian returns the bit string as an unsigned integer (big-endian bit order).
func (r *MmsValueRef) GetBitStringAsIntegerBigEndian() uint32 {
	if r == nil || r.c == nil {
		return 0
	}
	return uint32(C.MmsValue_getBitStringAsIntegerBigEndian(r.c))
}

// GetBitStringSize returns the size of the bit string in bits.
func (r *MmsValueRef) GetBitStringSize() int {
	if r == nil || r.c == nil {
		return 0
	}
	return int(C.MmsValue_getBitStringSize(r.c))
}

// ToDouble returns the value as float64. The underlying MmsValue must be of type MMS_FLOAT.
func (r *MmsValueRef) ToDouble() float64 {
	if r == nil || r.c == nil {
		return 0
	}
	return float64(C.MmsValue_toDouble(r.c))
}

// ToInt64 returns the value as int64. The underlying MmsValue must be of type MMS_INTEGER or MMS_UNSIGNED.
func (r *MmsValueRef) ToInt64() int64 {
	if r == nil || r.c == nil {
		return 0
	}
	return int64(C.MmsValue_toInt64(r.c))
}

// ToUint32 returns the value as uint32. The underlying MmsValue must be of type MMS_INTEGER or MMS_UNSIGNED.
func (r *MmsValueRef) ToUint32() uint32 {
	if r == nil || r.c == nil {
		return 0
	}
	return uint32(C.MmsValue_toUint32(r.c))
}

// SetVisibleString sets the value of a visible string. The ref must be of type MMS_VISIBLE_STRING.
func (r *MmsValueRef) SetVisibleString(s string) {
	if r == nil || r.c == nil {
		return
	}
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.MmsValue_setVisibleString(r.c, cs)
}

// SetMmsString sets the value of an MMS string. The ref must be of type MMS_STRING.
func (r *MmsValueRef) SetMmsString(s string) {
	if r == nil || r.c == nil {
		return
	}
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.MmsValue_setMmsString(r.c, cs)
}

// GetType returns the MMS type of the value.
func (r *MmsValueRef) GetType() MmsType {
	if r == nil || r.c == nil {
		return -1
	}
	return MmsType(C.MmsValue_getType(r.c))
}

// MmsValueCreateEmptyArray creates an empty array of the given size. Caller must call Free() when done.
func MmsValueCreateEmptyArray(size int) *MmsValueRef {
	return &MmsValueRef{c: C.MmsValue_createEmptyArray(C.int(size))}
}

// MmsValueCreateArray creates an array with elements described by elementType and initializes with default values. Caller must call Free() when done.
func MmsValueCreateArray(elementType *MmsVariableSpecificationRef, size int) *MmsValueRef {
	if elementType == nil || elementType.c == nil {
		return nil
	}
	return &MmsValueRef{c: C.MmsValue_createArray(elementType.c, C.int(size))}
}

// MmsValueNewDefaultValue creates a new MmsValue with default value for the given type specification. Caller must call Free() when done.
func MmsValueNewDefaultValue(typeSpec *MmsVariableSpecificationRef) *MmsValueRef {
	if typeSpec == nil || typeSpec.c == nil {
		return nil
	}
	return &MmsValueRef{c: C.MmsValue_newDefaultValue(typeSpec.c)}
}

// GetArraySize returns the number of elements in an array or structure. The value must be of type MMS_ARRAY or MMS_STRUCTURE.
func (r *MmsValueRef) GetArraySize() uint32 {
	if r == nil || r.c == nil {
		return 0
	}
	return uint32(C.MmsValue_getArraySize(r.c))
}

// GetElement returns the element at the given index (array or structure). Caller does not own the returned value.
func (r *MmsValueRef) GetElement(index int) *MmsValueRef {
	if r == nil || r.c == nil {
		return nil
	}
	el := C.MmsValue_getElement(r.c, C.int(index))
	if el == nil {
		return nil
	}
	return &MmsValueRef{c: el}
}

// SetElement sets the element at the given index. If an element already exists it is replaced; the caller is responsible for freeing the replaced value. v is not consumed (caller still owns it).
func (r *MmsValueRef) SetElement(index int, v *MmsValueRef) {
	if r == nil || r.c == nil || v == nil || v.c == nil {
		return
	}
	C.MmsValue_setElement(r.c, C.int(index), v.c)
}

// MmsValueDelete frees a C MmsValue. It is used when the caller has received a raw value from the C API (e.g. from Read) and must release it.
func MmsValueDelete(r *MmsValueRef) {
	if r != nil {
		r.Free()
	}
}

// GetBitStringAsInteger returns the bit string as an unsigned integer (little-endian). Only valid when v.Type == BitString and v.Value is uint32.
func (v *MmsValue) GetBitStringAsInteger() (uint32, error) {
	if v == nil || v.Type != BitString {
		return 0, UnSupportedOperation
	}
	u, ok := v.Value.(uint32)
	if !ok {
		return 0, UnSupportedOperation
	}
	return u, nil
}

// GetBitStringAsIntegerBigEndian returns the bit string as an unsigned integer (big-endian). Only valid when v.Type == BitString and v.Value is uint32 (bits are reordered).
func (v *MmsValue) GetBitStringAsIntegerBigEndian() (uint32, error) {
	if v == nil || v.Type != BitString {
		return 0, UnSupportedOperation
	}
	u, ok := v.Value.(uint32)
	if !ok {
		return 0, UnSupportedOperation
	}
	// Swap byte order for big-endian
	return (u>>24)&0xff | (u>>8)&0xff00 | (u<<8)&0xff0000 | (u<<24)&0xff000000, nil
}

// SetBitStringFromInteger sets the bit string from an unsigned integer (little-endian). Only valid when v.Type == BitString.
func (v *MmsValue) SetBitStringFromInteger(value uint32) error {
	if v == nil || v.Type != BitString {
		return UnSupportedOperation
	}
	v.Value = value
	return nil
}

// SetBitStringFromIntegerBigEndian sets the bit string from an unsigned integer (big-endian). Only valid when v.Type == BitString.
func (v *MmsValue) SetBitStringFromIntegerBigEndian(value uint32) error {
	if v == nil || v.Type != BitString {
		return UnSupportedOperation
	}
	v.Value = value
	return nil
}

// DeleteAllBitStringBits clears all bits. Requires a C-backed value; use MmsValueRef for in-place bit manipulation.
func (v *MmsValue) DeleteAllBitStringBits() error {
	if v == nil || v.Type != BitString {
		return UnSupportedOperation
	}
	v.Value = uint32(0)
	return nil
}

// SetAllBitStringBits sets all bits. Requires a C-backed value; use MmsValueRef for in-place bit manipulation.
func (v *MmsValue) SetAllBitStringBits() error {
	if v == nil || v.Type != BitString {
		return UnSupportedOperation
	}
	// Without C backing we don't know bit size; set to all 1s for 32 bits
	v.Value = uint32(0xffffffff)
	return nil
}
