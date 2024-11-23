package main

type Compute interface {
	InitCompute(uint32, uint32)

	RunInitKernel(int) error
	RunIdleKernel()
	RunGravitateKernel(int) error
	RunTransmoveKernel(int, float32) error
}
