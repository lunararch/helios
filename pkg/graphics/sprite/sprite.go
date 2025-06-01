package sprite

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/graphics/texture"
)

type Sprite struct {
	Texture  *texture.Texture
	Region   *texture.TextureRegion // Optional texture region for sprite sheets
	Position mgl32.Vec3
	Size     mgl32.Vec2
	Rotation float32
	Color    mgl32.Vec4
}

func NewSprite(tex *texture.Texture, position mgl32.Vec3, size mgl32.Vec2) *Sprite {
	return &Sprite{
		Texture:  tex,
		Region:   nil, // No region means use entire texture
		Position: position,
		Size:     size,
		Rotation: 0,
		Color:    mgl32.Vec4{1.0, 1.0, 1.0, 1.0},
	}
}

func NewSpriteWithRegion(region *texture.TextureRegion, position mgl32.Vec3, size mgl32.Vec2) *Sprite {
	sprite := &Sprite{
		Texture:  region.Texture,
		Region:   region,
		Position: position,
		Size:     size,
		Rotation: 0,
		Color:    mgl32.Vec4{1.0, 1.0, 1.0, 1.0},
	}

	if size.X() == 0 && size.Y() == 0 {
		sprite.Size = mgl32.Vec2{float32(region.GetWidth()), float32(region.GetHeight())}
	}

	return sprite
}

func (s *Sprite) SetRegion(region *texture.TextureRegion) {
	s.Region = region
	if region != nil {
		s.Texture = region.Texture
	}
}

func (s *Sprite) GetTextureCoords() (u1, v1, u2, v2 float32) {
	if s.Region != nil {
		return s.Region.U1, s.Region.V1, s.Region.U2, s.Region.V2
	}
	return 0.0, 0.0, 1.0, 1.0
}
