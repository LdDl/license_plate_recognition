package plates_recognition_grpc

import (
	"bytes"
	"image"
	"image/jpeg"
	"math"

	darknet "github.com/LdDl/go-darknet"
)

func round(v float64) int {
	if v >= 0 {
		return int(math.Floor(v + 0.5))
	}
	return int(math.Ceil(v - 0.5))
}

func imageToBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	return buf.Bytes(), err
}

func sqrStdDeviation(arr []float32) float32 {
	avg := averagePercent32(arr)
	aggSum := float32(0.0)
	for i := range arr {
		sqr := (arr[i] - avg) * (arr[i] - avg)
		aggSum += sqr
	}
	return float32(math.Sqrt(float64(aggSum / float32(len(arr)))))
}

func averagePercent32(arr []float32) float32 {
	sum := float32(0.0)
	for i := range arr {
		sum += arr[i]
	}
	return sum / float32(len(arr))
}

func averagePercent64(arr []float64) float64 {
	sum := 0.0
	for i := range arr {
		sum += arr[i]
	}
	return sum / float64(len(arr))
}

// Detections slice of image.Rectangle (for sorting)
type Detections []*darknet.Detection

func (r Detections) Len() int      { return len(r) }
func (r Detections) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r Detections) Less(i, j int) bool {
	return r[i].BoundingBox.StartPoint.X < r[j].BoundingBox.StartPoint.X
}
