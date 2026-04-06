package license_plate_recognition

import (
	"image"
)

func (net *YOLONetwork) detectPlates(imgSrc image.Image) ([]image.Rectangle, error) {
	detections, err := net.LicensePlates.Detect(imgSrc, net.PlatesThreshold, 0.4)
	if err != nil {
		return nil, err
	}
	rects := make([]image.Rectangle, 0, len(detections))
	for _, d := range detections {
		rects = append(rects, d.Rect)
	}
	return rects, nil
}
