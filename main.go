package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/webdav"
	"net/http"
	"strings"
)

var (
	path = flag.String("path", "./dav", "--path=/www/dav/")
	port = flag.Int("port", 80, "--port=80")
)

func main() {
	flag.Parse()
	mux := http.NewServeMux()
	handler := &webdav.Handler{
		Prefix:     "/dav/",
		FileSystem: webdav.Dir(*path),
		LockSystem: webdav.NewMemLS(),
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		if !ok {
			fmt.Println("no login")
		}
		fmt.Println("login: ", username, password)
		if strings.HasPrefix(req.RequestURI, handler.Prefix) {
			handler.ServeHTTP(w, req)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), mux); err != nil {
		fmt.Println(err)
	}
}
