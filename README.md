[![GoDoc](https://godoc.org/github.com/penkovski/graceful?status.svg)](https://pkg.go.dev/github.com/penkovski/graceful)
[![Go Report Card](https://goreportcard.com/badge/github.com/penkovski/graceful)](https://goreportcard.com/report/github.com/penkovski/graceful)
[![Test](https://github.com/penkovski/graceful/actions/workflows/go.yml/badge.svg)](https://github.com/penkovski/graceful/actions/workflows/go.yml)

# graceful
Graceful shutdown for HTTP servers with Go 1.8+

### What it is

A simple implementation wrapped in [~20 lines 
of code](./graceful.go). It doesn't support TLS listeners as 
I don't need it currently. 

### Installation

```shell
go get -u github.com/penkovski/graceful
```

Note that you should use Go 1.8+

> You can also just copy the `graceful.Shutdown` function in your project 
> and not depend on this repo.

### Usage

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/penkovski/graceful"
)

func main() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(hello),
	}

	if err := graceful.Shutdown(srv, 10 * time.Second); err != nil {
		log.Println(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	_, _ = w.Write([]byte("hello"))
}
```
