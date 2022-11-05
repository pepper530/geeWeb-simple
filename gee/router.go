package gee

import (
	"net/http"
	"strings"
)

// type router struct {
// 	handlers map[string]HandlerFunc
// }

// func newRouter() *router {
// 	return &router{handlers: make(map[string]HandlerFunc)}
// }

// func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
// 	log.Printf("Route %4s - %s", method, pattern)
// 	key := method + "-" + pattern
// 	r.handlers[key] = handler
// }

// func (r *router) handle(c *Context) {
// 	key := c.Method + "-" + c.Path
// 	if handler, ok := r.handlers[key]; ok {
// 		handler(c)
// 	} else {
// 		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
// 	}
// }

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed   解析方法
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 注册路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern // 请求方式-请求路径  构成一个key
	_, ok := r.roots[method]      // roots有没有这个方法的根节点
	if !ok {                      // 没有根节点就创建一个根节点
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0) // 有根节点就插入路径
	r.handlers[key] = handler
}

// 原来是直接根据key找到handler，现在先通过getRoute过一遍路由查找树，
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

// func (r *router) handle(c *Context) {
// 	n, params := r.getRoute(c.Method, c.Path)
// 	if n != nil {
// 		c.Params = params
// 		key := c.Method + "-" + n.pattern
// 		r.handlers[key](c)
// 	} else {
// 		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
// 	}
// }

// day5 添加中间件后
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
