package main

import (
	"fmt"
	"log"
	"net/http"
)

// handler echoes r.URL.Path   返回requests中的URL的path部分
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

/*
HTTP请求协议由首行、请求头、空行、正文组成。
首行=方法+URL+版本号

Request是个结构体，里面有Method, *url.URL, Header, Body等等

*/
// handler echoes r.URL.Header   返回request中的请求头部分
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}

/*  查阅了这个req.Header到底是什么，http请求头中的键值对信息
// A Header represents the key-value pairs in an HTTP header.
//
// The keys should be in canonical form, as returned by
// CanonicalHeaderKey.
type Header map[string][]string

*/

func main() {
	// 使用http中的路由映射规则 http.HandlerFunc实现了路由和规则的映射
	// 只能针对具体路由写处理逻辑
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))

	// 使用自定义路由映射规则，可以统一添加一些处理逻辑。
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))
}

// Engine is the uni handler for all requests
type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
