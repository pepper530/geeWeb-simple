package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//Default返回一个默认的路由引擎
	r := gin.Default()

	// 发起一个get请求，传入参数， url形式为：127.0.0.1：8080/user/search?username=小王子&address=沙河
	r.GET("/user/search", func(c *gin.Context) {
		username := c.DefaultQuery("username", "小王子")
		//username := c.Query("username")
		address := c.Query("address")
		//输出json结果给调用方
		c.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"username": username,
			"address":  address,
		})
	})
	r.Run(":8080")

	/*
	   ------------form表单获取参数--------------------------------
	   模拟用户登录时，传入的参数，比如需要用户名，密码，
	   url 里面只是有个  /login   点击回车后向服务器发请求，弹出登录页面，要你在对应位置输入用户名和密码，这就是完成了一次请求和一次响应
	   这里他写了html文件加载了的，在login.html文件写好内容，将login.html登录页面返回给用户（这就是这一次的响应）

	   	 r.LoadHTMLFiles("./login.html", "./index.html")
	   	 r.GET("/login", func(c *gin.Context){
	   		c.HTMl(http.StatusOK, "login.html", nil)
	   	})

	   	// 页面弹出后，输入用户名和密码，点击登录，就会触发第二次请求，
	   	但是这个登录的post请求，服务端这边无法处理，就会返回 404 not found。(即使是同一个url，使用请求方法不同，也会导致不能处理)
	   	因此，下面就写一个r.POST的请求-响应。
	   	r.POST("/login", func(c *gin.Context){
	   		username := c.PostForm("username")
	   		password := c.PostForm("password")
	   	})
	   	// 点击登录后，是触发了index.html中的form后面的action动作，发起这次请求，然后返回下面的html响应页面内容。
	   	c.HTML(http.StatusOK, "index.html", gin.H{
	   			"username": username,
	   			"password": password,
	   		})

	*/

	r.POST("/user/search", func(c *gin.Context) {
		// DefaultPostForm取不到值时会返回指定的默认值
		//username := c.DefaultPostForm("username", "小王子")
		username := c.PostForm("username")
		address := c.PostForm("address")
		//输出json结果给调用方
		c.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"username": username,
			"address":  address,
		})
	})

}
