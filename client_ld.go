package iec61850

// #include <iec61850_client.h>
import "C"
import (
	"fmt"
)

func (c *Client) GetLogicalDeviceList() DataModel {
	var clientError C.IedClientError
	deviceList := C.IedConnection_getLogicalDeviceList(c.conn, &clientError)

	var dataModel DataModel

	device := deviceList.next
	for device != nil {

		var ld LD
		ld.Data = C2GoStr((*C.char)(device.data))

		logicalNodes := C.IedConnection_getLogicalDeviceDirectory(c.conn, &clientError, (*C.char)(device.data))
		logicalNode := logicalNodes.next

		for logicalNode != nil {
			func() {
				var ln LN
				ln.Data = C2GoStr((*C.char)(logicalNode.data))
				lnRef := fmt.Sprintf("%s/%s", ld.Data, C2GoStr((*C.char)(logicalNode.data)))
				cRef, freeCRef := allocGo2CStr(lnRef)
				defer freeCRef()
				dataObjects := C.IedConnection_getLogicalNodeDirectory(c.conn, &clientError, cRef, C.ACSI_CLASS_DATA_OBJECT)
				dataObject := dataObjects.next
				for dataObject != nil {
					var do DO
					do.Data = C2GoStr((*C.char)(dataObject.data))

					dataObject = dataObject.next
					doRef := fmt.Sprintf("%s/%s.%s", C2GoStr((*C.char)(device.data)), C2GoStr((*C.char)(logicalNode.data)), do.Data)

					var das []DA
					c.GetDAs(doRef, das)

					do.DAs = das
					ln.DOs = append(ln.DOs, do)
				}
				C.LinkedList_destroy(dataObjects)
				clnRef, freeClnRef := allocGo2CStr(lnRef)
				defer freeClnRef()
				dataSets := C.IedConnection_getLogicalNodeDirectory(c.conn, &clientError, clnRef, C.ACSI_CLASS_DATA_SET)
				dataSet := dataSets.next
				for dataSet != nil {
					func() {
						var ds DS
						ds.Data = C2GoStr((*C.char)(dataSet.data))
						var isDeletable C.bool
						dataSetRef := fmt.Sprintf("%s.%s", lnRef, ds.Data)
						cdataSetRef, freeCdataSetRef := allocGo2CStr(dataSetRef)
						defer freeCdataSetRef()

						dataSetMembers := C.IedConnection_getDataSetDirectory(c.conn, &clientError, cdataSetRef, &isDeletable)
						if isDeletable {
							fmt.Printf("    Data set: %s (deletable)\n", ds.Data)
						} else {
							fmt.Printf("    Data set: %s (not deletable)\n", ds.Data)
						}
						dataSetMemberRef := dataSetMembers.next
						for dataSetMemberRef != nil {
							var dsRef DSRef
							dsRef.Data = C2GoStr((*C.char)(dataSetMemberRef.data))
							ds.DSRefs = append(ds.DSRefs, dsRef)

							dataSetMemberRef = dataSetMemberRef.next
						}
						C.LinkedList_destroy(dataSetMembers)
						dataSet = dataSet.next
						ln.DSs = append(ln.DSs, ds)
					}()
				}
				C.LinkedList_destroy(dataSets)

				clnRef1, freeClnRef1 := allocGo2CStr(lnRef)
				defer freeClnRef1()

				reports := C.IedConnection_getLogicalNodeDirectory(c.conn, &clientError, clnRef1, C.ACSI_CLASS_URCB)
				report := reports.next
				for report != nil {
					var r URReport
					r.Data = C2GoStr((*C.char)(report.data))
					ln.URReports = append(ln.URReports, r)

					report = report.next
				}
				C.LinkedList_destroy(reports)

				clnRef2, freeClnRef2 := allocGo2CStr(lnRef)
				defer freeClnRef2()

				reports = C.IedConnection_getLogicalNodeDirectory(c.conn, &clientError, clnRef2, C.ACSI_CLASS_BRCB)
				report = reports.next
				for report != nil {
					var r BRReport
					r.Data = C2GoStr((*C.char)(report.data))
					ln.BRReports = append(ln.BRReports, r)

					report = report.next
				}

				C.LinkedList_destroy(reports)

				ld.LNs = append(ld.LNs, ln)

				logicalNode = logicalNode.next
			}()
		}
		C.LinkedList_destroy(logicalNodes)

		dataModel.LDs = append(dataModel.LDs, ld)

		device = device.next
	}
	C.LinkedList_destroy(deviceList)
	return dataModel
}

func (c *Client) GetDAs(doRef string, das []DA) {

	var clientError C.IedClientError

	cdoRef, freeCdoRef := allocGo2CStr(doRef)
	defer freeCdoRef()

	dataAttributes := C.IedConnection_getDataDirectory(c.conn, &clientError, cdoRef)
	defer C.LinkedList_destroy(dataAttributes)
	if dataAttributes != nil {
		dataAttribute := dataAttributes.next

		for dataAttribute != nil {
			var da DA
			da.Data = C2GoStr((*C.char)(dataAttribute.data))
			das = append(das, da)

			dataAttribute = dataAttribute.next
			daRef := fmt.Sprintf("%s.%s", doRef, da.Data)
			c.GetDAs(daRef, das)
		}
	}

}
