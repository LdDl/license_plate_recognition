package license_plate_recognition

import (
	"image"
	"sort"

	"github.com/LdDl/go-darknet"
)

// Detections slice of image.Rectangle (for sorting)
type Detections []*darknet.Detection

func (r Detections) Len() int      { return len(r) }
func (r Detections) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r Detections) Less(i, j int) bool {
	return r[i].BoundingBox.StartPoint.X < r[j].BoundingBox.StartPoint.X
}

func (net *YOLONetwork) detectSymbols(imgSrc image.Image) ([]image.Rectangle, []int, string, float32, error) {
	img, err := darknet.Image2Float32(imgSrc)
	if err != nil {
		return nil, []int{}, "", 0.0, err
	}
	dr, err := net.OCR.Detect(img)
	if err != nil {
		return nil, []int{}, "nil", 0.0, err
	}
	img.Close()

	var recognizedText string
	var classIDs []int
	var probabilities []float32
	var rects []image.Rectangle

	sort.Sort(Detections(dr.Detections))
	for _, d := range dr.Detections {
		for i := range d.ClassIDs {
			probabilities = append(probabilities, d.Probabilities[i])
			recognizedText += d.ClassNames[i]
			classIDs = append(classIDs, d.ClassIDs[i])
			bBox := d.BoundingBox
			minX, minY := float64(bBox.StartPoint.X), float64(bBox.StartPoint.Y)
			maxX, maxY := float64(bBox.EndPoint.X), float64(bBox.EndPoint.Y)
			rect := image.Rect(round(minX), round(minY), round(maxX), round(maxY))
			rects = append(rects, rect)
		}
	}
	return rects, classIDs, recognizedText, sqrStdDeviation(probabilities), nil
}
