package texture

type TextureRegion struct {
	Texture *Texture
	U1, V1  float32 // Top-left UV coordinates
	U2, V2  float32 // Bottom-right UV coordinates
}

func NewTextureRegion(texture *Texture, u1, v1, u2, v2 float32) *TextureRegion {
	return &TextureRegion{
		Texture: texture,
		U1:      u1,
		V1:      v1,
		U2:      u2,
		V2:      v2,
	}
}

func NewTextureRegionFromPixels(texture *Texture, x, y, width, height int) *TextureRegion {
	u1 := float32(x) / float32(texture.Width)
	v1 := float32(y) / float32(texture.Height)
	u2 := float32(x+width) / float32(texture.Width)
	v2 := float32(y+height) / float32(texture.Height)

	return &TextureRegion{
		Texture: texture,
		U1:      u1,
		V1:      v1,
		U2:      u2,
		V2:      v2,
	}
}

func (tr *TextureRegion) GetWidth() int {
	return int((tr.U2 - tr.U1) * float32(tr.Texture.Width))
}

func (tr *TextureRegion) GetHeight() int {
	return int((tr.V2 - tr.V1) * float32(tr.Texture.Height))
}

func (tr *TextureRegion) GetUVs() [4]float32 {
	return [4]float32{tr.U1, tr.V1, tr.U2, tr.V2}
}

func (tr *TextureRegion) Bind(textureUnit uint32) {
	tr.Texture.Bind(textureUnit)
}
