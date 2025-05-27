package sprite

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/graphics/shader"
)

type Renderer struct {
	shader *shader.Shader
	vao    uint32
	vbo    uint32
}

func NewRenderer(shaderProgram *shader.Shader) *Renderer {
	renderer := &Renderer{
		shader: shaderProgram,
	}

	vertices := []float32{
		0.0, 0.0, 0.0, 0.0, 0.0,
		1.0, 0.0, 0.0, 1.0, 0.0,
		0.0, 1.0, 0.0, 0.0, 1.0,

		1.0, 0.0, 0.0, 1.0, 0.0,
		1.0, 1.0, 0.0, 1.0, 1.0,
		0.0, 1.0, 0.0, 0.0, 1.0,
	}

	gl.GenVertexArrays(1, &renderer.vao)
	gl.GenBuffers(1, &renderer.vbo)

	gl.BindVertexArray(renderer.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return renderer
}

func (r *Renderer) DrawSprite(sprite *Sprite) {
	r.shader.Use()

	model := mgl32.Translate3D(sprite.Position.X(), sprite.Position.Y(), sprite.Position.Z())

	if sprite.Rotation != 0 {
		model = model.Mul4(mgl32.Translate3D(sprite.Size.X()/2, sprite.Size.Y()/2, 0))
		model = model.Mul4(mgl32.HomogRotate3DZ(sprite.Rotation))
		model = model.Mul4(mgl32.Translate3D(-sprite.Size.X()/2, -sprite.Size.Y()/2, 0))
	}

	model = model.Mul4(mgl32.Scale3D(sprite.Size.X(), sprite.Size.Y(), 1.0))

	r.shader.SetMat4("model", model)
	r.shader.SetVec4("color", sprite.Color)

	if sprite.Region != nil {
		gl.BindVertexArray(r.vbo)

		texCoords := []float32{
			// Bottom left, bottom right, top left
			sprite.Region.U1, sprite.Region.V1,
			sprite.Region.U2, sprite.Region.V1,
			sprite.Region.U1, sprite.Region.V2,

			// Bottom right, top right, top left
			sprite.Region.U2, sprite.Region.V1,
			sprite.Region.U2, sprite.Region.V2,
			sprite.Region.U1, sprite.Region.V2,
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
		gl.BufferSubData(gl.ARRAY_BUFFER, 3*4, len(texCoords), gl.Ptr(texCoords))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		sprite.Region.Bind(0)
	} else {
		sprite.Texture.Bind(0)
	}

	gl.BindVertexArray(r.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)
}

func (r *Renderer) Delete() {
	gl.DeleteVertexArrays(1, &r.vao)
	gl.DeleteBuffers(1, &r.vbo)
}
