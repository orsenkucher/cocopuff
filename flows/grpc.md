### Protobuf
Для начала нужно установить компилятор [`protoc`](https://github.com/protocolbuffers/protobuf/releases).  
И добавить содержимое папки bin в `PATH`  

Плагины для Go  
```bash
go get -u -d google.golang.org/protobuf/cmd/protoc-gen-go &&
go install google.golang.org/protobuf/cmd/protoc-gen-go &&

go get -u -d google.golang.org/grpc/cmd/protoc-gen-go-grpc &&
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
`tools.go`
```go
// +build tools

package tools

import (
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

```

Убеждаемся что `protoc` найдет плагины в `PATH`
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

В `.proto` файле прописываем полный путь для импорта пакета содержащего сгенерированный код 
```protobuf
option go_package = "github.com/orsenkucher/cocopuff/<service>/pb";
```

Генерим сервер и клиент
```bash
mkdir ./pb
protoc -I ../api/proto \
    --go_out=./pb --go_opt=paths=source_relative \
    --go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
    ../api/proto/account.proto
```

Для обновления зависимостей
```bash
go mod tidy
go get -u -d ./... # тут аккуратно
```

Добавляем комментарий в main.go
```go
//go:generate mkdir ./pb -p
//go:generate protoc ../api/proto/account.proto -I ../api/proto --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative
```
```bash
go generate ./...
```
