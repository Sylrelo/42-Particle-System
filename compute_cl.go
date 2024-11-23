//go:build cl
// +build cl

package main

func GetCompute() Compute {
	var it Compute

	clCompute := ComputeCL{}

	it = &clCompute

	return it
}
