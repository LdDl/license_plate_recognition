// +build ignore

package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"time"

	darknet "github.com/LdDl/go-darknet"
)

func main() {
	cfg := "../data/license_plates_bench.cfg"
	weights := "../data/license_plates_100000.weights"
	imgPath := "../client/sample.jpg"

	fmt.Println("Loading go-darknet (plates only)...")
	net := darknet.YOLONetwork{
		GPUDeviceIndex:           0,
		WeightsFile:              weights,
		NetworkConfigurationFile: cfg,
		Threshold:                0.3,
	}
	if err := net.Init(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer net.Close()

	// Load image
	f, err := os.Open(imgPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	img, err := jpeg.Decode(f)
	f.Close()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	bounds := img.Bounds()
	nrgba := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			nrgba.Set(x, y, img.At(x, y))
		}
	}
	fmt.Printf("Image: %dx%d\n\n", bounds.Dx(), bounds.Dy())

	const warmup = 3
	const iters = 50

	fmt.Printf("=== Go→CGO→Darknet: Plate detection ONLY ===\n")

	// Warmup
	for i := 0; i < warmup; i++ {
		dImg, _ := darknet.Image2Float32(nrgba)
		dr, _ := net.Detect(dImg)
		dImg.Close()
		fmt.Printf("  Warmup %d: net=%v, detections=%d\n", i+1, dr.NetworkOnlyTimeTaken, len(dr.Detections))
	}

	// Bench: measure everything including Image2Float32 conversion
	fmt.Printf("\nBench (including Image2Float32 conversion):\n")
	start := time.Now()
	for i := 0; i < iters; i++ {
		dImg, _ := darknet.Image2Float32(nrgba)
		dr, _ := net.Detect(dImg)
		_ = dr
		dImg.Close()
	}
	elapsed := time.Since(start)
	avg := elapsed / time.Duration(iters)
	fps := float64(iters) / elapsed.Seconds()
	fmt.Printf("%d iters in %v, avg = %v/frame, %.1f FPS\n", iters, elapsed, avg, fps)

	// Bench: measure only Detect (pre-convert image once)
	fmt.Printf("\nBench (Detect only, image pre-converted):\n")
	dImg, _ := darknet.Image2Float32(nrgba)
	defer dImg.Close()
	start = time.Now()
	for i := 0; i < iters; i++ {
		dr, _ := net.Detect(dImg)
		_ = dr
	}
	elapsed = time.Since(start)
	avg = elapsed / time.Duration(iters)
	fps = float64(iters) / elapsed.Seconds()
	fmt.Printf("%d iters in %v, avg = %v/frame, %.1f FPS\n", iters, elapsed, avg, fps)
}
