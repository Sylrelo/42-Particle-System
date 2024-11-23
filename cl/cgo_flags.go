package clgo

/*

#cgo darwin CFLAGS: -framework OpenCL -framework OpenGL
#cgo darwin LDFLAGS: -framework OpenCL -framework OpenGL

#cgo linux CFLAGS: -Iincludes/
#cgo linux LDFLAGS: -L/usr/lib/x86_64-linux-gnu -Llib/ -Llib/vendors -lOpenCL

#cgo windows CFLAGS: -Iincludes/
#cgo windows LDFLAGS: -static -Llib/ -lOpenCL

#ifdef __linux__
#include <OpenCL/opencl.h>
#include <GL/gl.h>
#endif

#ifdef __windows__
#include <OpenCL/opencl.h>
#include <OpenGL/gl.h>
#endif

#ifdef __APPLE__
#include <OpenGL/OpenGL.h>
#include <OpenGL/CGLCurrent.h>
#endif


*/
import "C"
