package license_plate_recognition

const (
	gpuIndex = 0
)

// YOLONetwork Aggregate two neural networks: one is for finding license plates, another is for OCR
type YOLONetwork struct {
	LicensePlates *Network
	OCR           *Network
	// Class names loaded from .names files
	PlatesClassNames []string
	OCRClassNames    []string
	// Thresholds
	PlatesThreshold float32
	OCRThreshold    float32
}

// NewYOLONetwork creates a YOLONetwork from ONNX models.
func NewYOLONetwork(
	platesModel string, platesNames string,
	ocrModel string, ocrNames string,
	inputW, inputH int,
	platesThreshold, ocrThreshold float32,
) (*YOLONetwork, error) {
	plates, err := NewNetwork(platesModel, inputW, inputH)
	if err != nil {
		return nil, err
	}
	ocr, err := NewNetwork(ocrModel, inputW, inputH)
	if err != nil {
		plates.Close()
		return nil, err
	}

	platesClassNames, err := LoadClassNames(platesNames)
	if err != nil {
		plates.Close()
		ocr.Close()
		return nil, err
	}

	ocrClassNames, err := LoadClassNames(ocrNames)
	if err != nil {
		plates.Close()
		ocr.Close()
		return nil, err
	}

	return &YOLONetwork{
		LicensePlates:    plates,
		OCR:              ocr,
		PlatesClassNames: platesClassNames,
		OCRClassNames:    ocrClassNames,
		PlatesThreshold:  platesThreshold,
		OCRThreshold:     ocrThreshold,
	}, nil
}

// Close frees both underlying models.
func (net *YOLONetwork) Close() {
	if net.LicensePlates != nil {
		net.LicensePlates.Close()
	}
	if net.OCR != nil {
		net.OCR.Close()
	}
}
