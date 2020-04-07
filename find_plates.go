package plates_recognition_grpc

import (
	"image"
	"sort"
	"time"

	"github.com/LdDl/go-darknet"
	"github.com/disintegration/imaging"
)

func (net *YOLONetwork) detectPlates(imgSrc image.Image) ([]image.Rectangle, error) {

	img, err := darknet.Image2Float32(imgSrc)
	if err != nil {
		return nil, err
	}

	dr, err := net.LicensePlates.Detect(img)
	if err != nil {
		return nil, err
	}
	img.Close()

	var rects []image.Rectangle
	for _, d := range dr.Detections {
		for i := range d.ClassIDs {
			if d.ClassNames[i] != "car" && d.ClassNames[i] != "motorbike" && d.ClassNames[i] != "bus" && d.ClassNames[i] != "train" && d.ClassNames[i] != "truck" {
				// continue // Закоментировать, когда ищем плашки
			}
			bBox := d.BoundingBox
			minX, minY := float64(bBox.StartPoint.X), float64(bBox.StartPoint.Y)
			maxX, maxY := float64(bBox.EndPoint.X), float64(bBox.EndPoint.Y)
			rect := image.Rect(round(minX), round(minY), round(maxX), round(maxY))
			rects = append(rects, rect)
		}
	}

	return rects, nil
}