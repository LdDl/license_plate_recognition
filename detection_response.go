package plates_recognition_grpc

import (
	"fmt"
	"image"
	"time"
)

// YOLOResponse Neural net's response
type YOLOResponse struct {
	Plates  []PlateResponse
	Elapsed time.Duration
}

// PlateResponse Detected license plate information
type PlateResponse struct {
	Text          string
	Probability   float64
	Rect          image.Rectangle
	CroppedNumber *image.NRGBA
	OCRRects      []image.Rectangle
}

func (resp *YOLOResponse) String() string {
	result := ""
	for i := range resp.Plates {
		result += fmt.Sprintf("License plate #%d:\n%s\n", i, resp.Plates[i].String())
	}
	result += fmt.Sprintf("\nElapsed to find plate and read symbols: %v", resp.Elapsed)
	return result
}

func (presp *PlateResponse) String() string {
	result := fmt.Sprintf("\tText: %s\n\tDeviation (for detected symbols): %f\n\tRectangle's borders: %v", presp.Text, presp.Probability, presp.Rect)
	return result
}
