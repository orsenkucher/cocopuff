### Protobuf
Для начала нужно установить компилятор [`protoc`](https://github.com/protocolbuffers/protobuf/releases).  
И добавить содержимое папки bin в `PATH`  

Плагины для Go  
```bash
go get google.golang.org/protobuf/cmd/protoc-gen-go 
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
Убеждаемся что `protoc` найдет плагины в `PATH`
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

В `.proto` файле прописываем полный путь для импорта пакета содержащего сгенерированный код 
```protobuf
option go_package = "github.com/orsenkucher/cocopuff/pkg/api";
```

Генерим сервер и клиент
```bash
protoc -I ./api/proto \
    --go_out=./pkg/api --go_opt=paths=source_relative \
    --go-grpc_out=./pkg/api --go-grpc_opt=paths=source_relative \
    ./api/proto/adder.proto
```

Для обновления зависимостей
```bash
go mod tidy
go get -u ./... # тут аккуратно
```

Добавляем комментарий в main.go
```go
// go:generate protoc ./account.proto --go_out=plugins=grpc:./pb

protoc -I ./api/proto --go_out=plugins=grpc:pkg/api api/proto/adder.proto
```
