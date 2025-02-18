# gRPC

## Installation

### Protocol Buffer Compiler

```bash
sudo apt-get update
sudo apt-get install protobuf-compiler
```

### Protobuf and gRPC Generators

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
* Also installed with `make install`

### grpcurl

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```
* Also installed with `make install`

## Generate Protobuf Files

```bash
protoc --go_out=. --go-grpc_out=. contract/proto/user.proto
```

## Inspect gRPC Service

```bash
grpcurl -plaintext localhost:50051 list
grpcurl -import-path ./contract/proto/ -proto user.proto list
grpcurl -import-path ./contract/proto/ -proto user.proto list user.UserService
```

## Call gRPC Service

```bash
grpcurl -plaintext -import-path ./contract/proto -proto user.proto -d '{"page":1,"limit":100,"country":"GB"}' localhost:50051 user.UserService.List
grpcurl -plaintext -import-path ./contract/proto -proto user.proto -d '{"id":"28a9581d-9d49-4d3f-b0d6-7c49531e353d"}' localhost:50051 user.UserService.Get
grpcurl -plaintext -import-path ./contract/proto -proto user.proto -d '{"id":"28a9581d-9d49-4d3f-b0d6-7c49531e353d"}' localhost:50051 user.UserService.Delete
```
