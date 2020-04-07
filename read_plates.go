package plates_recognition_grpc

import (
	"image"
	"time"
)

// ReadLicensePlates Прогон изображения через нейронную сеть
func (net *YOLONetwork) ReadLicensePlates(imgSrc image.Image, carBox image.Rectangle) (*YOLOResponse, error) {
	resp := YOLOResponse{}
	st := time.Now()
	carimg := imaging.Crop(imgSrc, carBox)
	plates, err := net.detectPlates(carimg)
	if err != nil {
		return nil, err
	}
	for i := range plates {
		rectcropimg := imaging.Crop(carimg, plates[i])
		rects, text, prob, err := net.detectSymbols(rectcropimg)
		if err != nil {
			return nil, err
		}
		resp.Plates = append(resp.Plates, PlateResponse{
			Rect:        plates[i],
			Text:        text,
			Probability: float64(prob),
			OCRRects:    rects,
		})
	}
	resp.Elapsed = time.Since(st)
	return &resp, nil
}
