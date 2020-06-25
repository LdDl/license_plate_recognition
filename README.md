## Table of Contents

- [About](#about)
- [Requirements](#requirements)
- [Installation](#installation)
    - [Get source code](#get-source-code)
    - [Protobuf generation](#generate-protobuf-*.go-files-for-Go-server-and-Go-client)
    - [Neural net weights](#download-weights-and-configuration)
- [Usage](#usage)
    - [Server](#start-server)
    - [Client](#test-client-server)


## About
This is a gRPC server which accepts image and can make license plate recognition (using YOLOv3 or YOLOv4 neural network).

First server tries to find license plate. Then it does OCR (if it's possible).

Neural networks were trained on dataset of russian license plates. But you can train it on another dataset - read about process here https://github.com/AlexeyAB/darknet#how-to-train-to-detect-your-custom-objects


gRPC server accepts this data struct:
```protobuf
message CamInfo{
    string cam_id = 1; // id of camera (just to identify client app)
    int64 timestamp = 2; // timestamp of vehicle fixation (on client app)
    bytes image = 3; // bytes of full image in PNG-format
    Detection detection = 4; // BBox of detected vehicle (region of interest where License Plate Recognition is needed)
}
message Detection{
    int32 x_left = 1;
    int32 y_top = 2;
    int32 height = 3;
    int32 width = 4;
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

### Generate protobuf *.go files for Go-server and Go-client
```shell
protoc -I . yolo_grpc.proto --go_out=plugins=grpc:.
```

### Download weights and configuration
**Notice: please read [source code of *.sh script](cmd/download_data_RU.sh) before downloading. This script MAY NOT fit yours needs.**
```shell
cd cmd/
chmod +x download_data_RU.sh
./download_data_RU.sh
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
    ./recognition_server --port=50051 --platesConfig=../data/license_plates_inference.cfg --platesWeights=../data/license_plates_15000.weights --ocrConfig=../data/ocr_plates_inference.cfg --ocrWeights=../data/ocr_plates_7000.weights --saveDetected 1
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

