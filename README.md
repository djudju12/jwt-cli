```golang
go build -o jwt-cli src/main.go
```

You can pass the JWT with pipe or parameter

```console
$ get_token | jwt-cli
$ jwt-cli "<token>"
```

you can use the flag -t or --timestamp to parse
the timestamp seconds of the JWT to a more readable
string
```console
$ jwt-cli -t "<token>"
> {
.     "header": {
.         "alg": "HS256",
.         "typ": "JWT"
.     },
.     "payload": {
.         "exp": "2577817340 (2051-09-08 17:22:20)",
.         "iat": "1713817340 (2024-04-22 17:22:20)",
.         "iss": "example",
.         "rol": "example",
.         "sub": "example"
.     }
. }
```

