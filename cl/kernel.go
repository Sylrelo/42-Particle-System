package clgo

/*
#include <OpenCL/opencl.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"log"
	"os"
	"unsafe"
)

type CL_PROGRAM = C.cl_program
type CL_KERNEL = C.cl_kernel

//

type ClKernel struct {
	program CL_PROGRAM
	Kernel  CL_KERNEL

	kernels map[string]CL_KERNEL
}

func (kernel *ClKernel) Release() {
	C.clReleaseKernel(kernel.Kernel)
	C.clReleaseProgram(kernel.program)
}

// func (kernel *ClKernel) Use() {

// }

func SetKernelArgs(kernel CL_KERNEL, args ...CL_MEM) error {

	for _, arg := range args {

		// fmt.Println(arg, C.size_t(unsafe.Sizeof(arg)), unsafe.Pointer(&arg))
		errCode := C.clSetKernelArg(
			kernel,
			0,
			C.size_t(unsafe.Sizeof(arg)),
			unsafe.Pointer(&arg),
		)

		if errCode != C.CL_SUCCESS {
			log.Fatalf("Failed to set kernel argument: %d", errCode)
			return errors.New(ErrorString(int(errCode)))
		}

	}

	return nil
}

//

func CreateKernel(context CL_CONTEXT, device ClDevice, filepath string, funcName string) (ClKernel, error) {

	filebytes, err := os.ReadFile(filepath)
	if err != nil {
		return ClKernel{}, err
	}

	filecontent := string(filebytes)
	cSource := C.CString(filecontent)
	defer C.free(unsafe.Pointer(cSource))

	program := C.clCreateProgramWithSource(context, 1, &cSource, nil, nil)
	clErr := C.clBuildProgram(program, 1, &device.ID, nil, nil, nil)

	log.Printf("%s program built.", filepath)

	if clErr != C.CL_SUCCESS {
		getCompilationError(program, device)
		return ClKernel{}, errors.New(ErrorString(int(clErr)))
	}

	cFuncName := C.CString(funcName)
	defer C.free(unsafe.Pointer(cFuncName))

	kernel := C.clCreateKernel(program, cFuncName, nil)
	if clErr != C.CL_SUCCESS {
		return ClKernel{}, errors.New(ErrorString(int(clErr)))
	}
	log.Printf("%s (%s) kernel created.", filepath, funcName)

	return ClKernel{
		program: program,
		Kernel:  kernel,
		kernels: make(map[string]CL_KERNEL),
	}, nil
}

func getCompilationError(program CL_PROGRAM, device ClDevice) {
	var logSize C.size_t
	C.clGetProgramBuildInfo(program, device.ID, C.CL_PROGRAM_BUILD_LOG, 0, nil, &logSize)

	log := (*C.char)(C.malloc(logSize))
	defer C.free(unsafe.Pointer(log))

	C.clGetProgramBuildInfo(
		program,
		device.ID,
		C.CL_PROGRAM_BUILD_LOG,
		logSize,
		unsafe.Pointer(log),
		nil,
	)

	fmt.Printf("Build Log:\n%s\n", C.GoString(log))
}
