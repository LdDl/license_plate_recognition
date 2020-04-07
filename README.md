## Скачивание весов и конфигураций
```shell
cd cmd/
chmod +x download_data.sh
./download_data.sh
```

## Генерация protobuff файлов расширения *.go
```shell
protoc -I . yolo_grpc.proto --go_out=plugins=grpc:.
```

## Запуск сервера
```shell
cd cmd/server
go build -o recognition_server main.go
./recognition_server --port=50051 --platesConfig=../data/license_plates_inference.cfg --platesWeights=../data/license_plates_15000.weights --ocrConfig=../data/ocr_plates_inference.cfg --ocrWeights=../data/ocr_plates_7000.weights
```