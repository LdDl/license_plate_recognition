//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"time"

	"github.com/LdDl/license_plate_recognition/internal/bridge"
)

func main() {
	const dataDir = "../data"
	const sampleImage = "../client/sample.jpg"
	const warmup = 3
	const iters = 50

	// Load models
	platesModel, err := bridge.CreateModelCUDA(dataDir+"/license_plates.onnx", 416, 416)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer platesModel.Free()

	ocrModel, err := bridge.CreateModelCUDA(dataDir+"/ocr_plates.onnx", 416, 416)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer ocrModel.Free()

	// Load image once
	f, err := os.Open(sampleImage)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	rgb := imageToRGB(img)
	fmt.Printf("Image: %dx%d\n", w, h)

	// Plate detection
	fmt.Printf("\n=== Go→CGO→Rust: Plate detection (%dx%d) ===\n", w, h)
	for i := 0; i < warmup; i++ {
		platesModel.Detect(rgb, w, h, 0.3, 0.4)
	}
	start := time.Now()
	for i := 0; i < iters; i++ {
		platesModel.Detect(rgb, w, h, 0.3, 0.4)
	}
	elapsed := time.Since(start)
	avg := elapsed / time.Duration(iters)
	fps := float64(iters) / elapsed.Seconds()
	fmt.Printf("%d iters in %v, avg = %v/frame, %.1f FPS\n", iters, elapsed, avg, fps)

	// Get a plate crop for OCR bench
	plates, err := platesModel.Detect(rgb, w, h, 0.3, 0.4)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(plates) == 0 {
		fmt.Println("No plates found, skipping OCR bench")
		return
	}
	p := plates[0]
	x0 := max(p.BboxX, 0)
	y0 := max(p.BboxY, 0)
	x1 := min(x0+p.BboxW, w)
	y1 := min(y0+p.BboxH, h)
	cropImg := subImage(img, image.Rect(x0, y0, x1, y1))
	cropRGB := imageToRGB(cropImg)
	cw := x1 - x0
	ch := y1 - y0

	fmt.Printf("\n=== Go→CGO→Rust: OCR on crop (%dx%d) ===\n", cw, ch)
	for i := 0; i < warmup; i++ {
		ocrModel.Detect(cropRGB, cw, ch, 0.3, 0.4)
	}
	start = time.Now()
	for i := 0; i < iters; i++ {
		ocrModel.Detect(cropRGB, cw, ch, 0.3, 0.4)
	}
	elapsed = time.Since(start)
	avg = elapsed / time.Duration(iters)
	fps = float64(iters) / elapsed.Seconds()
	fmt.Printf("%d iters in %v, avg = %v/frame, %.1f FPS\n", iters, elapsed, avg, fps)

	// Full pipeline
	fmt.Printf("\n=== Go→CGO→Rust: Full pipeline (detect + OCR all) ===\n")
	for i := 0; i < warmup; i++ {
		benchFullPipeline(platesModel, ocrModel, img, rgb, w, h)
	}
	start = time.Now()
	for i := 0; i < iters; i++ {
		benchFullPipeline(platesModel, ocrModel, img, rgb, w, h)
	}
	elapsed = time.Since(start)
	avg = elapsed / time.Duration(iters)
	fps = float64(iters) / elapsed.Seconds()
	fmt.Printf("%d iters in %v, avg = %v/frame, %.1f FPS\n", iters, elapsed, avg, fps)
}

func benchFullPipeline(platesModel, ocrModel *bridge.Handle, img image.Image, rgb []byte, w, h int) {
	plates, _ := platesModel.Detect(rgb, w, h, 0.3, 0.4)
	for _, p := range plates {
		x0 := max(p.BboxX, 0)
		y0 := max(p.BboxY, 0)
		x1 := min(x0+p.BboxW, w)
		y1 := min(y0+p.BboxH, h)
		cropImg := subImage(img, image.Rect(x0, y0, x1, y1))
		cropRGB := imageToRGB(cropImg)
		cw := x1 - x0
		ch := y1 - y0
		ocrModel.Detect(cropRGB, cw, ch, 0.3, 0.4)
	}
}

func imageToRGB(img image.Image) []byte {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	rgb := make([]byte, w*h*3)
	idx := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := img.At(x, y).RGBA()
			rgb[idx] = uint8(r >> 8)
			rgb[idx+1] = uint8(g >> 8)
			rgb[idx+2] = uint8(bl >> 8)
			idx += 3
		}
	}
	return rgb
}

func subImage(img image.Image, rect image.Rectangle) image.Image {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	if si, ok := img.(subImager); ok {
		return si.SubImage(rect)
	}
	dst := image.NewNRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	for y := 0; y < rect.Dy(); y++ {
		for x := 0; x < rect.Dx(); x++ {
			dst.Set(x, y, img.At(rect.Min.X+x, rect.Min.Y+y))
		}
	}
	return dst
}
