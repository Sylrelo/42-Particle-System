package main

import (
	"clgo"
	"log"
	"os"
	"path"
	"unsafe"
)

type ComputeCL struct {
	queue       clgo.ClQueue
	device      clgo.ClDevice
	initProgram clgo.ClKernel

	clContext clgo.CL_CONTEXT

	clPositionBuffer clgo.CL_MEM
	clVelocityBuffer clgo.CL_MEM

	clProgram         clgo.CL_PROGRAM
	clInitKernel      clgo.CL_KERNEL
	clIdleKernel      clgo.CL_KERNEL
	clGravitateKernel clgo.CL_KERNEL
	clTransmove       clgo.CL_KERNEL
}

func InitClCompute(glPositionBuffer uint32, glVelocityBuffer uint32) ComputeCL {
	currentPath, err := os.Getwd()
	ExitOnError(err)
	INIT_KERNEL_PATH := path.Join(currentPath, "resources/compute/kernel.cl")

	clPlatforms, err := clgo.GetAvailablePlatforms()
	if err != nil {
		log.Fatalln(clPlatforms)
	}

	log.Println("Available OpenCL Platforms ====================")
	for _, platforn := range clPlatforms {
		log.Println(platforn.ID, platforn.Name)
	}

	clDevices, err := clgo.GetAvailableDevices(clPlatforms[0])
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Available OpenCL Devices ====================")
	for _, device := range clDevices {
		log.Println(device.ID, device.Name)
	}

	clContext, err := clgo.CreateSharedOpenglContext(clPlatforms[0], clDevices[0])
	if err != nil {
		log.Fatalln("CreateSharedOpenglContext", err)
	}
	log.Println("Shared OpenCL-GL Context Created.")

	clQueue, err := clgo.CreateCommandQueue(clContext, clDevices[0])
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("OpenCL Command Queue Created.")

	initProgram, err := clgo.InitKernels(clContext, clDevices[0], INIT_KERNEL_PATH)
	ExitOnError(err)

	clPositionBuffer, err := clgo.CreateSharedBuffer(clContext, glPositionBuffer)
	ExitOnError(err)
	clVelocityBuffer, err := clgo.CreateSharedBuffer(clContext, glVelocityBuffer)
	ExitOnError(err)

	return ComputeCL{
		clInitKernel:      initProgram.Kernels["initParticles"],
		clGravitateKernel: initProgram.Kernels["gravitateParticles"],
		clTransmove:       initProgram.Kernels["transmove"],

		device:      clDevices[0],
		clContext:   clContext,
		queue:       clQueue,
		initProgram: initProgram,

		clPositionBuffer: clPositionBuffer,
		clVelocityBuffer: clVelocityBuffer,
	}

}

/* -------------------------------------------------------------------------- */
/*                                 RUN KERNELS                                */
/* -------------------------------------------------------------------------- */

/* ---------------------------------- INIT ---------------------------------- */

func (ccl *ComputeCL) RunInitKernel(particleCount int) error {
	err := clgo.SetKernelArgs(ccl.clInitKernel, ccl.clPositionBuffer, ccl.clVelocityBuffer)
	if err != nil {
		return err
	}

	glObjects := []clgo.CL_MEM{
		(clgo.CL_MEM)(unsafe.Pointer(ccl.clPositionBuffer)),
		(clgo.CL_MEM)(unsafe.Pointer(ccl.clVelocityBuffer)),
	}

	err = ccl.queue.AcquireGLObjects(glObjects)
	if err != nil {
		return err
	}
	err = ccl.queue.EnqueueKernel(ccl.clInitKernel, particleCount)
	if err != nil {
		return err
	}

	err = ccl.queue.ReleaseGLObjects(glObjects)
	if err != nil {
		return err
	}

	ccl.queue.Finish()
	return nil
}

/* ---------------------------------- MOVE ---------------------------------- */

func (ccl *ComputeCL) RunTransmoveKernel(particleCount int, delta float32) error {
	err := clgo.SetKernelArgs(ccl.clTransmove, ccl.clPositionBuffer, ccl.clVelocityBuffer, delta)
	if err != nil {
		return err
	}

	glObjects := []clgo.CL_MEM{
		(clgo.CL_MEM)(unsafe.Pointer(ccl.clPositionBuffer)),
		(clgo.CL_MEM)(unsafe.Pointer(ccl.clVelocityBuffer)),
	}

	err = ccl.queue.AcquireGLObjects(glObjects)
	if err != nil {
		return err
	}
	err = ccl.queue.EnqueueKernel(ccl.clTransmove, particleCount)
	if err != nil {
		return err
	}

	err = ccl.queue.ReleaseGLObjects(glObjects)
	if err != nil {
		return err
	}

	ccl.queue.Finish()
	return nil
}

/* -------------------------------- GRAVITATE ------------------------------- */

func (ccl *ComputeCL) RunGravitateKernel(particleCount int) error {
	err := clgo.SetKernelArgs(ccl.clGravitateKernel, ccl.clPositionBuffer, ccl.clVelocityBuffer)
	if err != nil {
		return err
	}

	glObjects := []clgo.CL_MEM{
		(clgo.CL_MEM)(unsafe.Pointer(ccl.clPositionBuffer)),
		(clgo.CL_MEM)(unsafe.Pointer(ccl.clVelocityBuffer)),
	}

	err = ccl.queue.AcquireGLObjects(glObjects)
	if err != nil {
		return err
	}
	err = ccl.queue.EnqueueKernel(ccl.clGravitateKernel, particleCount)
	if err != nil {
		return err
	}

	err = ccl.queue.ReleaseGLObjects(glObjects)
	if err != nil {
		return err
	}

	ccl.queue.Finish()
	return nil

}

/* ---------------------------------- IDLE ---------------------------------- */

func (ccl *ComputeCL) RunIdleKernel() {}
