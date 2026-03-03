package iec61850

// #include <iec61850_server.h>
import "C"

import (
	"os"
	"unsafe"
)

type IedModel struct {
	Model *C.IedModel
}

// This is a little hacky but it works for calls from runtime_scl.
//
// The pointer must be a pointer to the C version of the IedModel.
func NewIedModelFromPointer(model unsafe.Pointer) *IedModel {
	return &IedModel{
		Model: (*C.IedModel)(model),
	}
}

type ModelNode struct {
	ObjectReference string
	_modelNode      unsafe.Pointer
}

func NewIedModel(name string) *IedModel {
	cname, freeCname := allocCString(name)
	defer freeCname()
	return &IedModel{
		Model: C.IedModel_create(cname),
	}
}

func (m *IedModel) Destroy() {
	C.IedModel_destroy(m.Model)
}

func (m *IedModel) GetModelNodeByObjectReference(objectRef string) *ModelNode {
	cObjectRef, freeCObjectRef := allocCString(objectRef)
	defer freeCObjectRef()

	do := C.IedModel_getModelNodeByObjectReference(m.Model, cObjectRef)
	if do == nil {
		return nil
	}
	return &ModelNode{_modelNode: unsafe.Pointer(do), ObjectReference: objectRef}
}

func (m *ModelNode) GetLogicalNode(node string) *LogicalNode {
	cNode, freeCNode := allocCString(node)
	defer freeCNode()

	logicalNode := C.LogicalDevice_getLogicalNode((*C.LogicalDevice)(m._modelNode), cNode)

	return &LogicalNode{
		node: logicalNode,
	}
}

func (m *ModelNode) ConvertToDataObject() *DataObject {
	return &DataObject{
		object: (*C.DataObject)(m._modelNode),
	}
}

func CreateModelFromConfigFileEx(filepath string) (*IedModel, error) {
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	cFilepath, freeCFilepath := allocCString(filepath)
	// Free memory
	defer freeCFilepath()
	model := &IedModel{
		Model: C.ConfigFileParser_createModelFromConfigFileEx(cFilepath),
	}
	return model, nil
}

type LogicalDevice struct {
	device *C.LogicalDevice
}

func (m *IedModel) CreateLogicalDevice(name string) *LogicalDevice {
	cname, freeCname := allocCString(name)
	defer freeCname()
	return &LogicalDevice{
		device: C.LogicalDevice_create(cname, m.Model),
	}
}

type LogicalNode struct {
	node *C.LogicalNode
}

func (d *LogicalDevice) CreateLogicalNode(name string) *LogicalNode {
	cname, freeCname := allocCString(name)
	defer freeCname()
	return &LogicalNode{
		node: C.LogicalNode_create(cname, d.device),
	}
}

type DataObject struct {
	object *C.DataObject
}

// ENS: EnumerationString
// VSS: Visible String Setting
// SAV: Sampled Value
// APC: Analogue Process Control

func (n *LogicalNode) CreateDataObjectCDC_ENS(name string) *DataObject {
	cname, freeCname := allocCString(name)
	defer freeCname()
	return &DataObject{
		object: C.CDC_ENS_create(cname, (*C.ModelNode)(n.node), 0),
	}
}

func (n *LogicalNode) CreateDataObjectCDC_VSS(name string) *DataObject {
	cname, freeCname := allocCString(name)
	defer freeCname()
	return &DataObject{
		object: C.CDC_VSS_create(cname, (*C.ModelNode)(n.node), 0),
	}
}

func (n *LogicalNode) CreateDataObjectCDC_SAV(name string, isInteger bool) *DataObject {
	cname, freeCname := allocCString(name)
	defer freeCname()
	return &DataObject{
		object: C.CDC_SAV_create(cname, (*C.ModelNode)(n.node), 0, C.bool(isInteger)),
	}
}

func (n *LogicalNode) CreateDataObjectCDC_APC(name string, ctlModel int) *DataObject {
	cname, freeCname := allocCString(name)
	defer freeCname()
	return &DataObject{
		object: C.CDC_APC_create(cname, (*C.ModelNode)(n.node), 0, C.uint(ctlModel), C.bool(false)),
	}
}

type DataAttribute struct {
	attribute *C.DataAttribute
}

func (do *DataObject) GetChild(name string) *DataAttribute {
	cname, freeCname := allocCString(name)
	defer freeCname()
	return &DataAttribute{
		attribute: (*C.DataAttribute)(unsafe.Pointer(C.ModelNode_getChild((*C.ModelNode)(unsafe.Pointer(do.object)), cname))),
	}
}

type DataSet struct {
	dataSet *C.DataSet
}

// CreateDataSet creates a new DataSet under this LogicalNode.
func (ln *LogicalNode) CreateDataSet(name string) *DataSet {
	cName, freeCName := allocCString(name)
	defer freeCName()

	cDataSet := C.DataSet_create(cName, ln.node)
	return &DataSet{dataSet: cDataSet}
}

// AddDataSetEntry adds a new DataSetEntry to this DataSet.
func (ds *DataSet) AddDataSetEntry(ref string) {
	cRef, freeCRef := allocCString(ref)
	defer freeCRef()

	C.DataSetEntry_create(ds.dataSet, cRef, -1, nil)
}
