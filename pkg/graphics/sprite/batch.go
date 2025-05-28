package sprite

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/graphics/shader"
	"github.com/lunararch/helios/pkg/graphics/texture"
)

const (
	MaxBatchSize      = 1000 // Maximum number of sprites in a batch
	VertecesPerSprite = 6    // Number of vertices per sprite (6 for two triangles)
	VertexSize        = 9    // Size of each vertex (3 for position, 2 for texture coords, 4 for color = 9 floats total, 36 bytes per vertex)
)

type SpriteBatch struct {
	shader      *shader.Shader
	vao         uint32
	vbo         uint32
	vertices    []float32
	spriteCount int
	currentTex  *texture.Texture
}

func NewSpriteBatch(shaderProgram *shader.Shader) *SpriteBatch {
	maxSize := MaxBatchSize * VertecesPerSprite * VertexSize

	batch := &SpriteBatch{
		shader:      shaderProgram,
		vertices:    make([]float32, 0, maxSize),
		spriteCount: 0,
	}

	gl.GenVertexArrays(1, &batch.vao)
	gl.GenBuffers(1, &batch.vbo)

	gl.BindVertexArray(batch.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, batch.vbo)

	gl.BufferData(gl.ARRAY_BUFFER, maxSize*4, nil, gl.DYNAMIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, VertexSize*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, VertexSize*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 4, gl.FLOAT, false, VertexSize*4, gl.PtrOffset(5*4))
	gl.EnableVertexAttribArray(2)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return batch
}

func (b *SpriteBatch) Begin() {
	b.vertices = b.vertices[:0]
	b.spriteCount = 0
	b.currentTex = nil
}

func (b *SpriteBatch) End() {
	b.Flush()
}

func (b *SpriteBatch) Flush() {
	if b.spriteCount == 0 {
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(b.vertices)*4, gl.Ptr(b.vertices))

	b.shader.Use()

	if b.currentTex != nil {
		b.currentTex.Bind(0)
	}

	gl.BindVertexArray(b.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(b.vertices)/VertexSize))
	gl.BindVertexArray(0)

	b.vertices = b.vertices[:0]
	b.spriteCount = 0
}

func (b *SpriteBatch) Draw(sprite *Sprite) {
	if (b.currentTex != nil && b.currentTex != sprite.Texture) || b.spriteCount >= MaxBatchSize {
		b.Flush()
		b.currentTex = sprite.Texture
	} else if b.currentTex == nil {
		b.currentTex = sprite.Texture
	}

	model := mgl32.Translate3D(sprite.Position.X(), sprite.Position.Y(), sprite.Position.Z())

	if sprite.Rotation != 0 {
		model = model.Mul4(mgl32.Translate3D(sprite.Size.X()/2, sprite.Size.Y()/2, 0))
		model = model.Mul4(mgl32.HomogRotate3DZ(sprite.Rotation))
		model = model.Mul4(mgl32.Translate3D(-sprite.Size.X()/2, -sprite.Size.Y()/2, 0))
	}

	model = model.Mul4(mgl32.Scale3D(sprite.Size.X(), sprite.Size.Y(), 1.0))

	var texU1, texV1, texU2, texV2 float32
	if sprite.Region != nil {
		texU1, texV1 = sprite.Region.U1, sprite.Region.V1
		texU2, texV2 = sprite.Region.U2, sprite.Region.V2
	} else {
		texU1, texV1 = 0.0, 0.0
		texU2, texV2 = 1.0, 1.0
	}

	// Add the vertices to the batch (2 triangles, 6 vertices)
	// Bottom left
	pos := model.Mul4x1(mgl32.Vec4{0, 0, 0, 1})
	b.vertices = append(b.vertices,
		pos.X(), pos.Y(), pos.Z(),
		texU1, texV1,
		sprite.Color.X(), sprite.Color.Y(), sprite.Color.Z(), sprite.Color.W())

	// Bottom right
	pos = model.Mul4x1(mgl32.Vec4{1, 0, 0, 1})
	b.vertices = append(b.vertices,
		pos.X(), pos.Y(), pos.Z(),
		texU2, texV1,
		sprite.Color.X(), sprite.Color.Y(), sprite.Color.Z(), sprite.Color.W())

	// Top left
	pos = model.Mul4x1(mgl32.Vec4{0, 1, 0, 1})
	b.vertices = append(b.vertices,
		pos.X(), pos.Y(), pos.Z(),
		texU1, texV2,
		sprite.Color.X(), sprite.Color.Y(), sprite.Color.Z(), sprite.Color.W())

	// Bottom right (second triangle)
	pos = model.Mul4x1(mgl32.Vec4{1, 0, 0, 1})
	b.vertices = append(b.vertices,
		pos.X(), pos.Y(), pos.Z(),
		texU2, texV1,
		sprite.Color.X(), sprite.Color.Y(), sprite.Color.Z(), sprite.Color.W())

	// Top right
	pos = model.Mul4x1(mgl32.Vec4{1, 1, 0, 1})
	b.vertices = append(b.vertices,
		pos.X(), pos.Y(), pos.Z(),
		texU2, texV2,
		sprite.Color.X(), sprite.Color.Y(), sprite.Color.Z(), sprite.Color.W())

	// Top left (second triangle)
	pos = model.Mul4x1(mgl32.Vec4{0, 1, 0, 1})
	b.vertices = append(b.vertices,
		pos.X(), pos.Y(), pos.Z(),
		texU1, texV2,
		sprite.Color.X(), sprite.Color.Y(), sprite.Color.Z(), sprite.Color.W())

	b.spriteCount++
}

func (b *SpriteBatch) Delete() {
	gl.DeleteBuffers(1, &b.vbo)
	gl.DeleteVertexArrays(1, &b.vao)
}
