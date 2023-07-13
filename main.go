package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/dgraph-io/badger/v3"
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
	db, err := badger.Open(badger.DefaultOptions("barger.db"))
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// 获取一个 client 实例
	client := getClientInfo()

	// 将 client 保存到数据库中
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(client.IP), encode(&client))
	})
	if err != nil {
		log.Println(err)
	}

	// 从数据库中获取 client
	var result Client
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(client.IP))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			result = decode(val)
			return nil
		})
		return err
	})
	if err != nil {
		log.Println(err)
	}

	// 输出 client 的信息
	fmt.Printf("%+v\n", result)
}

// 将 Client 结构体转换为字节数组
func encode(client *Client) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// 将 client 编码为字节序列
	err := enc.Encode(client)
	if err != nil {
		log.Println(err)
	}

	return buf.Bytes()
}

// 将字节数组转换为 Client 结构体
func decode(data []byte) Client {
	var client Client

	// 将字节序列解码为 client 结构体
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&client)
	if err != nil {
		log.Println(err)
	}

	return client
}

// 获取 client 信息
func getClientInfo() Client {
	client := Client{}

	// 执行 ip 命令
	cmd := exec.Command("bash", "-c", "ip a | grep inet | grep -v inet6 | awk -F 'inet ' '{print $2}' | awk -F '/' '{print $1}' | grep 10")
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	client.IP = strings.TrimSpace(string(out))

	// 执行 df 命令
	cmd = exec.Command("bash", "-c", "df -hl")
	out, err = cmd.Output()
	if err != nil {
		log.Println(err)
	}
	client.DiskInfo = strings.TrimSpace(string(out))

	// 执行 nvidia-smi 命令
	cmd = exec.Command("bash", "-c", "nvidia-smi")
	out, err = cmd.Output()
	if err != nil {
		log.Println(err)
	}
	client.Driver = strings.TrimSpace(string(out))

	// 执行 lscpu 命令
	cmd = exec.Command("bash", "-c", "lscpu | grep 'Model name'")
	out, err = cmd.Output()
	if err != nil {
		log.Println(err)
	}
	client.CPU = strings.TrimSpace(string(out))

	// 执行 lspci 命令
	cmd = exec.Command("bash", "-c", "lspci | grep NVI | grep VGA")
	out, err = cmd.Output()
	if err != nil {
		log.Println(err)
	}
	client.Card = strings.TrimSpace(string(out))

	// 执行 cat 命令
	cmd = exec.Command("bash", "-c", "cat /etc/lsb-release")
	out, err = cmd.Output()
	if err != nil {
		log.Println(err)
	}
	client.SystemInfo = strings.TrimSpace(string(out))

	// 执行 lsmem 命令
	cmd = exec.Command("bash", "-c", "lsmem")
	out, err = cmd.Output()
	if err != nil {
		log.Println(err)
	}
	client.Mem = strings.TrimSpace(string(out))

	return client
}
