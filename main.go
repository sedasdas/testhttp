package main

import (
	"fmt"
	"log"
	"testhttp/client"
	dbs "testhttp/db"
)

type Client struct {
	IP         string
	DiskInfo   string
	Driver     string
	CPU        string
	Card       string
	SystemInfo string
	Mem        string
}

func main() {
	// 打开数据库
	db, err := dbs.OpenDB("barger.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 获取一个 client 实例
	client, err := client.GetClientInfo()
	if err != nil {
		log.Fatal(err)
	}

	// 将 client 保存到数据库中

	err = dbs.SaveClient(db, &client)
	if err != nil {
		log.Fatal(err)
	}

	// 从数据库中获取 client
	result, err := dbs.GetClient(db, client.IP)
	if err != nil {
		log.Fatal(err)
	}

	// 输出 client 的信息
	fmt.Printf("%+v\n", result)
}
