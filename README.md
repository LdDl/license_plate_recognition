# WORK IN PROGRESS. DO NOT USE IT IN PRODUCTION

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
    - [Get source code](#get-source-code)
    - [Potobuf generation](#generate-protobuf-*.go-files-for-Go-server-and-Go-client)
    - [Neural net weights](#download-weights-and-configuration)
- [Usage](#usage)
    - [Server](#start-server)
    - [Client](#test-client-server)

## Requirements
Please follow instructions from [go-darknet](https://github.com/LdDl/go-darknet#go-darknet-go-bindings-for-darknet). There you will know how to install [AlexeyAB's darknet](https://github.com/AlexeyAB/darknet) and [Go-binding](https://github.com/LdDl/go-darknet) for it.

Please follow instructions from [google/protobuff](https://github.com/golang/protobuf) for installing protobuf for Go-ecosystem.

## Instalation

### Get source code
#### Notice: we are using Go-modules
```shell
go get https://github.com/LdDl/license_plate_recognition
```

### Generate protobuf *.go files for Go-server and Go-client
```shell
protoc -I . yolo_grpc.proto --go_out=plugins=grpc:.
```

### Download weights and configuration
#### Notice: please read source code of *.sh script before downloading. This script MAY NOT fit yours needs.
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
#### Notice: server should be started
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

