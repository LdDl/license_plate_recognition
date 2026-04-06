package license_plate_recognition

import (
	"image"
	"sort"
)

// detectionsByX sorts detections left-to-right by X coordinate.
type detectionsByX []Detection

func (r detectionsByX) Len() int           { return len(r) }
func (r detectionsByX) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r detectionsByX) Less(i, j int) bool { return r[i].Rect.Min.X < r[j].Rect.Min.X }

func (net *YOLONetwork) detectSymbols(imgSrc image.Image) ([]image.Rectangle, []int, string, float32, error) {
	detections, err := net.OCR.Detect(imgSrc, net.OCRThreshold, 0.4)
	if err != nil {
		return nil, nil, "", 0.0, err
	}

	sort.Sort(detectionsByX(detections))

	var recognizedText string
	var classIDs []int
	var probabilities []float32
	var rects []image.Rectangle

	for _, d := range detections {
		probabilities = append(probabilities, d.Confidence)
		classIDs = append(classIDs, d.ClassID)
		rects = append(rects, d.Rect)
		if d.ClassID < len(net.OCRClassNames) {
			recognizedText += net.OCRClassNames[d.ClassID]
		}
	}

	return rects, classIDs, recognizedText, sqrStdDeviation(probabilities), nil
}
