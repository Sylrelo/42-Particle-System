package main

import (
	"image"
	"log"
	"os"
	"path"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	_ "image/jpeg"
)

type GridViewUniforms struct {
	perspective int32
	camera      int32
	model       int32
}

type GridViewProgram struct {
	program uint32

	uniforms GridViewUniforms

	vbo uint32
	vao uint32
	tex uint32
}

func (gvp *GridViewProgram) Init() {
	currentPath, _ := os.Getwd()

	renderFragmentShader, err := CompileShader(currentPath+"/resources/shaders/gridview.frag", FRAGMENT)
	ExitOnError(err)

	renderVertexShader, err := CompileShader(currentPath+"/resources/shaders/gridview.vert", VERTEX)
	ExitOnError(err)

	gvp.program, err = CreateProgram(renderFragmentShader, renderVertexShader)
	ExitOnError(err)

	//

	gvp.LoadTexture()

	gvp.uniforms.perspective = gl.GetUniformLocation(gvp.program, gl.Str("perspective\x00"))
	gvp.uniforms.camera = gl.GetUniformLocation(gvp.program, gl.Str("camera\x00"))
	gvp.uniforms.model = gl.GetUniformLocation(gvp.program, gl.Str("model\x00"))

	// Generate Buffers and VAO
	gl.GenVertexArrays(1, &gvp.vao)
	gl.GenBuffers(1, &gvp.vbo)
	// fmt.Println(gl.GetError())

	gl.BindVertexArray(gvp.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, gvp.vbo)

	gl.BufferData(gl.ARRAY_BUFFER, len(QUAD_VERTICES)*4, gl.Ptr(QUAD_VERTICES), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (gvp *GridViewProgram) LoadTexture() {
	currentPath, _ := os.Getwd()

	file, err := os.Open(path.Join(currentPath, "resources", "grid.jpg"))
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalln(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	gl.GenTextures(1, &gvp.tex)
	gl.BindTexture(gl.TEXTURE_2D, gvp.tex)

	// Texture parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)

	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (gvp *GridViewProgram) Render(orbitCamera *OrbitCamera) {
	gl.UseProgram(gvp.program)
	gl.BindVertexArray(gvp.vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, gvp.tex)

	model := mgl32.Scale3D(10000, 10000, 10000)

	gvp.BindUniforms(orbitCamera, model)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	rotMatrix := mgl32.HomogRotate3DX(1.571)

	model = model.Mul4(rotMatrix)

	gvp.BindUniforms(orbitCamera, model)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (gvp *GridViewProgram) BindUniforms(orbitCamera *OrbitCamera, model mgl32.Mat4) {
	gl.UniformMatrix4fv(gvp.uniforms.camera, 1, false, &orbitCamera.cameraMatrix[0])
	gl.UniformMatrix4fv(gvp.uniforms.perspective, 1, false, &orbitCamera.perspectiveMatrix[0])
	gl.UniformMatrix4fv(gvp.uniforms.model, 1, false, &model[0])
}
