package clgo

/*
#cgo CFLAGS: -framework OpenCL -framework OpenGL
#cgo LDFLAGS: -framework OpenCL -framework OpenGL
#include <OpenCL/opencl.h>
#include <OpenGL/gl.h>
#include <OpenGL/OpenGL.h>
#include <OpenGL/CGLCurrent.h>
*/
import "C"
