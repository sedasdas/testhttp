package main

import (
	"log"
	dbs "testhttp/db"
	"testhttp/server"
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
	server.StartServer(db)

}
