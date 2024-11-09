//go:build darwin
// +build darwin

package clgo

/*
#include <OpenCL/opencl.h>
#include <OpenGL/gl.h>

#ifdef __APPLE__
#include <OpenGL/OpenGL.h>
#include <OpenGL/CGLCurrent.h>
#endif
*/
import "C"
import (
	"errors"
	"unsafe"
)

type CL_CONTEXT = C.cl_context

func CreateSharedOpenglContext(platform ClPlatform, device ClDevice) (CL_CONTEXT, error) {
	_ = platform
	_ = device

	var properties []C.cl_context_properties

	cglContext := C.CGLGetCurrentContext()
	cglShareGroup := C.CGLGetShareGroup(cglContext)

	properties = []C.cl_context_properties{
		C.cl_context_properties(C.CL_CONTEXT_PROPERTY_USE_CGL_SHAREGROUP_APPLE),
		C.cl_context_properties(uintptr(unsafe.Pointer(cglShareGroup))),
		0,
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
