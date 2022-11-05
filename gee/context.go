package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 为map[string]interface{}这种类型在gee里面起个别名而已，叫H
type H map[string]interface{}

// 封装http.ResponseWriter， *http.Request
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求相关信息
	Path   string     // 路径
	Method string     // 方法
	Params map[string]string
	// 响应相关信息
	StatusCode int   // 响应状态码
	// middleware
	handlers []HandlerFunc
	index    int
}

func (c *Context) Param(key string) string {
	value := c.Params[key]
	return value
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// 构造上下文，响应参数，请求参数，请求路径，请求方法
// func newContext(w http.ResponseWriter, req *http.Request) *Context {
// 	return &Context{
// 		Writer: w,
// 		Req:    req,
// 		Path:   req.URL.Path,
// 		Method: req.Method,
// 	}
// }

// 从request表单中的key获取对应的value
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 从url中获取相关参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 将状态码写入context，并且将状态码写入响应头
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 构建响应头部
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String调用我们的SetHeader和Status方法，构造string类型响应的状态码和头部，并写入响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain") // 消息格式
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

// day5 添加的  定义context的Next方法
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1, // index是记录当前执行到第几个中间件，
	}
}

// 当中间件调用Next方法时，控制权交给了下一个中间件，直到调用到最后一个中间件，
// 再从后往前，调用每个中间件在Next方法之后定义的部分。
// 回忆七米视频讲的，有点类似于串起来的栈。函数递归栈的操作。
func (c *Context) Next() {
	c.index++
	s := len(c.handlers) // []HandlerFunc的长度，计算有几个Handler
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}
