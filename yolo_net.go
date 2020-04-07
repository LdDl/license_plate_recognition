package plates_recognition_grpc

import (
	"github.com/LdDl/go-darknet"
)

const (
	gpuIndex = 0
)

// YOLONetwork Объединение двух сетей: для поиска плашек и для распознавания букв/цифр
type YOLONetwork struct {
	LicensePlates *darknet.YOLONetwork
	OCR           *darknet.YOLONetwork
}

// NewYOLONetwork Возвращет новую нейронную сеть для поиска номеров и букв на них
func NewYOLONetwork(platesCfg, platesWeights, ocrCfg, ocrWeights string) (*YOLONetwork, error) {
	plates := darknet.YOLONetwork{
		GPUDeviceIndex:           0,
		WeightsFile:              platesWeights,
		NetworkConfigurationFile: platesCfg,
		Threshold:                0.3,
	}

	ocr := darknet.YOLONetwork{
		GPUDeviceIndex:           0,
		WeightsFile:              ocrWeights,
		NetworkConfigurationFile: ocrCfg,
		Threshold:                0.4,
	}

	if err := plates.Init(); err != nil {
		return nil, err
	}
	if err := ocr.Init(); err != nil {
		return nil, err
	}

	return &YOLONetwork{
		LicensePlates: &plates,
		OCR:           &ocr,
	}, nil
}
