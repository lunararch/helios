package texture

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

type Texture struct {
	ID     uint32
	Width  int32
	Height int32
	Format uint32
}

func LoadFromFile(filePath string) (*Texture, error) {
	filePath = filepath.Clean(filePath)

	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open texture file: %w", err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return LoadFromImage(img)
}

func LoadFromImage(img image.Image) (*Texture, error) {
	bounds := img.Bounds()
	width := int32(bounds.Dx())
	height := int32(bounds.Dy())

	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return &Texture{
		ID:     textureID,
		Width:  width,
		Height: height,
		Format: gl.RGBA,
	}, nil
}

func (t *Texture) Bind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}

func (t *Texture) Delete() {
	gl.DeleteTextures(1, &t.ID)
}
