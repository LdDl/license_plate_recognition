package license_plate_recognition

import (
	"github.com/LdDl/go-darknet"
)

const (
	gpuIndex = 0
)

// YOLONetwork Aggregate two neural networks: one is for finding license plates, another is for OCR
type YOLONetwork struct {
	LicensePlates *darknet.YOLONetwork
	OCR           *darknet.YOLONetwork
}

// NewYOLONetwork Return pointer to YOLONetwork
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
		Threshold:                0.3,
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
