package clgo

/*
#include <OpenCL/opencl.h>
*/
import "C"
import (
	"errors"
	"fmt"
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

func (q *ClQueue) AcquireGLObjects(buffer CL_MEM) error {
	clErr := C.clEnqueueAcquireGLObjects(
		q.queue,
		1,
		&buffer,
		0,
		nil,
		nil,
	)

	if clErr != C.CL_SUCCESS {
		return errors.New(ErrorString(int(clErr)))
	}

	return nil
}

func (q *ClQueue) ReleaseGLObjects(buffer C.cl_mem) error {
	clErr := C.clEnqueueReleaseGLObjects(
		q.queue,
		1,
		&buffer,
		0,
		nil,
		nil,
	)

	if clErr != C.CL_SUCCESS {
		return errors.New(ErrorString(int(clErr)))
	}

	return nil
}

func (q *ClQueue) Finish() {
	C.clFinish(q.queue)
}
