package main

import (
	"flag"
	"log"
	"testhttp/client"
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

	var serverAddr string
	var isClient bool

	flag.StringVar(&serverAddr, "server", "", "server address")
	flag.BoolVar(&isClient, "client", false, "run as client")
	flag.Parse()

	if isClient {
		client.StartClient(serverAddr)
	} else {
		// 打开数据库
		db, err := dbs.OpenDB("barger.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		server.StartServer(db)

	}
}
