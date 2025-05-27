package shader

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"os"
	"strings"
)

type Shader struct {
	ID uint32
}

func New(vertexPath, fragmentPath string) (*Shader, error) {
	vertexCode, err := os.ReadFile(vertexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vertex shader: %w", err)
	}

	fragmentCode, err := os.ReadFile(fragmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read fragment shader: %w", err)
	}

	vertexShader, err := compileShader(string(vertexCode), gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := compileShader(string(fragmentCode), gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	defer gl.DeleteShader(fragmentShader)

	programID := gl.CreateProgram()
	gl.AttachShader(programID, vertexShader)
	gl.AttachShader(programID, fragmentShader)
	gl.LinkProgram(programID)

	var status int32
	gl.GetProgramiv(programID, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(programID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(programID, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("failed to link shader program: %v", log)
	}
	return &Shader{ID: programID}, nil
}

func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}

func (s *Shader) Delete() {
	gl.DeleteProgram(s.ID)
}

func (s *Shader) SetBool(name string, value bool) {
	var intValue int32
	if value {
		intValue = 1
	}
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), intValue)
}

func (s *Shader) SetInt(name string, value int32) {
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value)
}

func (s *Shader) SetFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value)
}

func (s *Shader) SetVec2(name string, value mgl32.Vec2) {
	gl.Uniform2fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, &value[0])
}

func (s *Shader) SetVec3(name string, value mgl32.Vec3) {
	gl.Uniform3fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, &value[0])
}

func (s *Shader) SetVec4(name string, value mgl32.Vec4) {
	gl.Uniform4fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, &value[0])
}

func (s *Shader) SetMat4(name string, value mgl32.Mat4) {
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")),
		1, false, &value[0],
	)
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile shader: %v", log)
	}

	return shader, nil
}
