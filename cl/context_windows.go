//go:build windows
// +build windows

package clgo

/*
#include <OpenCL/opencl.h>

#include <OpenGL/gl.h>
#include <OpenGL/wgl.h>
#include <OpenGL/wglext.h>


cl_context CreateSharedCLGLContext(cl_device_id *device, HGLRC glContext, HDC hDC) {
    cl_int 					err;
    cl_platform_id 	platform;
    // cl_device_id 		devices[1];
    size_t 					deviceSize;

    cl_context_properties props[] = {
        CL_CONTEXT_PLATFORM, (cl_context_properties)platform,
        CL_GL_CONTEXT_KHR, (cl_context_properties)glContext,
        CL_WGL_HDC_KHR, (cl_context_properties)hDC,
        0
    };

    return clCreateContext(props, 1, &device[0], NULL, NULL, &err);
}

*/
import "C"
import (
	"errors"
	"unsafe"
)

type CL_CONTEXT = C.cl_context

func CreateSharedOpenglContext(platform ClPlatform, device ClDevice) (CL_CONTEXT, error) {
	var properties []C.cl_context_properties

	glContext := C.wglGetCurrentContext()
	hDC := C.wglGetCurrentDC()

	properties = []C.cl_context_properties{
		C.CL_CONTEXT_PLATFORM, C.cl_context_properties(uintptr(unsafe.Pointer(platform.ID))),
		C.CL_GL_CONTEXT_KHR, C.cl_context_properties(uintptr(unsafe.Pointer(glContext))),
		C.CL_WGL_HDC_KHR, C.cl_context_properties(uintptr(unsafe.Pointer(hDC))),
		0,
	}

	var clErr C.cl_int
	clContext := C.clCreateContext(
		(*C.cl_context_properties)(&properties[0]),
		1, &device.ID,
		nil, nil,
		&clErr,
	)

	if clErr != C.CL_SUCCESS {
		return nil, errors.New(ErrorString(int(clErr)))
	}

	return clContext, nil
}
