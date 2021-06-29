# License Plate Recognition with [go-darknet](https://github.com/LdDl/go-darknet) [![GoDoc](https://godoc.org/github.com/LdDl/license_plate_recognition?status.svg)](https://godoc.org/github.com/LdDl/license_plate_recognition) [![Sourcegraph](https://sourcegraph.com/github.com/LdDl/license_plate_recognition/-/badge.svg)](https://sourcegraph.com/github.com/LdDl/license_plate_recognition?badge) [![Go Report Card](https://goreportcard.com/badge/github.com/LdDl/license_plate_recognition)](https://goreportcard.com/report/github.com/LdDl/license_plate_recognition) [![GitHub tag](https://img.shields.io/github/tag/LdDl/license_plate_recognition.svg)](https://github.com/LdDl/license_plate_recognition/releases)

## Table of Contents

- [About](#about)
- [Requirements](#requirements)
- [Installation](#installation)
    - [Get source code](#get-source-code)
    - [Protobuf generation](#generate-protobuf-*.go-files-for-Go-server-and-Go-client)
    - [Neural net weights](#download-weights-and-configuration)
    - [Custom handler](#custom-handler)
- [Usage](#usage)
    - [Server](#start-server)
    - [Client](#test-client-server)


## About
This is a gRPC server which accepts image and can make license plate recognition (using YOLOv3 or YOLOv4 neural network).

First server tries to find license plate. Then it does OCR (if it's possible).

Neural networks were trained on dataset of russian license plates. But you can train it on another dataset - read about process here https://github.com/AlexeyAB/darknet#how-to-train-to-detect-your-custom-objects

Darknet architecture for finding license plates - [Yolo V3](https://arxiv.org/abs/1804.02767)

Darknet architecture for doing OCR stuff - [Yolo V4](https://arxiv.org/abs/2004.10934)

gRPC server accepts this data struct:
```protobuf
message CamInfo{
    string cam_id = 1; // id of camera (just to identify client app)
    int64 timestamp = 2; // timestamp of vehicle fixation (on client app)
    bytes image = 3; // bytes of full image in PNG-format
    Detection detection = 4; // BBox of detected vehicle (region of interest where License Plate Recognition is needed)
    VirtualLineInfo virtual_line = 5; // Line which detected object has been crossed (not necessary field, but helpfull for real-time detection on road traffic)
}
message Detection{
    int32 x_left = 1;
    int32 y_top = 2;
    int32 height = 3;
    int32 width = 4;
}
message VirtualLineInfo{
    int32 id = 1;
    int32 left_x = 2;
    int32 left_y = 3;
    int32 right_x = 4;
    int32 right_y = 5;
}
```

## Requirements
Please follow instructions from [go-darknet](https://github.com/LdDl/go-darknet#go-darknet-go-bindings-for-darknet). There you will know how to install [AlexeyAB's darknet](https://github.com/AlexeyAB/darknet) and [Go-binding](https://github.com/LdDl/go-darknet) for it.

Please follow instructions from [google/protobuff](https://github.com/golang/protobuf) for installing protobuf for Go-ecosystem.

## Instalation

### Get source code
**Notice: we are using Go-modules**
```shell
go get https://github.com/LdDl/license_plate_recognition
```

### Download weights and configuration
**Notice: please read [source code of *.sh script](cmd/download_data_RU.sh) before downloading. This script MAY NOT fit yours needs.**
```shell
cd cmd/
chmod +x download_data_RU.sh
./download_data_RU.sh
```

### Custom Handler
Do not forget (if needed) to implement [AfterFunc](https://github.com/LdDl/license_plate_recognition/blob/master/cmd/server/main.go#L93)

Example is below:
```go
....
rs := &RecognitionServer{
    ....
    AfterFunction: doSomeStuff,
}
....
func doSomeStuff(data *PlateInfo, fileContents []byte) error {
	/*
		If you want, you can implement this function by yourself (and you can wrap this function also)
		Default behaviour: do nothing.
	*/
	return nil
}
....
```

## Usage
### Start server
* Navigate to folder with server application source code
    ```shell
    cd cmd/server
    ```
* Build source code of server application to executable
    ```shell
    go build -o recognition_server main.go
    ```
* Run server application
    ```shell
    ./recognition_server --port=50051 --platesConfig=../data/license_plates_inference.cfg --platesWeights=../data/license_plates_100000.weights --ocrConfig=../data/ocr_plates_inference.cfg --ocrWeights=../data/ocr_plates_140000.weights --saveDetected 1
    ```

### Test Client-Server
**Notice: server should be started**
* Navigate to folder with server application source code
    ```shell
    cd cmd/client
    ```
* Build source code of client application to executable
    ```shell
    go build -o client_app main.go
    ```
* Run client application
    ```shell
    ./client_app --host=localhost --port=50051 --file=sample.jpg -x 0 -y 0 --width=4032 --height=3024
    ```

* Check, if server can handle error (like negative _height_ parameter):
    ```shell
    ./client_app --host=localhost --port=50051 --file=sample.jpg -x 0 -y 0 --width=42 --height=-24
    ```

* On server's side there will be output something like this:
    ```shell
    2020/06/25 15:31:57
    License plate #0:
        Text: M288HO199
        Deviation (for detected symbols): 1.808632
        Rectangle's borders: (295,1057)-(608,1204)
    License plate #1:
        Text: A100CX777
        Deviation (for detected symbols): 2.295539
        Rectangle's borders: (2049,1384)-(2582,1618)
    Elapsed to find plate and read symbols: 372.108605ms
    ```
* On server's side the directory './detected' will appear also. Detected license plates will be stored there.

