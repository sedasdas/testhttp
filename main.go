package main

import (
	"encoding/json"
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
	client, _ := client.GetClientInfo()

	// 将 client 转换为 JSON 字符串
	jsonData, err := json.Marshal(client)
	if err != nil {
		log.Fatal(err)
	}

	// 将 JSON 字符串保存到数据库中
	err = dbs.SaveClient(db, client.IP, jsonData)
	if err != nil {
		log.Fatal(err)
	}

	// 从数据库中获取 client
	result, err := dbs.GetClient(db, client.IP)
	if err != nil {
		log.Fatal(err)
	}

	// 输出 client 的信息
	log.Printf("%+v\n", result)
}
