package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	coursesURL, _ := url.Parse("http://localhost:8081")
	sessionsURL, _ := url.Parse("http://localhost:8082")

	coursesProxy := httputil.NewSingleHostReverseProxy(coursesURL)
	sessionsProxy := httputil.NewSingleHostReverseProxy(sessionsURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/courses") {
			// Is it a request for sessions of a course?
			if strings.HasSuffix(r.URL.Path, "/sessions") {
				r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
				sessionsProxy.ServeHTTP(w, r)
				return
			}
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
			coursesProxy.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/api/sessions") {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
			sessionsProxy.ServeHTTP(w, r)
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	})

	fmt.Println("API Gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
