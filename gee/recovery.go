package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil { // recover 捕获panic
				message := fmt.Sprintf("%s", err)    // 打印堆栈信息
				log.Printf("%s\n\n", trace(message)) // trace函数用来获取触发panic的堆栈信息
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller
	// Callers用来返回调用栈的程序计数器，第0个Caller时Callers本身，第1个是上一层的trace
	// 第2个是再上一层的defer func， 为了简洁，我们跳过了前三个caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)   // 这个函数用来获取对应的函数
		file, line := fn.FileLine(pc) // 获取到调用该函数的文件名和行号，打印再日志中。
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
