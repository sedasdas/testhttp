package main

import (
	"github.com/gin-gonic/gin"
	"testhttp/client"
	"testhttp/middlewares"
)

func main() {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	client.StartClient()
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务

}
