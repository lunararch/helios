package camera

import "C"
import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type Camera struct {
	Position      mgl32.Vec2
	Size          mgl32.Vec2
	Zoom          float32
	Rotation      float32
	MinBounds     mgl32.Vec2
	MaxBounds     mgl32.Vec2
	BoundsEnabled bool
	target        *mgl32.Vec2
}

func New(width, height float32) *Camera {
	return &Camera{
		Position:      mgl32.Vec2{0, 0},
		Size:          mgl32.Vec2{width, height},
		Zoom:          1.0,
		Rotation:      0.0,
		BoundsEnabled: false,
	}
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	view := mgl32.Ident4()

	view = view.Mul4(mgl32.Translate3D(c.Size[0]/2.0, c.Size[1]/2.0, 0))

	view = view.Mul4(mgl32.Scale3D(c.Zoom, c.Zoom, 1.0))

	if c.Rotation != 0 {
		view = view.Mul4(mgl32.HomogRotate3DZ(c.Rotation))
	}

	view = view.Mul4(mgl32.Translate3D(-c.Position[0], -c.Position[1], 0))

	return view
}

func (c *Camera) SetBounds(minX, minY, maxX, maxY float32) {
	c.MinBounds = mgl32.Vec2{minX, minY}
	c.MaxBounds = mgl32.Vec2{maxX, maxY}
	c.BoundsEnabled = true

	c.ClampToBounds()
}

func (c *Camera) DisableBounds() {
	c.BoundsEnabled = false
}

func (c *Camera) ClampToBounds() {
	if !c.BoundsEnabled {
		return
	}

	effectiveWidth := c.Size[0] / c.Zoom
	effectiveHeight := c.Size[1] / c.Zoom
	halfWidth := effectiveWidth * 0.5
	halfHeight := effectiveHeight * 0.5

	minX := c.MinBounds[0] + halfWidth
	maxX := c.MaxBounds[0] - halfWidth

	if minX < maxX {
		if c.Position[0] < minX {
			c.Position[0] = minX
		} else if c.Position[0] > maxX {
			c.Position[0] = maxX
		}
	}

	minY := c.MinBounds[1] + halfHeight
	maxY := c.MaxBounds[1] - halfHeight
	if minY < maxY {
		if c.Position[1] < minY {
			c.Position[1] = minY
		} else if c.Position[1] > maxY {
			c.Position[1] = maxY
		}
	}
}

func (c *Camera) SetTarget(position *mgl32.Vec2) {
	c.target = position
}

func (c *Camera) ClearTarget() {
	c.target = nil
}

func (c *Camera) Update(deltaTime float32) {
	if c.target != nil {
		c.Position = *c.target
		c.ClampToBounds()
	}
}

func (c *Camera) ScreenToWorld(screenPos mgl32.Vec2) mgl32.Vec2 {
	worldX := screenPos[0]/c.Zoom + c.Position[0] - c.Size[0]/(2.0*c.Zoom)
	worldY := screenPos[1]/c.Zoom + c.Position[1] - c.Size[1]/(2.0*c.Zoom)

	if c.Rotation != 0 {
		x := worldX - c.Position[0]
		y := worldY - c.Position[1]

		cos := float32(math.Cos(float64(c.Rotation)))
		sin := float32(math.Sin(float64(c.Rotation)))
		rotatedX := x*cos - y*sin
		rotatedY := x*sin + y*cos

		worldX = rotatedX + c.Position[0]
		worldY = rotatedY + c.Position[1]
	}

	return mgl32.Vec2{worldX, worldY}
}

func (c *Camera) WorldToScreen(worldPos mgl32.Vec2) mgl32.Vec2 {
	var screenX, screenY float32

	if c.Rotation != 0 {
		x := worldPos[0] - c.Position[0]
		y := worldPos[1] - c.Position[1]

		cos := float32(math.Cos(float64(-c.Rotation)))
		sin := float32(math.Sin(float64(-c.Rotation)))
		rotatedX := x*cos - y*sin
		rotatedY := x*sin + y*cos

		screenX = (rotatedX + c.Size[0]/(2.0*c.Zoom)) * c.Zoom
		screenY = (rotatedY + c.Size[1]/(2.0*c.Zoom)) * c.Zoom
	} else {
		screenX = (worldPos[0] - c.Position[0] + c.Size[0]/(2.0*c.Zoom)) * c.Zoom
		screenY = (worldPos[1] - c.Position[1] + c.Size[1]/(2.0*c.Zoom)) * c.Zoom
	}

	return mgl32.Vec2{screenX, screenY}
}
