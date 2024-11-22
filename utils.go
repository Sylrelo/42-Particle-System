package main

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

func ExitOnError(err error) {
	if err == nil {
		return
	}

	log.Fatalln(err)
}

func abs(in float32) float32 {
	if in < 0 {
		return in * -1
	}

	return in
}

func clamp(value float32, max float32, min float32) float32 {
	if value > max {
		return max
	}

	if value < min {
		return min
	}

	return value
}

func clampn(value float32, max float32, min float32) float32 {
	sign := 1.0
	if value < 0 {
		sign = -1.0
	}

	return clamp(abs(value), max, min) * float32(sign)
}

func clampvec3n(v mgl32.Vec3, min float32, max float32) mgl32.Vec3 {
	v[0] = clampn(v[0], max, min)
	v[1] = clampn(v[1], max, min)
	v[2] = clampn(v[2], max, min)

	return v
}

var QUAD_VERTICES = []float32{
	-0.5, -0.5, 0.0, 0.0, 0.0, // Bottom-left
	0.5, -0.5, 0.0, 1.0, 0.0, // Bottom-right
	0.5, 0.5, 0.0, 1.0, 1.0, // Top-right

	-0.5, -0.5, 0.0, 0.0, 0.0, // Bottom-left
	0.5, 0.5, 0.0, 1.0, 1.0, // Top-right
	-0.5, 0.5, 0.0, 0.0, 1.0, // Top-left
}

var BIG_QUAD_VERTICES = []float32{
	-1000.0, -1000.0, 0.0, 0.0, 0.0, // Bottom-left
	1000.0, -1000.0, 0.0, 1.0, 0.0, // Bottom-right
	1000.0, 1000.0, 0.0, 1.0, 1.0, // Top-right

	-1000.0, -1000.0, 0.0, 0.0, 0.0, // Bottom-left
	1000.0, 1000.0, 0.0, 1.0, 1.0, // Top-right
	-1000.0, 1000.0, 0.0, 0.0, 1.0, // Top-left
}
