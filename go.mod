module 42-particle-system

go 1.23.2

require (
	clgo v0.0.0-00010101000000-000000000000
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240506104042-037f3cc74f2a
)

require github.com/go-gl/mathgl v1.2.0 // indirect

replace clgo => ./cl
