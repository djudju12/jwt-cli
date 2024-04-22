```golang
go build -o jwt-cli src/main.go
```

You can pass the JWT with pipe or parameter
```console
$ get_token | jwt-cli
$ jwt-cli "<token>"
```
