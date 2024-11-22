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
	"regexp"
	"unsafe"
)

type CL_PROGRAM = C.cl_program
type CL_KERNEL = C.cl_kernel

//

type ClKernel struct {
	program CL_PROGRAM
	Kernels map[string]CL_KERNEL
}

func (kernel *ClKernel) Release() {
	// C.clReleaseKernel(kernel.Kernel)
	C.clReleaseProgram(kernel.program)
}

// func (kernel *ClKernel) Use() {

// }

func SetKernelArgs(kernel CL_KERNEL, args ...interface{}) error {

	for i, arg := range args {

		var argSize C.size_t
		var argPtr unsafe.Pointer

		switch v := arg.(type) {
		case float32:
			tmp := arg.(float32)
			argSize = C.size_t(unsafe.Sizeof(tmp))
			argPtr = unsafe.Pointer(&tmp)
		case CL_MEM:
			argSize = C.size_t(unsafe.Sizeof(v))
			argPtr = unsafe.Pointer(&v)
		default:
			return fmt.Errorf("unsupported argument type: %T", v)
		}

		errCode := C.clSetKernelArg(
			kernel,
			(C.cl_uint)(i),
			argSize,
			argPtr,
			// C.size_t(unsafe.Sizeof(arg)),
			// unsafe.Pointer(&arg),
		)

		if errCode != C.CL_SUCCESS {
			log.Fatalf("Failed to set kernel argument: %d", errCode)
			return errors.New(ErrorString(int(errCode)))
		}

	}

	return nil
}

//

func InitKernels(context CL_CONTEXT, device ClDevice, filepath string) (ClKernel, error) {

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

	re := regexp.MustCompile(`__kernel\s*\w+\s*(\w+)\s*\(`)
	matches := re.FindAllStringSubmatch(filecontent, -1)

	kernels := make(map[string]CL_KERNEL)

	for _, match := range matches {
		if len(match) > 1 {
			funcName := match[1]
			createdKernel, err := createKernel(program, funcName)
			if err != nil {
				log.Fatalln(err)
			}

			kernels[funcName] = createdKernel

		}
	}

	return ClKernel{
		program: program,
		Kernels: kernels,
	}, nil
}

func createKernel(program CL_PROGRAM, fnName string) (CL_KERNEL, error) {
	var clErr C.cl_int

	cFuncName := C.CString(fnName)
	defer C.free(unsafe.Pointer(cFuncName))

	kernel := C.clCreateKernel(program, cFuncName, nil)
	if clErr != C.CL_SUCCESS {
		return nil, errors.New(ErrorString(int(clErr)))
	}
	log.Printf("%s kernel created.", fnName)

	return kernel, nil
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
