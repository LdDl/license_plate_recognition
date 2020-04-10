## Generate protobuf *.go files
```shell
protoc -I . yolo_grpc.proto --go_out=plugins=grpc:.
```

## Download weights and configuration
```shell
cd cmd/
chmod +x download_data.sh
./download_data.sh
```

## Start server
```shell
cd cmd/server
go build -o recognition_server main.go
./recognition_server --port=50051 --platesConfig=../data/license_plates_inference.cfg --platesWeights=../data/license_plates_15000.weights --ocrConfig=../data/ocr_plates_inference.cfg --ocrWeights=../data/ocr_plates_7000.weights --saveDetected 1
```

## Test "client" - "server"
### Notice: server should be started
```shell
cd cmd/client
go build -o client_app main.go
./client_app --host=localhost --port=50051 --file=sample.jpg -x 0 -y 0 --width=4032 --height=3024
```

### Check, if server can handle error:
```shell
./client_app --host=localhost --port=50051 --file=sample.jpg -x 0 -y 0 --width=42 --height=-24
```

