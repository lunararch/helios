package sprite

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/graphics/texture"
)

type Sprite struct {
	Position mgl32.Vec3
	Size     mgl32.Vec2
	Color    mgl32.Vec4
	Texture  *texture.Texture
	Region   *texture.TextureRegion
	Rotation float32
}

func NewSprite(tex *texture.Texture, position mgl32.Vec3, size mgl32.Vec2) *Sprite {
	return &Sprite{
		Position: position,
		Size:     size,
		Color:    mgl32.Vec4{1.0, 1.0, 1.0, 1.0},
		Texture:  tex,
		Rotation: 0.0,
	}
}

func (s *Sprite) WithRegion(region *texture.TextureRegion) *Sprite {
	s.Region = region
	return s
}

func (s *Sprite) WithColor(color mgl32.Vec4) *Sprite {
	s.Color = color
	return s
}
