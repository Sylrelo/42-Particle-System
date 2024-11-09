package main

type Compute interface {
	InitCompute(uint32, uint32)
	RunInitKernel(int)
	RunGravitateKernel(int)
}
