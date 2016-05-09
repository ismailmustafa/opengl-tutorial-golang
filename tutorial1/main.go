package main

import (
  "runtime"
  "github.com/go-gl/gl/v4.1-core/gl"
  "github.com/go-gl/glfw/v3.1/glfw"
  "log"
  "fmt"
)

func init() {
  runtime.LockOSThread()
}

func main() {

  // Initialize glfw
  if err := glfw.Init(); err != nil {
    log.Fatal(err)
  }
  defer glfw.Terminate()

  // Window hints
  glfw.WindowHint(glfw.Resizable, glfw.True)
  glfw.WindowHint(glfw.ContextVersionMajor, 4)
  glfw.WindowHint(glfw.ContextVersionMinor, 1)
  glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
  glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

  // Create window
  window, err := glfw.CreateWindow(640, 480, "Window", nil, nil)
  if err != nil {
    log.Fatal(err)
  }
  window.MakeContextCurrent()

  // Initialize gl
  if err := gl.Init(); err != nil {
    log.Fatal(err)
  }

  version := gl.GoStr(gl.GetString(gl.VERSION))
  fmt.Println("OpenGL version", version)

  gl.ClearColor(1.0, 1.0, 1.0, 1.0)

  for !window.ShouldClose() {

    // Maintenance
    window.SwapBuffers()
    glfw.PollEvents()
  }
}

