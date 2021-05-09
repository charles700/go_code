net/http 模块介绍
----

## DefaultServeMux  
```golang
func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "首页")
	})
	http.ListenAndServe(":8080", nil)
}
```
- `HandleFunc` 内部是调用了 DefaultServeMux  的 HandleFunc, 内容 使用了  DefaultServeMux 的 Handle 方法
```golang 
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, http.HandlerFunc(handler))
}
```




## http.HandlerFunc 的实现
- `http.HandlerFunc` 实现了 ServeHTTP 方法
```golang 
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```

## http.NewServeMux
- 可以通过自定义一个 实现了 ServeHTTP 方法的结构体， 来自定义路由管理器
```golang
func main() {
	mux := http.NewServeMux()

    // 使用内置的类型
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, r.URL.Path)
	}))

    // 自定义结构体
	thwelcome := &textHandler{"TextHandle!"}
	mux.Handle("/text", thwelcome)

	http.ListenAndServe(":8000", mux)

}
```

## 中间件实现
```golang 
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
	fmt.Fprintln(w, r.URL.Path)
}

func main() {
	http.Handle("/", logHandler(http.HandlerFunc(index)))
	http.ListenAndServe(":8080", nil)
}
```
