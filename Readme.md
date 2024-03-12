# httpbuf

[![Go Reference](https://pkg.go.dev/badge/github.com/matthewmueller/httpbuf.svg)](https://pkg.go.dev/github.com/matthewmueller/httpbuf)

Properly wrap `http.ResponseWriter`. A simple wrapper around the excellent [httpsnoop](https://github.com/felixge/httpsnoop) package that provides slightly better ergonomics.

## Install

```sh
go get github.com/matthewmueller/httpbuf
```

## Usage

```go
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httpbuf.Wrap(w)
		defer rw.Flush()
		next.ServeHTTP(rw, r)
		fmt.Println("captured", rw.Body.String())
	})
}
```

## Contributors

- Matt Mueller ([@mattmueller](https://twitter.com/mattmueller))

## License

MIT
