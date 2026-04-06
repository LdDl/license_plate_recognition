//go:build ignore
// +build ignore

package main

/*
#include <darknet.h>

// Optimized: NRGBA pixels → darknet CHW float32, single pass, no intermediate alloc
static void nrgba_to_darknet_image(image* im, int w, int h, const uint8_t* pix, int stride) {
    int pixel_count = w * h;
    int idx = 0;
    for (int y = 0; y < h; y++) {
        const uint8_t* row = pix + y * stride;
        for (int x = 0; x < w; x++) {
            int src = x * 4;
            im->data[(pixel_count*0) + idx] = (float)row[src]   / 255.0f;
            im->data[(pixel_count*1) + idx] = (float)row[src+1] / 255.0f;
            im->data[(pixel_count*2) + idx] = (float)row[src+2] / 255.0f;
            idx++;
        }
    }
}
*/
import "C"

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"time"
	"unsafe"

	darknet "github.com/LdDl/go-darknet"
)

// FastImage2Float32 converts *image.NRGBA directly to darknet image
// without the intermediate draw.Copy → image.RGBA allocation.
func FastImage2Float32(img *image.NRGBA) *C.image {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	cImg := C.make_image(C.int(w), C.int(h), 3)
	C.nrgba_to_darknet_image(&cImg, C.int(w), C.int(h),
		(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])), C.int(img.Stride))
	return &cImg
}

func main() {
	cfg := "../data/license_plates_bench.cfg"
	weights := "../data/license_plates_100000.weights"
	imgPath := "../client/sample.jpg"

	fmt.Println("Loading go-darknet (plates only, FIXED Image2Float32)...")
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

	// --- Original Image2Float32 ---
	fmt.Println("=== Original Image2Float32 (with draw.Copy) ===")
	for i := 0; i < warmup; i++ {
		dImg, _ := darknet.Image2Float32(nrgba)
		dr, _ := net.Detect(dImg)
		dImg.Close()
		_ = dr
	}
	start := time.Now()
	for i := 0; i < iters; i++ {
		dImg, _ := darknet.Image2Float32(nrgba)
		dr, _ := net.Detect(dImg)
		_ = dr
		dImg.Close()
	}
	elapsed := time.Since(start)
	fmt.Printf("%d iters in %v, avg = %v/frame, %.1f FPS\n\n",
		iters, elapsed, elapsed/time.Duration(iters), float64(iters)/elapsed.Seconds())

	// --- Fixed: direct NRGBA → darknet ---
	fmt.Println("=== Fixed Image2Float32 (direct NRGBA, no copy) ===")
	for i := 0; i < warmup; i++ {
		cImg := FastImage2Float32(nrgba)
		dImg := &darknet.DarknetImage{
			Width:  bounds.Dx(),
			Height: bounds.Dy(),
		}
		// We need to pass cImg to Detect, but DarknetImage.image is unexported...
		// So we measure conversion + detect separately
		C.free_image(*cImg)
		_ = dImg
	}

	// Measure conversion time alone
	fmt.Println("Conversion time only:")
	start = time.Now()
	for i := 0; i < iters; i++ {
		cImg := FastImage2Float32(nrgba)
		C.free_image(*cImg)
	}
	elapsed = time.Since(start)
	convAvg := elapsed / time.Duration(iters)
	fmt.Printf("%d iters in %v, avg = %v/conversion\n\n", iters, elapsed, convAvg)

	// Compare with original conversion time
	fmt.Println("Original conversion time only:")
	start = time.Now()
	for i := 0; i < iters; i++ {
		dImg, _ := darknet.Image2Float32(nrgba)
		dImg.Close()
	}
	elapsed = time.Since(start)
	origConvAvg := elapsed / time.Duration(iters)
	fmt.Printf("%d iters in %v, avg = %v/conversion\n", iters, elapsed, origConvAvg)
	fmt.Printf("\nSpeedup: %.1fx\n", float64(origConvAvg)/float64(convAvg))
}
