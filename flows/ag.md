```sh
choco install ag
```
```bash
ag <text> <path>
```
### Example
When vscode "Find in files" `Ctrl`+`Shift`+`F` sucks...  
Use `ag`
```
$ ag GetReqID
graphql/authentication/authentication.go
25:                     reqID := middleware.GetReqID(ctx)

pub/log/log.go
24:                     reqID := middleware.GetReqID(ctx)
```
