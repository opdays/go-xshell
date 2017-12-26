package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	consolePkg "github.com/yy584089958/go-xshell/console"
)

type WebServer struct {
	Engine *gin.Engine
	Debug  bool
}

func (this *WebServer) init() {
	this.Engine = gin.Default()
	if this.Debug {
		gin.SetMode(gin.DebugMode)
	}
}

func (this *WebServer) static() {
	this.Engine.LoadHTMLGlob("templates/*")
	this.Engine.Static("/static/", "static")
	this.Engine.StaticFile("/favicon.ico", "static/favicon.ico")
}

func (this *WebServer) html() {
	this.Engine.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
}

func AccessAllow(c *gin.Context) {
	//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
	//c.Writer.Header().Set("Access-Control-Allow-Headers", "content-type, x-requested-with")
	//if c.Request.Method == http.MethodOptions {
	//	c.AbortWithStatus(200)
	//	return
	//}
	c.Next()
}
func (this *WebServer) ConsoleRouter() {
	console := this.Engine.Group("/console")
	console.GET("/", consolePkg.ConsoleHtml)
	console.GET("/Ws", consolePkg.ConsoleWs)
	console.POST("/aesEn", consolePkg.AesEncrypt) //请求加密接口 会带上时间戳加密

}
func (this *WebServer) Start() {
	this.init()
	this.Engine.Use(AccessAllow)
	this.static()
	this.html()

	this.ConsoleRouter()

	this.Engine.Run("0.0.0.0:18080")

}
func init() {
	fmt.Println("init webserver")
}

var Server = &WebServer{Debug: true}

func main() {
	Server.Start()
}
