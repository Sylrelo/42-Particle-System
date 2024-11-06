package clgo

/*
#cgo darwin CFLAGS: -framework OpenCL -framework OpenGL
#cgo darwin LDFLAGS: -framework OpenCL -framework OpenGL

#cgo linux CFLAGS: -Iincludes/
#cgo linux LDFLAGS: -static

#cgo windows CFLAGS: -Iincludes/
#cgo windows LDFLAGS: -static

#include <OpenCL/opencl.h>
#include <OpenGL/gl.h>
#include <OpenGL/OpenGL.h>
#include <OpenGL/CGLCurrent.h>
*/
import "C"
