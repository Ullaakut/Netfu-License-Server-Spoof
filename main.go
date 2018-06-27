package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func makeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("X-IPLB-Instance", "17342")
		w.Header().Set("X-Powered-By", "PHP/5.4.45")
		w.Header().Set("Accept", "text/html")
		w.Header().Set("Accept-Encoding", "gzip, deflate")
		w.Header().Set("Accept-Language", "fr,fr-FR")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Host", "netfu.net")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func validLicense(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "1")
	fmt.Printf("Received request from %s\n", r.RemoteAddr)
}

func main() {
	http.HandleFunc("/", makeGzipHandler(validLicense))
	http.ListenAndServe(":80", nil)
	fmt.Println("server running at http://localhost:80")

	for {
		time.Sleep(10 * time.Second)
	}
}
