package iec61850

/*
#include <iec61850_client.h>
#include <mms_value.h>
#include <mms_type_spec.h>
#include <stdlib.h>
#include <string.h>

// Allocate and fill MmsVariableSpecification for simple types. Size: bit width for integer/unsigned/bitstring, max bytes for string types, exponent width for float, 4 or 6 for binary time.
static MmsVariableSpecification* goMmsVarSpecCreateSimple(MmsType type, const char* name, int size) {
	MmsVariableSpecification* s = (MmsVariableSpecification*)malloc(sizeof(MmsVariableSpecification));
	if (!s) return NULL;
	memset(s, 0, sizeof(MmsVariableSpecification));
	s->type = type;
	s->name = name ? strdup(name) : NULL;
	switch (type) {
		case MMS_BOOLEAN: break;
		case MMS_INTEGER: s->typeSpec.integer = size; break;
		case MMS_UNSIGNED: s->typeSpec.unsignedInteger = size; break;
		case MMS_BIT_STRING: s->typeSpec.bitString = size; break;
		case MMS_OCTET_STRING: s->typeSpec.octetString = size; break;
		case MMS_VISIBLE_STRING: s->typeSpec.visibleString = size; break;
		case MMS_STRING: s->typeSpec.mmsString = size; break;
		case MMS_FLOAT: s->typeSpec.floatingpoint.formatWidth = 64; s->typeSpec.floatingpoint.exponentWidth = (uint8_t)(size & 0xff); break;
		case MMS_BINARY_TIME: s->typeSpec.binaryTime = size; break;
		default: break;
	}
	return s;
}

// Create structure; takes ownership of elements array and each element. Name may be NULL.
static MmsVariableSpecification* goMmsVarSpecCreateStructure(const char* name, MmsVariableSpecification** elements, int count) {
	MmsVariableSpecification* s = (MmsVariableSpecification*)malloc(sizeof(MmsVariableSpecification));
	if (!s) return NULL;
	memset(s, 0, sizeof(MmsVariableSpecification));
	s->type = MMS_STRUCTURE;
	s->name = name ? strdup(name) : NULL;
	s->typeSpec.structure.elementCount = count;
	s->typeSpec.structure.elements = elements;
	return s;
}

// Create array; takes ownership of elementType. Name may be NULL.
static MmsVariableSpecification* goMmsVarSpecCreateArray(const char* name, MmsVariableSpecification* elementType, int elementCount) {
	MmsVariableSpecification* s = (MmsVariableSpecification*)malloc(sizeof(MmsVariableSpecification));
	if (!s) return NULL;
	memset(s, 0, sizeof(MmsVariableSpecification));
	s->type = MMS_ARRAY;
	s->name = name ? strdup(name) : NULL;
	s->typeSpec.array.elementCount = elementCount;
	s->typeSpec.array.elementTypeSpec = elementType;
	return s;
}

// Recursively free a spec tree created by the goMmsVarSpecCreate* functions. Do not use for specs from the library.
static void goMmsVarSpecFreeOwned(MmsVariableSpecification* s) {
	if (!s) return;
	if (s->name) free((void*)s->name);
	if (s->type == MMS_STRUCTURE) {
		int i;
		for (i = 0; i < s->typeSpec.structure.elementCount; i++)
			goMmsVarSpecFreeOwned(s->typeSpec.structure.elements[i]);
		free(s->typeSpec.structure.elements);
	} else if (s->type == MMS_ARRAY) {
		goMmsVarSpecFreeOwned(s->typeSpec.array.elementTypeSpec);
	}
	free(s);
}
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
// Client.GetVariableAccessAttributes, MmsConnection.GetVariableAccessAttributes(Async), or from
// NewMmsVariableSpecification / CreateStructure / CreateArray. Caller must call Free() on the root ref when done.
// Do not call Free() on refs returned by GetChildSpecificationByIndex, GetChildSpecificationByName, or GetArrayElementSpecification.
type MmsVariableSpecificationRef struct {
	c            *C.MmsVariableSpecification
	owned        bool // true if this ref is a root that should be freed by the user
	libraryOwned bool // if owned, true = use library destroy; false = use our free for Go-created specs
}

// Free releases the C memory for a root ref. It is a no-op for child refs and safe to call multiple times.
func (r *MmsVariableSpecificationRef) Free() {
	if r == nil || r.c == nil {
		return
	}
	if !r.owned {
		r.c = nil
		return
	}
	if r.libraryOwned {
		C.MmsVariableSpecification_destroy(r.c)
	} else {
		C.goMmsVarSpecFreeOwned(r.c)
	}
	r.c = nil
	r.owned = false
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

// GetChildSpecificationByIndex returns the child variable specification at the given index (for structure or array). Do not call Free() on the returned ref.
func (r *MmsVariableSpecificationRef) GetChildSpecificationByIndex(index int) *MmsVariableSpecificationRef {
	if r == nil || r.c == nil {
		return nil
	}
	child := C.MmsVariableSpecification_getChildSpecificationByIndex(r.c, C.int(index))
	if child == nil {
		return nil
	}
	return &MmsVariableSpecificationRef{c: child, owned: false}
}

// GetChildSpecificationByName returns the child variable specification with the given name. Do not call Free() on the returned ref.
func (r *MmsVariableSpecificationRef) GetChildSpecificationByName(name string) *MmsVariableSpecificationRef {
	if r == nil || r.c == nil {
		return nil
	}
	cs, freecs := allocCString(name)
	defer freecs()
	child := C.MmsVariableSpecification_getChildSpecificationByName(r.c, cs, nil)
	if child == nil {
		return nil
	}
	return &MmsVariableSpecificationRef{c: child, owned: false}
}

// GetArrayElementSpecification returns the element type specification for an array. Do not call Free() on the returned ref.
func (r *MmsVariableSpecificationRef) GetArrayElementSpecification() *MmsVariableSpecificationRef {
	if r == nil || r.c == nil {
		return nil
	}
	el := C.MmsVariableSpecification_getArrayElementSpecification(r.c)
	if el == nil {
		return nil
	}
	return &MmsVariableSpecificationRef{c: el, owned: false}
}

// IsValueOfType checks whether the given value has exactly the same type as this variable specification.
func (r *MmsVariableSpecificationRef) IsValueOfType(v *MmsValueRef) bool {
	if r == nil || r.c == nil || v == nil || v.c == nil {
		return false
	}
	return bool(C.MmsVariableSpecification_isValueOfType(r.c, v.c))
}

// GetChildValue returns the child of value corresponding to the relative MMS path childId (use "$" as separator). The ref must be for a structure; value must be the corresponding MmsValue. Caller does not own the returned value.
func (r *MmsVariableSpecificationRef) GetChildValue(value *MmsValueRef, childId string) *MmsValueRef {
	if r == nil || r.c == nil || value == nil || value.c == nil {
		return nil
	}
	cChild, freecChild := allocCString(childId)
	defer freecChild()
	el := C.MmsVariableSpecification_getChildValue(r.c, value.c, cChild)
	if el == nil {
		return nil
	}
	return &MmsValueRef{c: el}
}

// GetNamedVariableRecursive returns the variable specification of the child specified by the relative MMS name nameId (use "$" as separator). Do not call Free() on the returned ref.
func (r *MmsVariableSpecificationRef) GetNamedVariableRecursive(nameId string) *MmsVariableSpecificationRef {
	if r == nil || r.c == nil {
		return nil
	}
	cName, freecName := allocCString(nameId)
	defer freecName()
	child := C.MmsVariableSpecification_getNamedVariableRecursive(r.c, cName)
	if child == nil {
		return nil
	}
	return &MmsVariableSpecificationRef{c: child, owned: false}
}

// GetExponentWidth returns the exponent width for floating-point types; returns a meaningful value only for MMS_FLOAT/MMS_VISIBLE_STRING etc. as defined by the library.
func (r *MmsVariableSpecificationRef) GetExponentWidth() int {
	if r == nil || r.c == nil {
		return 0
	}
	return int(C.MmsVariableSpecification_getExponentWidth(r.c))
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

// NewMmsVariableSpecification creates a simple type specification. Caller must call Free() when done.
// Size meaning: bit width for MMS_INTEGER, MMS_UNSIGNED, MMS_BIT_STRING; max bytes for MMS_OCTET_STRING,
// MMS_VISIBLE_STRING, MMS_STRING; exponent width for MMS_FLOAT (e.g. 11 for double); 4 or 6 for MMS_BINARY_TIME.
// Name may be empty. For MMS_BOOLEAN and other types with no size, pass 0.
func NewMmsVariableSpecification(typ MmsType, name string, size int) *MmsVariableSpecificationRef {
	var cName *C.char
	var freeCName func()
	if name != "" {
		cName, freeCName = allocCString(name)
		defer freeCName()
	}
	c := C.goMmsVarSpecCreateSimple(C.MmsType(typ), cName, C.int(size))
	if c == nil {
		return nil
	}
	return &MmsVariableSpecificationRef{c: c, owned: true, libraryOwned: false}
}

// CreateStructure creates a structure type specification containing the given element specs.
// The element refs are incorporated into the structure; do not call Free() on them after.
// Caller must call Free() on the returned ref when done.
func CreateStructure(name string, elements []*MmsVariableSpecificationRef) *MmsVariableSpecificationRef {
	if len(elements) == 0 {
		return nil
	}
	n := C.size_t(len(elements)) * C.size_t(unsafe.Sizeof(uintptr(0)))
	cArr, freeCArr := allocCMalloc(n)
	defer freeCArr()
	base := (*[1 << 20]*C.MmsVariableSpecification)(cArr)
	for i, el := range elements {
		if el != nil && el.c != nil {
			base[i] = el.c
		}
	}
	var cName *C.char
	var freeCName func()
	if name != "" {
		cName, freeCName = allocCString(name)
		defer freeCName()
	}
	c := C.goMmsVarSpecCreateStructure(cName, (**C.MmsVariableSpecification)(cArr), C.int(len(elements)))
	if c == nil {
		return nil
	}
	for _, el := range elements {
		if el != nil {
			el.owned = false
		}
	}
	return &MmsVariableSpecificationRef{c: c, owned: true, libraryOwned: false}
}

// CreateArray creates an array type specification with the given element type and length.
// The elementType ref is incorporated; do not call Free() on it after.
// Caller must call Free() on the returned ref when done.
func CreateArray(name string, elementType *MmsVariableSpecificationRef, elementCount int) *MmsVariableSpecificationRef {
	if elementType == nil || elementType.c == nil || elementCount < 0 {
		return nil
	}
	var cName *C.char
	var freeCName func()
	if name != "" {
		cName, freeCName = allocCString(name)
		defer freeCName()
	}
	c := C.goMmsVarSpecCreateArray(cName, elementType.c, C.int(elementCount))
	if c == nil {
		return nil
	}
	elementType.owned = false
	return &MmsVariableSpecificationRef{c: c, owned: true, libraryOwned: false}
}
