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
