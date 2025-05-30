package entity

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	*BaseComponent
	Position mgl32.Vec3
	Rotation float32
	Scale    mgl32.Vec2
}

func NewTransform() *Transform {
	return &Transform{
		BaseComponent: NewBaseComponent(ComponentTypeTransform),
		Position:      mgl32.Vec3{0, 0, 0},
		Rotation:      0,
		Scale:         mgl32.Vec2{1, 1},
	}
}

func NewTransformWithValues(position mgl32.Vec3, rotation float32, scale mgl32.Vec2) *Transform {
	return &Transform{
		BaseComponent: NewBaseComponent(ComponentTypeTransform),
		Position:      position,
		Rotation:      rotation,
		Scale:         scale,
	}
}

func (t *Transform) SetPosition(position mgl32.Vec3) {
	t.Position = position
}

func (t *Transform) SetPosition2D(x, y float32) {
	t.Position = mgl32.Vec3{x, y, t.Position.Z()}
}

func (t *Transform) SetRotation(rotation float32) {
	t.Rotation = rotation
}

func (t *Transform) SetScale(scale mgl32.Vec2) {
	t.Scale = scale
}

func (t *Transform) SetUniformScale(scale float32) {
	t.Scale = mgl32.Vec2{scale, scale}
}

func (t *Transform) Translate(offset mgl32.Vec3) {
	t.Position = t.Position.Add(offset)
}

func (t *Transform) Translate2D(x, y float32) {
	t.Position = t.Position.Add(mgl32.Vec3{x, y, 0})
}

func (t *Transform) Rotate(angle float32) {
	t.Rotation += angle
}

func (t *Transform) GetModelMatrix() mgl32.Mat4 {
	model := mgl32.Ident4()

	model = model.Mul4(mgl32.Translate3D(t.Position.X(), t.Position.Y(), t.Position.Z()))

	if t.Rotation != 0 {
		model = model.Mul4(mgl32.HomogRotate3DZ(t.Rotation))
	}

	model = model.Mul4(mgl32.Scale3D(t.Scale.X(), t.Scale.Y(), 1.0))

	return model
}

func (t *Transform) GetWorldMatrix() mgl32.Mat4 {
	if t.entity == nil || t.entity.GetParent() == nil {
		return t.GetModelMatrix()
	}

	parentTransform := t.entity.GetParent().GetTransform()
	return parentTransform.GetWorldMatrix().Mul4(t.GetModelMatrix())
}
