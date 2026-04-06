package license_plate_recognition

type Configuration struct {
	ServerConf serverInstanceConfiguration `toml:"server"`
	YOLOPlates yoloConfiguration           `toml:"yolo_plates"`
	YOLOOCR    yoloConfiguration           `toml:"yolo_ocr"`
}

type serverInstanceConfiguration struct {
	Host         string `toml:"host"`
	Port         int32  `toml:"port"`
	SaveDetected bool   `toml:"save_detected"`
	QueueLimit   int    `toml:"queue_limit"`
}

type yoloConfiguration struct {
	Model     string  `toml:"model"`
	Names     string  `toml:"names"`
	InputW    int     `toml:"input_w"`
	InputH    int     `toml:"input_h"`
	Threshold float32 `toml:"threshold"`
}
