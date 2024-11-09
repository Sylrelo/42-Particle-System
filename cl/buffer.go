package clgo

/*
#include <OpenCL/opencl.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

func CreateSharedBuffer(clContext CL_CONTEXT, glVbo uint32) (CL_MEM, error) {
	var clErr C.cl_int

	clBuffer := C.clCreateFromGLBuffer(
		(C.cl_context)(unsafe.Pointer(clContext)),
		C.CL_MEM_READ_WRITE,
		C.cl_uint(glVbo),
		&clErr,
	)

	if clErr != C.CL_SUCCESS {
		return nil, errors.New(ErrorString(int(clErr)))
	}

	return clBuffer, nil

}
