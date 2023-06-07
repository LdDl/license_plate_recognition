package license_plate_recognition

import (
	"bytes"
	"image"
	"image/jpeg"
	"math"
	"os"
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

func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, 0777)
	if err == nil || os.IsExist(err) {
		return nil
	}
	return err
}

func pascalVOC2YOLO(x1, y1, x2, y2, imgW, imgH float64) [4]float64 {
	return [4]float64{
		((x2 + x1) / (2 * imgW)),
		((y2 + y1) / (2 * imgH)),
		(x2 - x1) / imgW,
		(y2 - y1) / imgH,
	}
}
