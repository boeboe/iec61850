package iec61850

// #include <iec61850_client.h>
import "C"
import (
	"fmt"

	"github.com/spf13/cast"
)

const (
	ActDA  = "%s/%s.SGCB.ActSG"
	EditDA = "%s/%s.SGCB.EditSG"
	CnfDA  = "%s/%s.SGCB.CnfEdit"
)

type SettingGroup struct {
	NumOfSG int
	ActSG   int
	EditSG  int
	CnfEdit bool
}

// WriteSG writes the SettingGroup
func (c *Client) WriteSG(ld, ln, objectRef string, fc FC, actSG int, value interface{}) error {
	// Set active setting group
	if err := c.Write(fmt.Sprintf(ActDA, ld, ln), SP, actSG); err != nil {
		return err
	}

	// Set edit setting group
	if err := c.Write(fmt.Sprintf(EditDA, ld, ln), SP, actSG); err != nil {
		return err
	}

	// Change a setting group value
	if err := c.Write(objectRef, fc, value); err != nil {
		return err
	}

	// Confirm new setting group values
	if err := c.Write(fmt.Sprintf(CnfDA, ld, ln), SP, true); err != nil {
		return err
	}
	return nil
}

// GetSG gets the SettingGroup
func (c *Client) GetSG(objectRef string) (*SettingGroup, error) {
	var clientError C.IedClientError
	cObjectRef, freeCObjectRef := allocCString(objectRef)
	defer freeCObjectRef()

	// Get type
	sgcbVarSpec := C.IedConnection_getVariableSpecification(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(SP))
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	defer C.MmsVariableSpecification_destroy(sgcbVarSpec)

	// Read SGCB
	sgcbVal := C.IedConnection_readObject(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(SP))
	if err := GetIedClientError(clientError); err != nil {
		return nil, err
	}
	//defer C.MmsValue_delete(sgcbVal)

	numOfSGValue, err := c.getSubElementValue(sgcbVal, sgcbVarSpec, "NumOfSG")
	if err != nil {
		return nil, err
	}

	actSGValue, err := c.getSubElementValue(sgcbVal, sgcbVarSpec, "ActSG")
	if err != nil {
		return nil, err
	}

	editSGValue, err := c.getSubElementValue(sgcbVal, sgcbVarSpec, "EditSG")
	if err != nil {
		return nil, err
	}

	cnfEditValue, err := c.getSubElementValue(sgcbVal, sgcbVarSpec, "CnfEdit")
	if err != nil {
		return nil, err
	}

	sg := &SettingGroup{
		NumOfSG: cast.ToInt(numOfSGValue),
		ActSG:   cast.ToInt(actSGValue),
		EditSG:  cast.ToInt(editSGValue),
		CnfEdit: cast.ToBool(cnfEditValue),
	}
	return sg, nil
}
