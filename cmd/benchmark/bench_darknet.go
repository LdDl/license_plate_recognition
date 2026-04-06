// +build ignore

package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"time"

	lpr "github.com/LdDl/license_plate_recognition"
)

const (
	dkWarmup = 3
	dkIters  = 50
)

func main() {
	cfg := "../data/license_plates_bench.cfg"
	weights := "../data/license_plates_100000.weights"
	ocrCfg := "../data/ocr_plates_bench.cfg"
	ocrWeights := "../data/ocr_plates_140000.weights"
	imgPath := "../client/sample.jpg"

	fmt.Println("Loading go-darknet network...")
	t0 := time.Now()
	net, err := lpr.NewYOLONetwork(cfg, weights, ocrCfg, ocrWeights, 0.3, 0.3)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("Network loaded in %v\n\n", time.Since(t0))

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

	// Convert to NRGBA (same as real pipeline)
	bounds := img.Bounds()
	nrgba := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			nrgba.Set(x, y, img.At(x, y))
		}
	}
	fmt.Printf("Image: %dx%d\n", bounds.Dx(), bounds.Dy())

	// Warmup
	fmt.Printf("\n=== Go→CGO→Darknet: Full pipeline (detect + OCR) ===\n")
	fmt.Printf("Warmup (%d runs)...\n", dkWarmup)
	for i := 0; i < dkWarmup; i++ {
		resp, err := net.ReadLicensePlates(nrgba, false)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("  Warmup %d: %v, plates=%d", i+1, resp.Elapsed, len(resp.Plates))
		for _, p := range resp.Plates {
			fmt.Printf(" [%s %.1f%%]", p.Text, p.Probability)
		}
		fmt.Println()
	}

	// Benchmark
	fmt.Printf("\nBenchmark (%d runs)...\n", dkIters)
	start := time.Now()
	for i := 0; i < dkIters; i++ {
		_, err := net.ReadLicensePlates(nrgba, false)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}
	elapsed := time.Since(start)
	avg := elapsed / time.Duration(dkIters)
	fps := float64(dkIters) / elapsed.Seconds()
	fmt.Printf("%d iters in %v, avg = %v/frame, %.1f FPS\n", dkIters, elapsed, avg, fps)
}
