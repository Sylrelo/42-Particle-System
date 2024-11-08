package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// glfw.WindowHint(glfw.AlphaBits, 8)

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

	particlesCount := 1000000

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, particlesCount*4*4, nil, gl.DYNAMIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	///

	var vbo_velocity uint32
	gl.GenBuffers(1, &vbo_velocity)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo_velocity)
	gl.BufferData(gl.ARRAY_BUFFER, particlesCount*4*4, nil, gl.DYNAMIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// fmt.Println(gl.CreateShader(gl.COMPUTE_SHADER))

	fmt.Println("OpenGL VBO created:", vbo)

	cl_compute := InitClCompute(vbo, vbo_velocity)

	fmt.Println("OpenCL-OpenGL interop buffer created successfully on macOS")

	////////////////////////
	perspectiveMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), 16.0/9.0, 0.01, 10000)
	cameraMatrix := mgl32.LookAtV(mgl32.Vec3{0, 0, 10}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	// cameraMatrix := mgl32.Ident4()

	// perspectiveMatrix = mgl32.Ident4()

	perspectiveUniform := gl.GetUniformLocation(program, gl.Str("perspective\x00"))
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))

	gl.UseProgram(program)
	gl.UniformMatrix4fv(perspectiveUniform, 1, false, &perspectiveMatrix[0])
	gl.UniformMatrix4fv(cameraUniform, 1, false, &cameraMatrix[0])

	cl_compute.RunInitKernel(particlesCount)

	var vao uint32
	gl.UseProgram(program)

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo_velocity)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, nil)

	// gl.EnableVertexArrayAttrib(0)

	// gl.PointSize(2)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0, 0, 0, 1.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.Finish()

		cl_compute.RunGravitateKernel(particlesCount)
		// gl.UseProgram(program)
		// gl.EnableVertexAttribArray(0)
		// gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)

		// gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
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
