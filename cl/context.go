package clgo

/*
#include <OpenCL/opencl.h>
#include <OpenGL/gl.h>
#include <OpenGL/OpenGL.h>
#include <OpenGL/CGLCurrent.h>
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type CL_CONTEXT = C.cl_context

func CreateSharedOpenglContext() (CL_CONTEXT, error) {
	var properties []C.cl_context_properties

	cglContext := C.CGLGetCurrentContext()
	cglShareGroup := C.CGLGetShareGroup(cglContext)

	if runtime.GOOS == "darwin" {
		properties = []C.cl_context_properties{
			C.cl_context_properties(C.CL_CONTEXT_PROPERTY_USE_CGL_SHAREGROUP_APPLE),
			C.cl_context_properties(uintptr(unsafe.Pointer(cglShareGroup))),
			0,
		}
	} else {
		return nil, errors.ErrUnsupported
	}

	var clErr C.cl_int
	clContext := C.clCreateContext(
		(*C.cl_context_properties)(&properties[0]),
		0, nil, nil, nil, &clErr,
	)

	if clErr != C.CL_SUCCESS {
		return nil, errors.New(ErrorString(int(clErr)))
	}

	return clContext, nil
}
