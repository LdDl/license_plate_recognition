## Скачивание весов и конфигураций
```shell
cd cmd/
chmod +x download_data.sh
./download_data.sh
```

## Запуск сервера
```shell
cd cmd/server
go build -o recognition_server main.go
./recognition_server
```