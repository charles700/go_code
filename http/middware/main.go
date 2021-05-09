package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// func middlewareHandler(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// 执行 handler 之前的逻辑

// 		next.ServeHTTP(w, r)

// 		// 执行 handler 之后的逻辑
// 	})
// }

// logHandler 日志输出中间件
func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Start %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Complete %s in  %v", r.URL.Path, time.Since(start))
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	html := `<doctype html>
        <html>
        <head>
          <title>Hello World</title>
        </head>
        <body>
        <p>
          <a href="/welcome">Welcome</a> |  <a href="/message">Message</a>
        </p>
        </body>
</html>`
	fmt.Fprintln(w, html)
}

func main() {
	http.Handle("/", logHandler(http.HandlerFunc(index)))
	http.ListenAndServe(":8080", nil)
}
