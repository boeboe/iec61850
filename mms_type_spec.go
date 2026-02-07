package iec61850

/*
#include <iec61850_client.h>
#include <mms_value.h>
#include <mms_type_spec.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// MmsVariableSpecification is the C MmsVariableSpecification struct type (for type mapping).
// Use MmsVariableSpecificationRef for API calls that return or take variable specifications.
type MmsVariableSpecification C.MmsVariableSpecification

// MmsTypeSpecification: in libiec61850 type information is represented by MmsVariableSpecification.
// This alias exists for API compatibility with the MMS type specification concept.
type MmsTypeSpecification = MmsVariableSpecification

// MmsNamedVariableList is the C named variable list pointer type (opaque; used on server side).
type MmsNamedVariableList C.MmsNamedVariableList

// MmsVariableSpecificationRef wraps a C MmsVariableSpecification pointer, typically obtained from
// Client.GetVariableAccessAttributes or from the server model. Caller must call Free() when done.
type MmsVariableSpecificationRef struct {
	c *C.MmsVariableSpecification
}

// Free releases the C memory. It is safe to call multiple times.
func (r *MmsVariableSpecificationRef) Free() {
	if r != nil && r.c != nil {
		C.MmsVariableSpecification_destroy(r.c)
		r.c = nil
	}
}

// GetType returns the MMS type of the variable.
func (r *MmsVariableSpecificationRef) GetType() MmsType {
	if r == nil || r.c == nil {
		return -1
	}
	return MmsType(C.MmsVariableSpecification_getType(r.c))
}

// GetName returns the variable name. The returned string is only valid while the specification exists.
func (r *MmsVariableSpecificationRef) GetName() string {
	if r == nil || r.c == nil {
		return ""
	}
	n := C.MmsVariableSpecification_getName(r.c)
	if n == nil {
		return ""
	}
	return C.GoString(n)
}

// GetSize returns the number of elements for structures/arrays, or bit/byte size for other types. Returns -1 if not applicable.
func (r *MmsVariableSpecificationRef) GetSize() int {
	if r == nil || r.c == nil {
		return -1
	}
	return int(C.MmsVariableSpecification_getSize(r.c))
}

// GetChildSpecificationByIndex returns the child variable specification at the given index (for structure or array). Caller does not own the returned ref.
func (r *MmsVariableSpecificationRef) GetChildSpecificationByIndex(index int) *MmsVariableSpecificationRef {
	if r == nil || r.c == nil {
		return nil
	}
	child := C.MmsVariableSpecification_getChildSpecificationByIndex(r.c, C.int(index))
	if child == nil {
		return nil
	}
	return &MmsVariableSpecificationRef{c: child}
}

// GetChildSpecificationByName returns the child variable specification with the given name. Caller does not own the returned ref.
func (r *MmsVariableSpecificationRef) GetChildSpecificationByName(name string) *MmsVariableSpecificationRef {
	if r == nil || r.c == nil {
		return nil
	}
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	child := C.MmsVariableSpecification_getChildSpecificationByName(r.c, cs, nil)
	if child == nil {
		return nil
	}
	return &MmsVariableSpecificationRef{c: child}
}

// GetArrayElementSpecification returns the element type specification for an array. Caller does not own the returned ref.
func (r *MmsVariableSpecificationRef) GetArrayElementSpecification() *MmsVariableSpecificationRef {
	if r == nil || r.c == nil {
		return nil
	}
	el := C.MmsVariableSpecification_getArrayElementSpecification(r.c)
	if el == nil {
		return nil
	}
	return &MmsVariableSpecificationRef{c: el}
}

// IsValueOfType checks whether the given value has exactly the same type as this variable specification.
func (r *MmsVariableSpecificationRef) IsValueOfType(v *MmsValueRef) bool {
	if r == nil || r.c == nil || v == nil || v.c == nil {
		return false
	}
	return bool(C.MmsVariableSpecification_isValueOfType(r.c, v.c))
}

// GetStructureElements returns a list of structure element names for MMS_STRUCTURE types. Caller must not free the returned strings (they are valid while the spec exists).
func (r *MmsVariableSpecificationRef) GetStructureElements() []string {
	if r == nil || r.c == nil {
		return nil
	}
	list := C.MmsVariableSpecification_getStructureElements(r.c)
	if list == nil {
		return nil
	}
	defer C.LinkedList_destroy(list)
	var out []string
	for node := list; node != nil; node = C.LinkedList_getNext(node) {
		data := C.LinkedList_getData(node)
		if data != nil {
			out = append(out, C.GoString((*C.char)(data)))
		}
	}
	return out
}
