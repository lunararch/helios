package texture

type TextureRegion struct {
	Texture *Texture
	U1      float32
	V1      float32
	U2      float32
	V2      float32
}

func NewTextureRegion(texture *Texture, x, y, width, height float32) *TextureRegion {
	texWidth := float32(texture.Width)
	texHeight := float32(texture.Height)

	return &TextureRegion{
		Texture: texture,
		U1:      x / texWidth,
		V1:      y / texHeight,
		U2:      (x + width) / texWidth,
		V2:      (y + height) / texHeight,
	}
}

func (tr *TextureRegion) Bind(unit uint32) {
	tr.Texture.Bind(unit)
}
