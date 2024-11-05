package clgo

/*
#include <OpenCL/opencl.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type CL_DEVICE_ID = C.cl_device_id

type ClDevice struct {
	ID     CL_DEVICE_ID
	Name   string
	Vendor string
}

func GetAvailableDevices(platform ClPlatform) ([]ClDevice, error) {
	var numDevices C.cl_uint
	var devicesList []ClDevice

	if C.clGetDeviceIDs(platform.ID, C.CL_DEVICE_TYPE_ALL, 0, nil, &numDevices) != C.CL_SUCCESS {
		fmt.Println("Error getting number of devices")
		return nil, errors.New("Error getting number of devices")
	}

	deviceIDs := make([]C.cl_device_id, numDevices)

	err := C.clGetDeviceIDs(platform.ID, C.CL_DEVICE_TYPE_ALL, numDevices, &deviceIDs[0], nil)
	if err != C.CL_SUCCESS {
		return nil, errors.New(ErrorString(int(err)))
	}

	for _, deviceId := range deviceIDs {
		device := ClDevice{
			ID:     deviceId,
			Name:   getDeviceInfoString(deviceId, C.CL_DEVICE_NAME),
			Vendor: getDeviceInfoString(deviceId, C.CL_DEVICE_VENDOR),
		}

		devicesList = append(devicesList, device)
	}
	return devicesList, nil
}

func getDeviceInfoString(device C.cl_device_id, paramName C.cl_device_info) string {
	var size C.size_t
	if errcode := C.clGetDeviceInfo(device, paramName, 0, nil, &size); errcode != C.CL_SUCCESS {
		return ErrorString(int(errcode))
	}

	buffer := make([]byte, size)
	if errcode := C.clGetDeviceInfo(device, paramName, size, unsafe.Pointer(&buffer[0]), nil); errcode != C.CL_SUCCESS {
		return ErrorString(int(errcode))
	}

	return string(buffer[:size-1])
}
