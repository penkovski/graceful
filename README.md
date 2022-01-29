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

> You can also just copy the `graceful.Start` or `graceful.StartCtx` functions 
> in your project and not depend on this repo.

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

	if err := graceful.Start(srv, 10*time.Second); err != nil {
		log.Println(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	_, _ = w.Write([]byte("hello"))
}
```

### Usage with `errgroup` and Context

```go
package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/penkovski/graceful"
	"golang.org/x/sync/errgroup"
)

func main() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(hello),
	}

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		if err := graceful.StartCtx(ctx, srv, 10*time.Second); err != nil {
			log.Println("server shutdown error:", err)
			return err
		}
		return errors.New("server stopped successfully")
	})
	if err := g.Wait(); err != nil {
		log.Println("errgroup stopped:", err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	_, _ = w.Write([]byte("hello"))
}
```
