### 99designs/gqlgen
```bash
go get -u -d github.com/99designs/gqlgen &&
go install github.com/99designs/gqlgen
```
`tools.go`
```go
// +build tools

package tools

import _ "github.com/99designs/gqlgen"

```

Добавляем комментарий в main.go
```go
//go:generate gqlgen
```
```bash
go generate ./...
```

### vektah/dataloaden
```bash
go get -u -d github.com/vektah/dataloaden &&
go install github.com/vektah/dataloaden
```
`tools.go`
```go
// +build tools

package tools

import _ "github.com/vektah/dataloaden"

```

Добавляем комментарий в dataloader.go
```go
//go:generate dataloaden <Type>Loader string *github.com/orsenkucher/cocopuff/graphql.<Type>

//go:generate <Type>SliceLoader int []*github.com/orsenkucher/cocopuff/graphql.<Type>
```
```bash
go generate ./...
```
