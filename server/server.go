package server

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"testhttp/client"
	dbs "testhttp/db"
	"testhttp/middlewares"
)

func StartServer(db *badger.DB) {
	// 创建 HTTP 服务器
	router := gin.Default()
	router.Use(middlewares.Cors())
	// 查询一个 client
	router.GET("/client/:ip", func(c *gin.Context) {
		ip := c.Param("ip")
		client, err := dbs.GetClient(db, ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if client == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "client not found",
			})
			return
		}
		c.JSON(http.StatusOK, client)
	})

	// 查询多个 client
	router.GET("/clients", func(c *gin.Context) {
		clients, err := dbs.GetAllClients(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, clients)
	})

	// 新增 client
	router.POST("/client", func(c *gin.Context) {
		var client client.Client
		err := c.BindJSON(&client)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		err = dbs.SaveClient(db, client.IP, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "client created",
		})
	})

	// 修改 client
	router.PUT("/client/:ip", func(c *gin.Context) {
		ip := c.Param("ip")
		var client client.Client
		err := c.BindJSON(&client)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		err = dbs.UpdateClient(db, ip, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "client updated",
		})
	})

	// 删除 client
	router.DELETE("/client/:ip", func(c *gin.Context) {
		ip := c.Param("ip")
		err := dbs.DeleteClient(db, ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "client deleted",
		})
	})

	// 启动服务器
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
