package animation

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/graphics/texture"
)

type SpriteSheet struct {
	Texture     *texture.Texture
	FrameWidth  int32
	FrameHeight int32
	Columns     int32
	Rows        int32
	TotalFrames int32
}

func NewSpriteSheet(tex *texture.Texture, frameWidth, frameHeight int32) *SpriteSheet {
	columns := tex.Width / frameWidth
	rows := tex.Height / frameHeight
	totalFrames := columns * rows

	return &SpriteSheet{
		Texture:     tex,
		FrameWidth:  frameWidth,
		FrameHeight: frameHeight,
		Columns:     columns,
		Rows:        rows,
		TotalFrames: totalFrames,
	}
}

func (ss *SpriteSheet) GetFrameRegion(frameIndex int32) (*texture.TextureRegion, error) {
	if frameIndex < 0 || frameIndex >= ss.TotalFrames {
		return nil, fmt.Errorf("frame index %d out of range (0-%d)", frameIndex, ss.TotalFrames-1)
	}

	col := frameIndex % ss.Columns
	row := frameIndex / ss.Columns

	u1 := float32(col*ss.FrameWidth) / float32(ss.Texture.Width)
	v1 := float32(row*ss.FrameHeight) / float32(ss.Texture.Height)
	u2 := float32((col+1)*ss.FrameWidth) / float32(ss.Texture.Width)
	v2 := float32((row+1)*ss.FrameHeight) / float32(ss.Texture.Height)

	return &texture.TextureRegion{
		Texture: ss.Texture,
		U1:      u1,
		V1:      v1,
		U2:      u2,
		V2:      v2,
	}, nil
}

func (ss *SpriteSheet) GetFrameRegions(startFrame, endFrame int32) ([]*texture.TextureRegion, error) {
	if startFrame < 0 || endFrame >= ss.TotalFrames || startFrame > endFrame {
		return nil, fmt.Errorf("invalid frame range %d-%d", startFrame, endFrame)
	}

	regions := make([]*texture.TextureRegion, 0, endFrame-startFrame+1)
	for i := startFrame; i <= endFrame; i++ {
		region, err := ss.GetFrameRegion(i)
		if err != nil {
			return nil, err
		}
		regions = append(regions, region)
	}

	return regions, nil
}

func (ss *SpriteSheet) CreateAnimation(name string, startFrame, endFrame int32, frameDuration float32, loop bool) (*AnimationClip, error) {
	regions, err := ss.GetFrameRegions(startFrame, endFrame)
	if err != nil {
		return nil, err
	}

	clip := NewAnimationClip(name, loop)
	for _, region := range regions {
		frame := NewFrame(region, frameDuration)
		clip.AddFrame(frame)
	}

	return clip, nil
}

func (ss *SpriteSheet) CreateAnimationWithDurations(name string, frameIndices []int32, durations []float32, loop bool) (*AnimationClip, error) {
	if len(frameIndices) != len(durations) {
		return nil, fmt.Errorf("frame indices and durations arrays must have the same length")
	}

	clip := NewAnimationClip(name, loop)
	for i, frameIndex := range frameIndices {
		region, err := ss.GetFrameRegion(frameIndex)
		if err != nil {
			return nil, err
		}

		frame := NewFrame(region, durations[i])
		clip.AddFrame(frame)
	}

	return clip, nil
}

type AnimationBuilder struct {
	spriteSheet *SpriteSheet
	clip        *AnimationClip
}

func NewAnimationBuilder(spriteSheet *SpriteSheet, name string, loop bool) *AnimationBuilder {
	return &AnimationBuilder{
		spriteSheet: spriteSheet,
		clip:        NewAnimationClip(name, loop),
	}
}

func (ab *AnimationBuilder) AddFrame(frameIndex int32, duration float32) *AnimationBuilder {
	region, err := ab.spriteSheet.GetFrameRegion(frameIndex)
	if err != nil {
		return ab // Silently ignore invalid frames
	}

	frame := NewFrame(region, duration)
	ab.clip.AddFrame(frame)
	return ab
}

func (ab *AnimationBuilder) AddFrameWithOffset(frameIndex int32, duration float32, offset mgl32.Vec2) *AnimationBuilder {
	region, err := ab.spriteSheet.GetFrameRegion(frameIndex)
	if err != nil {
		return ab // Silently ignore invalid frames
	}

	frame := NewFrameWithOffset(region, duration, offset)
	ab.clip.AddFrame(frame)
	return ab
}

func (ab *AnimationBuilder) AddFrameRange(startFrame, endFrame int32, duration float32) *AnimationBuilder {
	for i := startFrame; i <= endFrame; i++ {
		ab.AddFrame(i, duration)
	}
	return ab
}

func (ab *AnimationBuilder) AddFrames(frames []int32, durations []float32) *AnimationBuilder {
	if len(frames) != len(durations) {
		return ab // Silently ignore mismatched arrays
	}

	for i, frameIndex := range frames {
		ab.AddFrame(frameIndex, durations[i])
	}
	return ab
}

func (ab *AnimationBuilder) Build() *AnimationClip {
	return ab.clip
}
