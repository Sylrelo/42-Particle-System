package main

/*
#cgo CFLAGS: -framework OpenCL -framework OpenGL
#cgo LDFLAGS: -framework OpenCL -framework OpenGL
#include <OpenCL/opencl.h>
#include <OpenGL/gl.h>
#include <OpenGL/OpenGL.h>
#include <OpenGL/CGLCurrent.h>
*/
import "C"
import (
	"clgo"
	"fmt"
	"log"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Platform struct {
	ID      C.cl_platform_id
	Name    string
	Vendor  string
	Version string
}

func init() {
	runtime.LockOSThread()
}

func main() {
	var clErr C.cl_int

	// Initialize GLFW for OpenGL context creation
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "OpenCL-OpenGL Interop", nil, nil)
	if err != nil {
		log.Fatalln("failed to create glfw window:", err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize gl:", err)
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 1024*4, nil, gl.DYNAMIC_DRAW)

	fmt.Println("OpenGL VBO created:", vbo)

	clPlatforms, err := clgo.GetAvailablePlatforms()
	if err != nil {
		log.Fatalln(clPlatforms)
	}

	clContext, err := clgo.CreateSharedOpenglContext()
	if err != nil {
		log.Fatalln(err)
	}

	clDevices, err := clgo.GetAvailableDevices(clPlatforms[0])
	if err != nil {
		log.Fatalln(err)
	}

	clQueue, err := clgo.CreateCommandQueue(clContext, clDevices[0])
	if err != nil {
		log.Fatalln(err)
	}

	// Step 5: Create OpenCL Buffer from OpenGL Buffer
	clBuffer := C.clCreateFromGLBuffer(
		(C.cl_context)(unsafe.Pointer(clContext)),
		C.CL_MEM_READ_WRITE,
		C.GLuint(vbo),
		&clErr,
	)
	if clErr != C.CL_SUCCESS {
		fmt.Println("Failed to create OpenCL buffer from OpenGL buffer")
		return
	}

	fmt.Println("OpenCL-OpenGL interop buffer created successfully on macOS")

	// Main Loop: Rendering and OpenCL Operations
	for !window.ShouldClose() {
		clQueue.AcquireGLObjects((clgo.CL_MEM)(unsafe.Pointer(clBuffer)))
		clQueue.ReleaseGLObjects((clgo.CL_MEM)(unsafe.Pointer(clBuffer)))
		clQueue.Finish()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		window.SwapBuffers()
		glfw.PollEvents()
	}

	//var numPlatforms C.cl_uint
	//if C.clGetPlatformIDs(0, nil, &numPlatforms) != C.CL_SUCCESS {
	//	fmt.Println("Error getting number of platforms")
	//	return
	//}
	//
	//platforms := make([]C.cl_platform_id, numPlatforms)
	//if C.clGetPlatformIDs(numPlatforms, &platforms[0], nil) != C.CL_SUCCESS {
	//	fmt.Println("Error getting platform IDs")
	//	return
	//}
	//
	//var platformList []Platform
	//for _, platform := range platforms {
	//	platformStruct := Platform{
	//		ID:      platform,
	//		Name:    getPlatformInfoString(platform, C.CL_PLATFORM_NAME),
	//		Vendor:  getPlatformInfoString(platform, C.CL_PLATFORM_VENDOR),
	//		Version: getPlatformInfoString(platform, C.CL_PLATFORM_VERSION),
	//	}
	//	platformList = append(platformList, platformStruct)
	//}

}

func getPlatformInfoString(platform C.cl_platform_id, paramName C.cl_platform_info) string {
	var size C.size_t
	if C.clGetPlatformInfo(platform, paramName, 0, nil, &size) != C.CL_SUCCESS {
		return ""
	}

	buffer := make([]byte, size)
	if C.clGetPlatformInfo(platform, paramName, size, unsafe.Pointer(&buffer[0]), nil) != C.CL_SUCCESS {
		return ""
	}

	return string(buffer[:size-1])
}
