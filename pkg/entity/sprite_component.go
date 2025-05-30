package entity

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/graphics/sprite"
	"github.com/lunararch/helios/pkg/graphics/texture"
)

type SpriteComponent struct {
	*BaseComponent
	sprite      *sprite.Sprite
	texture     *texture.Texture
	color       mgl32.Vec4
	visible     bool
	layer       int
	spriteBatch *sprite.SpriteBatch
}

func NewSpriteComponent(tex *texture.Texture, spriteBatch *sprite.SpriteBatch) *SpriteComponent {
	return &SpriteComponent{
		BaseComponent: NewBaseComponent(ComponentTypeSprite),
		texture:       tex,
		color:         mgl32.Vec4{1.0, 1.0, 1.0, 1.0},
		visible:       true,
		layer:         0,
		spriteBatch:   spriteBatch,
	}
}

func (sc *SpriteComponent) Initialize() error {
	if err := sc.BaseComponent.Initialize(); err != nil {
		return err
	}

	if sc.entity != nil && sc.texture != nil {
		transform := sc.entity.GetTransform()
		size := mgl32.Vec2{float32(sc.texture.Width), float32(sc.texture.Height)}

		sc.sprite = sprite.NewSprite(sc.texture, transform.Position, size)
		sc.sprite.Color = sc.color
	}

	return nil
}

func (sc *SpriteComponent) Update(deltaTime float32) {
	if !sc.active || sc.sprite == nil || sc.entity == nil {
		return
	}

	transform := sc.entity.GetTransform()
	sc.sprite.Position = transform.Position
	sc.sprite.Rotation = transform.Rotation
	sc.sprite.Size = mgl32.Vec2{
		sc.sprite.Size.X() * transform.Scale.X(),
		sc.sprite.Size.Y() * transform.Scale.Y(),
	}
	sc.sprite.Color = sc.color
}

func (sc *SpriteComponent) Render(alpha float32) {
	if !sc.active || !sc.visible || sc.sprite == nil || sc.spriteBatch == nil {
		return
	}

	sc.spriteBatch.Draw(sc.sprite)
}

func (sc *SpriteComponent) SetTexture(tex *texture.Texture) {
	sc.texture = tex
	if sc.sprite != nil {
		sc.sprite.Texture = tex
		sc.sprite.Size = mgl32.Vec2{float32(tex.Width), float32(tex.Height)}
	}
}

func (sc *SpriteComponent) GetTexture() *texture.Texture {
	return sc.texture
}

func (sc *SpriteComponent) SetColor(color mgl32.Vec4) {
	sc.color = color
	if sc.sprite != nil {
		sc.sprite.Color = color
	}
}

func (sc *SpriteComponent) GetColor() mgl32.Vec4 {
	return sc.color
}

func (sc *SpriteComponent) SetVisible(visible bool) {
	sc.visible = visible
}

func (sc *SpriteComponent) IsVisible() bool {
	return sc.visible
}

func (sc *SpriteComponent) SetLayer(layer int) {
	sc.layer = layer
}

func (sc *SpriteComponent) GetLayer() int {
	return sc.layer
}

func (sc *SpriteComponent) GetSprite() *sprite.Sprite {
	return sc.sprite
}
