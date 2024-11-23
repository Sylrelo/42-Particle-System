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

type CL_PLATFORM_ID = C.cl_platform_id

type ClPlatform struct {
	ID      CL_PLATFORM_ID
	Name    string
	Vendor  string
	Version string
}

func GetAvailablePlatforms() ([]ClPlatform, error) {
	var numPlatforms C.cl_uint

	if errCl := C.clGetPlatformIDs(0, nil, &numPlatforms); errCl != C.CL_SUCCESS {
		fmt.Println("Error getting number of platforms")
		return nil, errors.New(ErrorString(int(errCl)))
	}

	platforms := make([]CL_PLATFORM_ID, numPlatforms)
	if C.clGetPlatformIDs(numPlatforms, &platforms[0], nil) != C.CL_SUCCESS {
		fmt.Println("Error getting platform IDs")
		return nil, errors.New("Error getting platforms IDs")
	}

	var platformList []ClPlatform
	for _, platform := range platforms {
		platformStruct := ClPlatform{
			ID:      platform,
			Name:    getPlatformInfoString(platform, C.CL_PLATFORM_NAME),
			Vendor:  getPlatformInfoString(platform, C.CL_PLATFORM_VENDOR),
			Version: getPlatformInfoString(platform, C.CL_PLATFORM_VERSION),
		}
		platformList = append(platformList, platformStruct)
	}

	return platformList, nil
}

func getPlatformInfoString(platform CL_PLATFORM_ID, paramName C.cl_platform_info) string {
	var size C.size_t
	clErr := C.clGetPlatformInfo(platform, paramName, 0, nil, &size)

	if clErr != C.CL_SUCCESS {
		return ErrorString(int(clErr))
	}

	buffer := make([]byte, size)

	clErr = C.clGetPlatformInfo(platform, paramName, size, unsafe.Pointer(&buffer[0]), nil)
	if clErr != C.CL_SUCCESS {
		return ErrorString(int(clErr))
	}

	return string(buffer[:size-1])
}
