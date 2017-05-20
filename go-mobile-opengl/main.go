package main

import (
	"encoding/binary"
	"log"

	"math"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

var (
	images *glutil.Images
	fps    *debug.FPS

	program gl.Program

	buf gl.Buffer

	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform

	world gl.Uniform
	scale gl.Uniform

	posX float32
	posY float32

	an float32
	ap float32 = 5

	sc float32 = 50.0
	sp float32 = 5
)

func main() {
	app.Main(func(a app.App) {

		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				sz = e
				posX = float32(sz.WidthPx / 2)
				posY = float32(sz.HeightPx / 2)
			case paint.Event:
				if glctx == nil || e.External {
					// skip any paint events sent by the system.
					continue
				}
				onPaint(glctx, sz)
				a.Publish()
				// directly paint the next frame
				a.Send(paint.Event{})

			}
		}
	})
}

func onStart(glctx gl.Context) {
	var err error
	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	buf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	position = glctx.GetAttribLocation(program, "position")
	color = glctx.GetUniformLocation(program, "color")
	offset = glctx.GetUniformLocation(program, "offset")

	world = glctx.GetUniformLocation(program, "world")
	scale = glctx.GetUniformLocation(program, "scale")

	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)
}

func onStop(glctx gl.Context) {
	glctx.DeleteProgram(program)
	glctx.DeleteBuffer(buf)
	fps.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {

	glctx.ClearColor(0, 0, 0, 0)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	glctx.UseProgram(program)
	glctx.Uniform4f(color, 1, 1, 1, 1)

	glctx.Uniform2f(offset, posX/float32(sz.WidthPx), posY/float32(sz.HeightPx))

	if an = an + ap; an >= 360.0 {
		an = 0.0
	}
	glctx.UniformMatrix4fv(world, getrotmat(an))

	if sc = sc + sp; sc <= 20 || sc >= 195 {
		sp = sp * -1.0
	}
	glctx.UniformMatrix4fv(scale, getscalemat(sc))

	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.EnableVertexAttribArray(position)
	glctx.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)

	glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	glctx.DisableVertexAttribArray(position)

	fps.Draw(sz)
}

var triangleData = f32.Bytes(binary.LittleEndian,
	0.0, 0.5, 0.0, // top
	-0.25, -0.5, 0.0, // left
	0.25, -0.5, 0.0, // right
)

func getrotmat(a float32) []float32 {

	a64 := float64(math.Pi / 180 * a)

	cosa := float32(math.Cos(a64))
	sina := float32(math.Sin(a64))

	var rot4mat = []float32{
		cosa, -sina, 0.0, 0.0,
		sina, cosa, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}

	return rot4mat
}

func getscalemat(a float32) (mat []float32) {

	sca := a / 100.0

	mat = []float32{
		sca, 0.0, 0.0, 0.0,
		0.0, sca, 0.0, 0.0,
		0.0, 0.0, sca, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}

	return
}

const (
	coordsPerVertex = 3
	vertexCount     = 3
)

const vertexShader = `
#version 100
uniform vec2 offset;
uniform mat4 world;
uniform mat4 scale;
attribute vec4 position;

void main() {
	// offset comes in with x/y values between 0 and 1.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	gl_Position = world * position * scale + offset4 ;
}`

const fragmentShader = `
#version 100
precision lowp float; // mediump or highp not required
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`
