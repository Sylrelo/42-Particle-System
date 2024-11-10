package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	MAX_CAM_SMOOTH = 100
)

func init() {
	runtime.LockOSThread()
}

type OrbitCamera struct {
	Position     mgl32.Vec3
	Target       mgl32.Vec3
	TargetOffset mgl32.Vec3
	Up           mgl32.Vec3
	Azimuth      float32
	Elevation    float32
	Distance     float32

	PanDelta    mgl32.Vec3
	RotateDelta mgl32.Vec3
}

type InputData struct {
	mouseButton [8]glfw.Action
	_oldMouseX  float64
	_oldMouseY  float64
}

type ParticleSystem struct {
	windowHeight int
	windowWidth  int

	cameraMatrix mgl32.Mat4
	orbitCamera  OrbitCamera

	inputs InputData
}

func (ps *ParticleSystem) SmoothMovement() {

	if abs(ps.orbitCamera.RotateDelta[0]) > 0.1 || abs(ps.orbitCamera.RotateDelta[1]) > 0.1 {
		valX := ps.orbitCamera.RotateDelta[0] * 0.05
		valY := ps.orbitCamera.RotateDelta[1] * 0.05

		ps.orbitCamera.RotateDelta[0] -= valX
		ps.orbitCamera.RotateDelta[1] -= valY

		ps.HandleCameraRotation(valX, valY, false)
	}

	if abs(ps.orbitCamera.PanDelta[0]) > 0.1 || abs(ps.orbitCamera.PanDelta[1]) > 0.1 {
		tmp := ps.orbitCamera.PanDelta.Mul(0.05)

		ps.orbitCamera.PanDelta = ps.orbitCamera.PanDelta.Sub(tmp)

		ps.HandleCameraMovement(tmp[0], tmp[1], false)
	}
}

func (ps *ParticleSystem) RecalculateCamera() {
	azimuthRad := ps.orbitCamera.Azimuth * float32(math.Pi) / 180.0
	elevationRad := ps.orbitCamera.Elevation * float32(math.Pi) / 180.0

	xOffset := ps.orbitCamera.Distance * float32(math.Cos(float64(elevationRad))) * float32(math.Sin(float64(azimuthRad)))
	yOffset := ps.orbitCamera.Distance * float32(math.Sin(float64(elevationRad)))
	zOffset := ps.orbitCamera.Distance * float32(math.Cos(float64(elevationRad))) * float32(math.Cos(float64(azimuthRad)))

	ps.orbitCamera.Position = mgl32.Vec3{
		ps.orbitCamera.Target.X() + ps.orbitCamera.TargetOffset[0] + xOffset,
		ps.orbitCamera.Target.Y() + ps.orbitCamera.TargetOffset[1] + yOffset,
		ps.orbitCamera.Target.Z() + ps.orbitCamera.TargetOffset[2] + zOffset,
	}

	actualTarget := ps.orbitCamera.Target.Add(ps.orbitCamera.TargetOffset)

	ps.cameraMatrix = mgl32.LookAtV(
		ps.orbitCamera.Position,
		actualTarget,
		mgl32.Vec3{0, 1, 0},
	)
}

func (ps *ParticleSystem) HandleCameraRotation(x float32, y float32, storeDelta bool) {
	ps.orbitCamera.Azimuth += x
	ps.orbitCamera.Elevation += y

	if storeDelta {
		ps.orbitCamera.RotateDelta = clampvec3n(
			ps.orbitCamera.RotateDelta.Add(mgl32.Vec3{x, y, 0}),
			0,
			MAX_CAM_SMOOTH,
		)
	}

	ps.RecalculateCamera()
}

func (ps *ParticleSystem) HandleCameraMovement(x float32, y float32, storeDelta bool) {
	forward := ps.orbitCamera.Target.
		Add(ps.orbitCamera.TargetOffset).
		Sub(ps.orbitCamera.Position).
		Normalize()

	right := forward.
		Cross(ps.orbitCamera.Up).
		Normalize()

	up := right.
		Cross(forward).
		Normalize()

	ps.orbitCamera.TargetOffset = ps.orbitCamera.TargetOffset.Add(right.Mul(x * 0.1)).Add(up.Mul(y * 0.1))

	if storeDelta {
		ps.orbitCamera.PanDelta = clampvec3n(
			ps.orbitCamera.PanDelta.Add(mgl32.Vec3{x, y, 0}),
			0,
			MAX_CAM_SMOOTH,
		)
	}

	ps.RecalculateCamera()
}

func (inp *ParticleSystem) MouseButtonEvent(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	inp.inputs.mouseButton[button] = action
}

func (ps *ParticleSystem) MouseMotionEvent(w *glfw.Window, xpos float64, ypos float64) {
	diffX := (ps.inputs._oldMouseX - xpos) * 0.05
	diffY := (ps.inputs._oldMouseY - -ypos) * 0.05
	ps.inputs._oldMouseX = xpos
	ps.inputs._oldMouseY = -ypos

	if ps.inputs.mouseButton[glfw.MouseButtonLeft] == glfw.Press {
		ps.HandleCameraRotation(float32(diffX), float32(diffY), true)
	}

	if ps.inputs.mouseButton[glfw.MouseButtonRight] == glfw.Press {
		ps.HandleCameraMovement(float32(diffX), float32(diffY), true)
	}

}

func (ps *ParticleSystem) MouseScrollEvent(w *glfw.Window, xoff float64, yoff float64) {
	ps.orbitCamera.Distance += float32(yoff)
	ps.orbitCamera.Position[2] += float32(yoff)
	ps.RecalculateCamera()
}

func (inp *ParticleSystem) WindowResizeEvent(w *glfw.Window, width int, height int) {
	inp.windowHeight = height
	inp.windowWidth = width
}

///

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	system := ParticleSystem{}

	system.orbitCamera = OrbitCamera{
		Position:     mgl32.Vec3{0, 0, 10},
		Target:       mgl32.Vec3{0, 0, 0},
		TargetOffset: mgl32.Vec3{0, 0, 0},
		Up:           mgl32.Vec3{0, 1, 0},
		Azimuth:      0,
		Elevation:    0,
		Distance:     10,
	}
	system.RecalculateCamera()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(1280, 720, "OpenCL-OpenGL Interop", nil, nil)
	if err != nil {
		log.Fatalln("failed to create glfw window:", err)
	}
	window.MakeContextCurrent()

	window.SetSizeCallback(system.WindowResizeEvent)
	window.SetCursorPosCallback(system.MouseMotionEvent)
	window.SetMouseButtonCallback(system.MouseButtonEvent)
	window.SetScrollCallback(system.MouseScrollEvent)

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
	// cameraMatrix := mgl32.Ident4()

	// perspectiveMatrix = mgl32.Ident4()

	perspectiveUniform := gl.GetUniformLocation(program, gl.Str("perspective\x00"))
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))

	gl.UseProgram(program)
	gl.UniformMatrix4fv(perspectiveUniform, 1, false, &perspectiveMatrix[0])

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
		gl.UniformMatrix4fv(cameraUniform, 1, false, &system.cameraMatrix[0])

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.Finish()

		cl_compute.RunGravitateKernel(particlesCount)

		system.SmoothMovement()

		// gl.UseProgram(program)
		// gl.EnableVertexAttribArray(0)
		// gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)

		// gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.POINTS, 0, int32(particlesCount))

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
