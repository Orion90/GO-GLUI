// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Originally put together by github.com/segfault88, but
// I thought it might be useful to somebody else too.

// It took me quite a lot of frustration and messing around
// to get a basic example of glfw3 with modern OpenGL (3.3)
// with shaders etc. working. Hopefully this will save you
// some trouble. Enjoy!

package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glu"
	"runtime"
	"time"
)

const (
	vertex = `#version 330

in vec2 position;
in vec4 color;

out vec4 Color;

void main()
{
    Color = color;
    gl_Position = vec4(position, 0.0, 1.0);
}`

	fragment = `#version 330

in vec4 Color;

out vec4 outColor;

void main()
{
    outColor = vec4(Color);
}`
)

func main() {
	// lock glfw/gl calls to a single thread
	runtime.LockOSThread()

	glfw.Init()
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	glfw.WindowHint(glfw.DepthBits, 24)
	window, err := glfw.CreateWindow(800, 600, "Example", nil, nil)
	if err != nil {
		panic(err)
	}

	defer window.Destroy()

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Init()

	vao := gl.GenVertexArray()
	vao.Bind()

	vbo := gl.GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)

	verticies := []float32{
		0.0, 0.5, 1.0, 0.0, 0.0, 0.1, // Vertex 1: Red
		0.5, 0, 0.0, 1.0, 0.0, 0.1, // Vertex 2: Green
		-0.5, 0, 0.0, 0.0, 1.0, 0.1, // Vertex 3: Blue

		0.5, 0, 0.0, 1.0, 0.0, 0.1, // Vertex 2: Green
		-0.5, 0, 0.0, 0.0, 1.0, 0.1, // Vertex 3: Blue
		0, -0.5, 1.0, 0.0, 0.0, 0.1, // Vertex 1: Red

		-0.5, -0.5, 0.8, 0.8, 0.8, 0.1, // Vertex 1: Red
		0.5, -0.5, 0.8, 0.8, 0.8, 1, // Vertex 1: Red
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(verticies)*4, verticies, gl.STATIC_DRAW)

	vertex_shader := gl.CreateShader(gl.VERTEX_SHADER)
	vertex_shader.Source(vertex)
	vertex_shader.Compile()
	fmt.Println(vertex_shader.GetInfoLog())
	defer vertex_shader.Delete()

	fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragment_shader.Source(fragment)
	fragment_shader.Compile()
	fmt.Println(fragment_shader.GetInfoLog())
	defer fragment_shader.Delete()

	program := gl.CreateProgram()
	program.AttachShader(vertex_shader)
	program.AttachShader(fragment_shader)

	program.BindFragDataLocation(0, "outColor")
	program.Link()
	program.Use()
	defer program.Delete()

	positionAttrib := program.GetAttribLocation("position")
	positionAttrib.AttribPointer(2, gl.FLOAT, false, 24, nil)
	positionAttrib.EnableArray()
	defer positionAttrib.DisableArray()
	var col uintptr
	col = 8
	colorAttrib := program.GetAttribLocation("color")
	colorAttrib.AttribPointer(4, gl.FLOAT, false, 24, col)
	colorAttrib.EnableArray()
	defer colorAttrib.DisableArray()
	// uniColor := program.GetUniformLocation("triangleColor")
	for !window.ShouldClose() {
		gl.ClearColor(0, 0, 0, 0.2)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		gl.DrawArrays(gl.LINES, 6, 2)
		window.SwapBuffers()
		glfw.PollEvents()

		if window.GetMouseButton(glfw.MouseButton1) == glfw.Press {
			fmt.Println(window.GetCursorPosition())
			time.Sleep(150 * time.Millisecond)
		}
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}

func checkGLerror() {
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		string, _ := glu.ErrorString(glerr)
		panic(string)
	}
}
