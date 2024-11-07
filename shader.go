package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gl/gl/all-core/gl"
)

type GL_SHADER_TYPE = uint32
type GL_SHADER = uint32
type GL_PROGRAM = uint32

const (
	FRAGMENT GL_SHADER_TYPE = gl.FRAGMENT_SHADER
	VERTEX   GL_SHADER_TYPE = gl.VERTEX_SHADER
	COMPUTE  GL_SHADER_TYPE = gl.COMPUTE_SHADER
)

func CompileShader(path string, shaderType GL_SHADER_TYPE) (GL_SHADER, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read shader file: %v", err)
	}

	source := string(data) + "\x00"

	shader := gl.CreateShader(shaderType)
	cstr, free := gl.Strs(source)
	defer free()

	gl.ShaderSource(shader, 1, cstr, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength+1)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])

		return 0, fmt.Errorf("failed to compile shader: %s", log)
	}

	return shader, nil
}

func CreateProgram(shaders ...GL_SHADER) (uint32, error) {

	defer func() {
		for _, shader := range shaders {
			gl.DeleteShader(shader)
		}
	}()

	program := gl.CreateProgram()

	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength+1)
		gl.GetProgramInfoLog(program, logLength, nil, &log[0])
		return 0, fmt.Errorf("failed to link program: %s", log)
	}

	return program, nil
}
