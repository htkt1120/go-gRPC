# gRPC
## 概要
gRPCによるapi作成
Unary, サーバーストリーム, クライアントストリーム, 双方向ストリームのサーバー

## setup
Protocol Buffersをダウンロード、解凍し、`C:\Program Files\protoc-24.0-win64`として配置


```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

cd api
protoc --go_out=../pkg/grpc --go_opt=paths=source_relative --go-grpc_out=../pkg/grpc --go-grpc_opt=paths=source_relative hello.proto
```


## serverコマンド
```
cd cmd/server
go run main.go
```

## clientコマンド
```
cd cmd/client
go run main.go
```

## greeter testコマンド
```
cd d cmd/server/greeter
go test
```

