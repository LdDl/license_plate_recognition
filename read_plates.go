package license_plate_recognition

import (
	"image"
	"time"

	"github.com/disintegration/imaging"
)

// ReadLicensePlates Returns found license plates with information about each one
func (net *YOLONetwork) ReadLicensePlates(imgSrc image.Image, saveCrop bool) (*YOLOResponse, error) {
	resp := YOLOResponse{}
	st := time.Now()
	plates, err := net.detectPlates(imgSrc)
	if err != nil {
		return nil, err
	}
	for i := range plates {
		rectcropimg := imaging.Crop(imgSrc, plates[i])
		rects, text, prob, err := net.detectSymbols(rectcropimg)
		if err != nil {
			return nil, err
		}
		plResp := PlateResponse{
			Rect:        plates[i],
			Text:        text,
			Probability: float64(prob),
			OCRRects:    rects,
		}
		if saveCrop {
			plResp.CroppedNumber = rectcropimg
		}
		resp.Plates = append(resp.Plates, plResp)
	}
	resp.Elapsed = time.Since(st)
	return &resp, nil
}
