package license_plate_recognition

import (
	"bufio"
	"fmt"
	"image"
	"os"

	"github.com/LdDl/license_plate_recognition/internal/bridge"
)

// Detection represents a single object detection result.
type Detection struct {
	Rect       image.Rectangle
	ClassID    int
	Confidence float32
}

// Network wraps a single ORT model via od-bridge FFI.
type Network struct {
	handle *bridge.Handle
}

// NewNetwork creates a model from an ONNX file.
func NewNetwork(modelPath string, inputW, inputH int) (*Network, error) {
	h, err := bridge.CreateModel(modelPath, inputW, inputH)
	if err != nil {
		return nil, err
	}
	return &Network{handle: h}, nil
}

// Close frees the underlying Rust model.
func (n *Network) Close() {
	if n.handle != nil {
		n.handle.Free()
		n.handle = nil
	}
}

// Detect runs inference on an image and returns detections.
func (n *Network) Detect(img image.Image, confThreshold, nmsThreshold float32) ([]Detection, error) {
	if n.handle == nil {
		return nil, fmt.Errorf("network is closed")
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	rgb := imageToRGB(img)

	raw, err := n.handle.Detect(rgb, w, h, confThreshold, nmsThreshold)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, nil
	}

	results := make([]Detection, len(raw))
	for i, d := range raw {
		results[i] = Detection{
			Rect:       image.Rect(d.BboxX, d.BboxY, d.BboxX+d.BboxW, d.BboxY+d.BboxH),
			ClassID:    d.ClassID,
			Confidence: d.Confidence,
		}
	}
	return results, nil
}

// imageToRGB converts any image.Image to a flat RGB byte slice (HWC, row-major).
func imageToRGB(img image.Image) []byte {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	rgb := make([]byte, h*w*3)

	// Fast path for NRGBA (most common in this project)
	if nrgba, ok := img.(*image.NRGBA); ok {
		idx := 0
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			off := (y - nrgba.Rect.Min.Y) * nrgba.Stride
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				pOff := off + (x-nrgba.Rect.Min.X)*4
				rgb[idx] = nrgba.Pix[pOff]
				rgb[idx+1] = nrgba.Pix[pOff+1]
				rgb[idx+2] = nrgba.Pix[pOff+2]
				idx += 3
			}
		}
		return rgb
	}

	// Generic fallback
	idx := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rgb[idx] = uint8(r >> 8)
			rgb[idx+1] = uint8(g >> 8)
			rgb[idx+2] = uint8(b >> 8)
			idx += 3
		}
	}
	return rgb
}

// LoadClassNames reads a .names file (one class name per line).
func LoadClassNames(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var names []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	return names, scanner.Err()
}
