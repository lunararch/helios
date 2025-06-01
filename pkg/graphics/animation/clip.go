package animation

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/graphics/texture"
)

type Frame struct {
	TextureRegion *texture.TextureRegion
	Duration      float32    // Duration in seconds
	Offset        mgl32.Vec2 // Offset for this frame (useful for sprite positioning)
}

func NewFrame(region *texture.TextureRegion, duration float32) *Frame {
	return &Frame{
		TextureRegion: region,
		Duration:      duration,
		Offset:        mgl32.Vec2{0, 0},
	}
}

func NewFrameWithOffset(region *texture.TextureRegion, duration float32, offset mgl32.Vec2) *Frame {
	return &Frame{
		TextureRegion: region,
		Duration:      duration,
		Offset:        offset,
	}
}

type AnimationClip struct {
	Name       string
	Frames     []*Frame
	Loop       bool
	TotalTime  float32
	FrameCount int
}

func NewAnimationClip(name string, loop bool) *AnimationClip {
	return &AnimationClip{
		Name:       name,
		Frames:     make([]*Frame, 0),
		Loop:       loop,
		TotalTime:  0,
		FrameCount: 0,
	}
}

func (ac *AnimationClip) AddFrame(frame *Frame) {
	ac.Frames = append(ac.Frames, frame)
	ac.TotalTime += frame.Duration
	ac.FrameCount++
}

func (ac *AnimationClip) AddFrames(frames []*Frame) {
	for _, frame := range frames {
		ac.AddFrame(frame)
	}
}

func (ac *AnimationClip) GetFrameAt(time float32) (*Frame, error) {
	if ac.FrameCount == 0 {
		return nil, fmt.Errorf("animation clip '%s' has no frames", ac.Name)
	}

	if ac.Loop && time >= ac.TotalTime {
		time = float32(int(time) % int(ac.TotalTime))
		if time < 0 {
			time += ac.TotalTime
		}
	}

	if time < 0 {
		time = 0
	} else if time >= ac.TotalTime {
		time = ac.TotalTime - 0.001 // Just before the end
	}

	currentTime := float32(0)
	for _, frame := range ac.Frames {
		if time >= currentTime && time < currentTime+frame.Duration {
			return frame, nil
		}
		currentTime += frame.Duration
	}

	return ac.Frames[ac.FrameCount-1], nil
}

func (ac *AnimationClip) GetFrameIndex(time float32) int {
	if ac.FrameCount == 0 {
		return 0
	}

	if ac.Loop && time >= ac.TotalTime {
		time = float32(int(time) % int(ac.TotalTime))
		if time < 0 {
			time += ac.TotalTime
		}
	}

	if time < 0 {
		return 0
	} else if time >= ac.TotalTime {
		return ac.FrameCount - 1
	}

	currentTime := float32(0)
	for i, frame := range ac.Frames {
		if time >= currentTime && time < currentTime+frame.Duration {
			return i
		}
		currentTime += frame.Duration
	}

	return ac.FrameCount - 1
}

func (ac *AnimationClip) IsComplete(time float32) bool {
	if ac.Loop {
		return false
	}
	return time >= ac.TotalTime
}

func (ac *AnimationClip) Reset() float32 {
	return 0
}

func (ac *AnimationClip) Clone(newName string) *AnimationClip {
	clone := NewAnimationClip(newName, ac.Loop)
	for _, frame := range ac.Frames {
		clone.AddFrame(&Frame{
			TextureRegion: frame.TextureRegion,
			Duration:      frame.Duration,
			Offset:        frame.Offset,
		})
	}
	return clone
}
