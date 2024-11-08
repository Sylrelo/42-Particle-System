package clgo

/*
#include <OpenCL/opencl.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"log"
	"unsafe"
)

type CL_MEM = C.cl_mem

type ClQueue struct {
	queue C.cl_command_queue
}

func CreateCommandQueue(clContext C.cl_context, device ClDevice) (ClQueue, error) {
	var clErr C.cl_int

	var queue C.cl_command_queue = C.clCreateCommandQueue(
		(C.cl_context)(unsafe.Pointer(clContext)),
		device.ID,
		0,
		&clErr,
	)

	if clErr != C.CL_SUCCESS {
		fmt.Println("Failed to create command queue")
		return ClQueue{}, errors.New(ErrorString(int(clErr)))
	}

	return ClQueue{
		queue: queue,
	}, nil
}

func (q *ClQueue) AcquireGLObjects(buffers []CL_MEM) error {
	numObjects := C.cl_uint(len(buffers))

	clErr := C.clEnqueueAcquireGLObjects(
		q.queue,
		numObjects,
		(*C.cl_mem)(unsafe.Pointer(&buffers[0])),
		0,
		nil,
		nil,
	)

	if clErr != C.CL_SUCCESS {
		return errors.New(ErrorString(int(clErr)))
	}

	return nil
}

func (q *ClQueue) ReleaseGLObjects(buffers []CL_MEM) error {
	numObjects := C.cl_uint(len(buffers))

	clErr := C.clEnqueueReleaseGLObjects(
		q.queue,
		numObjects,
		(*C.cl_mem)(unsafe.Pointer(&buffers[0])),
		0,
		nil,
		nil,
	)

	if clErr != C.CL_SUCCESS {
		return errors.New(ErrorString(int(clErr)))
	}

	return nil
}

func (q *ClQueue) EnqueueKernel(kernel CL_KERNEL, workSize int) error {
	globalWorkSize := C.size_t(workSize)

	errCode := C.clEnqueueNDRangeKernel(
		q.queue,
		kernel,
		1, nil,
		&globalWorkSize,
		nil, 0, nil, nil,
	)
	if errCode != C.CL_SUCCESS {
		log.Fatalf("Failed to enqueue kernel: %d (%s)", errCode, ErrorString(int(errCode)))
		return errors.New(ErrorString(int(errCode)))
	}

	return nil
}

func (q *ClQueue) Finish() {
	clErr := C.clFinish(q.queue)

	if clErr != C.CL_SUCCESS {
		log.Fatalln(ErrorString(int(clErr)))
	}
}
