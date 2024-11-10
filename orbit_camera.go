package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

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

	perspectiveMatrix mgl32.Mat4
	cameraMatrix      mgl32.Mat4
}

func (oc *OrbitCamera) RecalculatePerspectiveMatrix(width int, height int) {
	ratio := float32(width) / float32(height)

	oc.perspectiveMatrix = mgl32.Perspective(mgl32.DegToRad(45.0), ratio, 0.01, 10000)
}

func (oc *OrbitCamera) RecalculateCamera() {
	azimuthRad := oc.Azimuth * float32(math.Pi) / 180.0
	elevationRad := oc.Elevation * float32(math.Pi) / 180.0

	xOffset := oc.Distance * float32(math.Cos(float64(elevationRad))) * float32(math.Sin(float64(azimuthRad)))
	yOffset := oc.Distance * float32(math.Sin(float64(elevationRad)))
	zOffset := oc.Distance * float32(math.Cos(float64(elevationRad))) * float32(math.Cos(float64(azimuthRad)))

	oc.Position = mgl32.Vec3{
		oc.Target.X() + oc.TargetOffset[0] + xOffset,
		oc.Target.Y() + oc.TargetOffset[1] + yOffset,
		oc.Target.Z() + oc.TargetOffset[2] + zOffset,
	}

	actualTarget := oc.Target.Add(oc.TargetOffset)

	oc.cameraMatrix = mgl32.LookAtV(
		oc.Position,
		actualTarget,
		mgl32.Vec3{0, 1, 0},
	)
}

func (oc *OrbitCamera) HandleCameraRotation(x float32, y float32, storeDelta bool) {
	oc.Azimuth += x
	oc.Elevation += y

	if storeDelta {
		oc.RotateDelta = clampvec3n(
			oc.RotateDelta.Add(mgl32.Vec3{x, y, 0}),
			0,
			MAX_CAM_SMOOTH,
		)
	}

	oc.RecalculateCamera()
}

func (oc *OrbitCamera) HandleCameraDistance(z float32, storeDelta bool) {
	oc.Distance += float32(z)
	oc.Position[2] += float32(z)
	oc.RecalculateCamera()

	if storeDelta {
		oc.PanDelta[2] = clampn(oc.PanDelta[2]+z, MAX_CAM_SMOOTH, 0)
	}
}

func (oc *OrbitCamera) HandleCameraMovement(x float32, y float32, storeDelta bool) {
	forward := oc.Target.
		Add(oc.TargetOffset).
		Sub(oc.Position).
		Normalize()

	right := forward.
		Cross(oc.Up).
		Normalize()

	up := right.
		Cross(forward).
		Normalize()

	oc.TargetOffset = oc.TargetOffset.Add(right.Mul(x * 0.1)).Add(up.Mul(y * 0.1))

	if storeDelta {
		oc.PanDelta = clampvec3n(
			oc.PanDelta.Add(mgl32.Vec3{x, y, 0}),
			0,
			MAX_CAM_SMOOTH,
		)
	}

	oc.RecalculateCamera()
}

func (oc *OrbitCamera) SmoothMovement() {

	if abs(oc.RotateDelta[0]) > 0.1 || abs(oc.RotateDelta[1]) > 0.1 {
		valX := oc.RotateDelta[0] * 0.05
		valY := oc.RotateDelta[1] * 0.05

		oc.RotateDelta[0] -= valX
		oc.RotateDelta[1] -= valY

		oc.HandleCameraRotation(valX, valY, false)
	}

	if abs(oc.PanDelta[0]) > 0.1 || abs(oc.PanDelta[1]) > 0.1 || abs(oc.PanDelta[2]) > 0.1 {
		tmp := oc.PanDelta.Mul(0.05)

		oc.PanDelta = oc.PanDelta.Sub(tmp)

		oc.HandleCameraMovement(tmp[0], tmp[1], false)
		oc.HandleCameraDistance(tmp[2], false)
	}

}
