package plates_recognition_grpc

import (
	"image"
	"sort"

	"github.com/LdDl/go-darknet"
)

func (net *YOLONetwork) detectSymbols(imgSrc image.Image) ([]image.Rectangle, string, float32, error) {

	scaleWidth, scaleHeight := float64(imgSrc.Bounds().Dx())/416.0, float64(imgSrc.Bounds().Dy())/416.0
	// imgResized := resize.Resize(416, 416, imgSrc, resize.Bicubic) // Должно быть для плашек в случае YOLOV3

	img, err := darknet.Image2Float32(imgSrc)
	if err != nil {
		return nil, "", 0.0, err
	}
	dr, err := net.OCR.Detect(img)
	if err != nil {
		return nil, "nil", 0.0, err
	}
	img.Close()

	var recognizedText string
	var probabilities []float32
	var rects []image.Rectangle

	sort.Sort(Detections(dr.Detections))
	// log.Println("detect symbols len", len(dr.Detections))

	for _, d := range dr.Detections {
		// log.Println("detect symbols", d)

		for i := range d.ClassIDs {

			// recognizedText += "|" + d.ClassNames[i] + "|"
			probabilities = append(probabilities, d.Probabilities[i])
			recognizedText += d.ClassNames[i]
			bBox := d.BoundingBox
			minX, minY := float64(bBox.StartPoint.X)*scaleWidth, float64(bBox.StartPoint.Y)*scaleHeight
			maxX, maxY := float64(bBox.EndPoint.X)*scaleWidth, float64(bBox.EndPoint.Y)*scaleHeight
			// minX, minY := float64(bBox.StartPoint.X), float64(bBox.StartPoint.Y)
			// maxX, maxY := float64(bBox.EndPoint.X), float64(bBox.EndPoint.Y)
			rect := image.Rect(round(minX), round(minY), round(maxX), round(maxY))
			rects = append(rects, rect)
			// gocv.Rectangle(mat, rect, color.RGBA{255, 255, 0, 0}, 2)
		}
	}
	// log.Println("OCR text", recognizedText)
	return rects, recognizedText, sqrStdDeviation(probabilities), nil
}
