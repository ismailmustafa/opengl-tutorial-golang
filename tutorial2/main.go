package main

import (
  "runtime"
  "github.com/go-gl/gl/v4.1-core/gl"
  "github.com/go-gl/glfw/v3.1/glfw"
  "log"
  "bufio"
  "os"
  "fmt"
  "strings"
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
  window, err := glfw.CreateWindow(640, 480, "Triangle", nil, nil)
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

  // Triangle verticies
  vertexBufferData := []float32{
    -1.0, -1.0, 0.0,
     1.0, -1.0, 0.0,
     0.0,  1.0, 0.0,
  }

  // Create Vertex array object
  var vertexArrayID uint32
  gl.GenVertexArrays(1, &vertexArrayID)
  gl.BindVertexArray(vertexArrayID)

  var vertexBuffer uint32
  gl.GenBuffers(1, &vertexBuffer)
  gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
  gl.BufferData(gl.ARRAY_BUFFER, len(vertexBufferData)*4, gl.Ptr(vertexBufferData), gl.STATIC_DRAW)


  // load shaders
  programID, err := newProgram("vertexShader.vertexshader", "fragmentShader.fragmentshader")
  if err != nil {
    log.Fatal(err)
  }

  gl.ClearColor(0.11, 0.545, 0.765, 0.0)

  for !window.ShouldClose() {

    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.UseProgram(programID)

    gl.EnableVertexAttribArray(0)
    gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

    gl.DrawArrays(gl.TRIANGLES, 0, 3)

    gl.DisableVertexAttribArray(0)

    // Maintenance
    window.SwapBuffers()
    glfw.PollEvents()
  }
}

// Create a new program to run. Requires path to vertex shader and fragment
// shader files
func newProgram(vertexFilePath, fragmentFilePath string) (uint32, error) {

  // Load both shaders
  vertexShaderID, fragmentShaderID, err := loadShaders(vertexFilePath, fragmentFilePath)
  if err != nil {
    return 0, err
  }

  // Create new program
  programID := gl.CreateProgram()
  gl.AttachShader(programID, vertexShaderID)
  gl.AttachShader(programID, fragmentShaderID)
  gl.LinkProgram(programID)

  // Check status of program
  var status int32
  gl.GetProgramiv(programID, gl.LINK_STATUS, &status)
  if status == gl.FALSE {
    var logLength int32
    gl.GetProgramiv(programID, gl.INFO_LOG_LENGTH, &logLength)

    log := strings.Repeat("\x00", int(logLength+1))
    gl.GetProgramInfoLog(programID, logLength, nil, gl.Str(log))

    return 0, fmt.Errorf("failed to link program: %v", log)
  }

  // Detach shaders
  gl.DetachShader(programID, vertexShaderID)
  gl.DetachShader(programID, fragmentShaderID)

  // Delete shaders
  gl.DeleteShader(vertexShaderID)
  gl.DeleteShader(fragmentShaderID)

  return programID, nil
}

// Load both shaders and return
func loadShaders(vertexFilePath, fragmentFilePath string) (uint32, uint32, error) {

  // Compile vertex shader
  vertexShaderID, err := compileShader(readShaderCode(vertexFilePath), gl.VERTEX_SHADER)
  if err != nil {
    return 0, 0, nil
  }

  // Compile fragment shader
  fragmentShaderID, err := compileShader(readShaderCode(fragmentFilePath), gl.FRAGMENT_SHADER)
  if err != nil {
    return 0, 0, nil
  }

  return vertexShaderID, fragmentShaderID, nil
}

// Compile shader. Source is null terminated c string. shader type is self
// explanatory
func compileShader(source string, shaderType uint32) (uint32, error) {

  // Create new shader 
  shader := gl.CreateShader(shaderType)
  // Convert shader string to null terminated c string
  shaderCode, free := gl.Strs(source)
  defer free()
  gl.ShaderSource(shader, 1, shaderCode, nil)

  // Compile shader
  gl.CompileShader(shader)

  // Check shader status
  var status int32
  gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
  if status == gl.FALSE {
    var logLength int32
    gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

    log := strings.Repeat("\x00", int(logLength+1))
    gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

    return 0, fmt.Errorf("failed to compile %v: %v", source, log)
  }
  return shader, nil
}

// Read shader code from file
func readShaderCode(filePath string) string {
  code := ""
  f, err := os.Open(filePath)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  scanner := bufio.NewScanner(f)
  for scanner.Scan() {
    code += "\n" + scanner.Text()
  }
  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }
  code += "\x00"
  return code
}
