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
	"os"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
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

// func setInitKernelArgs() {

// }

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

	window, err := glfw.CreateWindow(1280, 720, "OpenCL-OpenGL Interop", nil, nil)
	if err != nil {
		log.Fatalln("failed to create glfw window:", err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize gl:", err)
	}

	currentPath, _ := os.Getwd()

	renderFragmentShader, err := CompileShader(currentPath+"/shaders/base.frag", FRAGMENT)
	ExitOnError(err)
	renderVertexShader, err := CompileShader(currentPath+"/shaders/base.vert", VERTEX)
	ExitOnError(err)
	program, err := CreateProgram(renderFragmentShader, renderVertexShader)
	ExitOnError(err)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	particlesCount := 500000
	gl.BufferData(gl.ARRAY_BUFFER, particlesCount*4*4, nil, gl.DYNAMIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// fmt.Println(gl.CreateShader(gl.COMPUTE_SHADER))

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

	clKernel, err := clgo.CreateKernel(clContext, clDevices[0], currentPath+"/compute/kernel.cl", "initParticles")
	ExitOnError(err)

	_ = clKernel

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

	////////////////////////
	perspectiveMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), 16.0/9.0, 0.01, 10000)
	cameraMatrix := mgl32.LookAtV(mgl32.Vec3{0, 0, 4}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	// cameraMatrix := mgl32.Ident4()

	// perspectiveMatrix = mgl32.Ident4()

	perspectiveUniform := gl.GetUniformLocation(program, gl.Str("perspective\x00"))
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))

	gl.UseProgram(program)
	gl.UniformMatrix4fv(perspectiveUniform, 1, false, &perspectiveMatrix[0])
	gl.UniformMatrix4fv(cameraUniform, 1, false, &cameraMatrix[0])

	///////////////////////
	gl.Finish()
	errCl := clgo.SetKernelArgs(clKernel.Kernel, (clgo.CL_MEM)(unsafe.Pointer(clBuffer)))
	ExitOnError(errCl)
	errCl = clQueue.AcquireGLObjects((clgo.CL_MEM)(unsafe.Pointer(clBuffer)))
	ExitOnError(errCl)
	errCl = clQueue.EnqueueKernel(clKernel.Kernel, particlesCount)
	ExitOnError(errCl)
	errCl = clQueue.ReleaseGLObjects((clgo.CL_MEM)(unsafe.Pointer(clBuffer)))
	ExitOnError(errCl)
	clQueue.Finish()

	// _ = clBuffer
	// _ = clQueue
	///////////////////////

	var vao uint32
	gl.UseProgram(program)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	// gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)

	// gl.EnableVertexArrayAttrib(0)

	// gl.PointSize(2)

	for !window.ShouldClose() {
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// clQueue.Finish()
		// gl.UseProgram(program)
		// gl.EnableVertexAttribArray(0)
		// gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)

		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.POINTS, 0, int32(particlesCount))

		// clQueue.AcquireGLObjects((clgo.CL_MEM)(unsafe.Pointer(clBuffer)))
		// clQueue.ReleaseGLObjects((clgo.CL_MEM)(unsafe.Pointer(clBuffer)))
		// clQueue.Finish()

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
