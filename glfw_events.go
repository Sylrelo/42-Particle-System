package main

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type InputData struct {
	mouseButton [8]glfw.Action
	_oldMouseX  float64
	_oldMouseY  float64
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
		ps.orbitCamera.HandleCameraRotation(float32(diffX), float32(diffY), true)
	}

	if ps.inputs.mouseButton[glfw.MouseButtonRight] == glfw.Press {
		ps.orbitCamera.HandleCameraMovement(float32(diffX), float32(diffY), true)
	}

}

func (ps *ParticleSystem) MouseScrollEvent(w *glfw.Window, xoff float64, yoff float64) {
	ps.orbitCamera.HandleCameraDistance(float32(yoff), true)
}

func (inp *ParticleSystem) WindowResizeEvent(w *glfw.Window, width int, height int) {
	inp.windowHeight = height
	inp.windowWidth = width

	inp.orbitCamera.RecalculatePerspectiveMatrix(width, height)
	gl.Viewport(0, 0, int32(width), int32(height))

}
