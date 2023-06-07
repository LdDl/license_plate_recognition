# Generate *.pb.go files
```shell
protoc -I service/rpc/protos service/rpc/protos/*.proto --go_out=service/rpc/protos/ --go-grpc_out=service/rpc/protos/ --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional
```

_Note: flag `experimental_allow_proto3_optional` is not required but could be usefull in future works_

# Generate documentation
```shell
mkdir -p ./service/rpc/docs
protoc -I service/rpc/protos service/rpc/protos/*.proto --doc_out=./service/rpc/docs --doc_opt=html,service.html  --experimental_allow_proto3_optional
protoc -I service/rpc/protos service/rpc/protos/*.proto --doc_out=./service/rpc/docs --doc_opt=markdown,service.md  --experimental_allow_proto3_optional
```