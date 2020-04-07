package plates_recognition_grpc

import (
	"fmt"
	"image"
	"time"
)

// YOLOResponse Ответ нейронной сети
type YOLOResponse struct {
	Plates  []PlateResponse
	Elapsed time.Duration
}

// PlateResponse Информация по детектируемой плашке
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
		result += fmt.Sprintf("Плашка #%d:\n%s\n", i, resp.Plates[i].String())
	}
	result += fmt.Sprintf("\nУшло времени на детектирование и распознавание: %v", resp.Elapsed)
	return result
}

func (presp *PlateResponse) String() string {
	result := fmt.Sprintf("\tТекст: %s\n\tОтклонение (для обнаруженных символов): %f\n\tГраницы прямоугольника: %v", presp.Text, presp.Probability, presp.Rect)
	return result
}
