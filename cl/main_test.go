package clgo

import (
	"testing"
)

func TestGetPlatforms(t *testing.T) {
	platforms, err := GetAvailablePlatforms()

	if err != nil {
		t.Fatal(err)
	}

	for i, platform := range platforms {
		t.Log(i, platform.Name, platform.Vendor, platform.Version)
	}
}

func TestGetDevices(t *testing.T) {
	platforms, _ := GetAvailablePlatforms()

	devices, err := GetAvailableDevices(platforms[0])

	if err != nil {
		t.Fatal(err)
	}

	for i, device := range devices {
		t.Log(i, device.Name, device.Vendor)
	}
}
