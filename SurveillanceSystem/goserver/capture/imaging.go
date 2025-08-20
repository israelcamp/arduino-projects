package capture

import (
	"bytes"
	"image"
	jpeg "image/jpeg"
	"github.com/disintegration/imaging"
)

func resize(img image.Image, w, h int) *image.NRGBA {
	// If w or h is 0 imaging keeps aspect-ratio automatically.
	return imaging.Resize(img, w, h, imaging.Lanczos)
}

func decode(imgBytes []byte) (image.Image, string, error) {
	return image.Decode(bytes.NewReader(imgBytes))
}

func encodeJPEG(img image.Image, quality int) ([]byte, error) {
	buf := new(bytes.Buffer)
	opts := &jpeg.Options{Quality: quality} // 80-85 is a good default
	if err := jpeg.Encode(buf, img, opts); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ResizeImageBytes resizes any JPEG/PNG/WebP image in memory.
// Pass height or width = 0 to resize proportionally.
func ResizeImageBytes(in []byte, width, height, quality int) ([]byte, error) {
	src, _, err := decode(in)            // step 1
	if err != nil { return nil, err }

	dst := resize(src, width, height) // step 2 (choose your lib)

	return encodeJPEG(dst, quality)   // step 3
}
