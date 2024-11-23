//go:build linux
// +build linux

package clgo

/*
#include <OpenCL/opencl.h>

#ifdef __linux__
#include <GL/gl.h>
#include <GL/glext.h>
#endif

#ifdef __APPLE__
#include <OpenGL/OpenGL.h>
#include <OpenGL/CGLCurrent.h>
#endif
*/
import "C"
import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type CL_CONTEXT = C.cl_context


func CreateSharedOpenglContext(platform ClPlatform, device ClDevice) (CL_CONTEXT, error) {
	_ = platform
	_ = device

	// currentContext := glfw.GetCurrentContext();
	
	// errCl := C.clGetGLContextInfoKHR(
	// 	(*C.cl_context_properties)(unsafe.Pointer(&platform.ID)), 
	// 	C.CL_CONTEXT_PLATFORM, 
	// 	C.size_t(uintptr(unsafe.Pointer(currentContext))), 
	// 	C.size_t(unsafe.Pointer(1)),
	// 	nil,
	// )

	// if errCl != C.CL_SUCCESS {
	// 	return nil, errors.New(ErrorString(int(errCl)))
	// }

	display := glfw.GetX11Display()
	ctx := glfw.GetCurrentContext()
	_ = display
	_ = ctx


	// aled2 := C.glXGetCurrentDisplay()
	// aled := C.glXGetCurrentContext()

	// properties := []C.cl_context_properties{
	// 	C.CL_GL_CONTEXT_KHR, C.cl_context_properties(uintptr(unsafe.Pointer(currentContext))),
	// 	C.CL_CONTEXT_PLATFORM, C.cl_context_properties(uintptr(unsafe.Pointer(platform.ID))),
	// 	0,
	// }

	// clContext := C.clCreateContext(
	// 	(*C.cl_context_properties)(&properties[0]),
	// 	0, nil, nil, nil, &errCl,
	// )

	// if errCl != C.CL_SUCCESS {
	// 	return nil, errors.New(ErrorString(errCl))
	// }

	return nil, nil
	// return clContext, nil
}