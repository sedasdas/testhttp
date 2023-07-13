package main

import (
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
		log.Fatal(err)
	}
	defer db.Close()

	// 获取一个 client 实例
	client := getClientInfo()

	// 将 client 保存到数据库中
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(client.IP), encode(client))
	})
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	// 输出 client 的信息
	fmt.Printf("%+v\n", result)
}

// 将 Client 结构体转换为字节数组
func encode(client *Client) []byte {
	// 这里可以使用任意一种序列化方式，比如 json、gob、protobuf 等
	return []byte(client.IP + "\n" + client.DiskInfo + "\n" + client.Driver + "\n" + client.CPU + "\n" + client.Card + "\n" + client.SystemInfo + "\n" + client.Mem)
}

// 将字节数组转换为 Client 结构体
func decode(data []byte) Client {
	client := Client{}
	fields := strings.Split(string(data), "\n")
	client.IP = fields[0]
	client.DiskInfo = fields[1]
	client.Driver = fields[2]
	client.CPU = fields[3]
	client.Card = fields[4]
	client.SystemInfo = fields[5]
	client.Mem = fields[6]
	return client
}

// 获取 client 信息
func getClientInfo() *Client {
	client := &Client{}

	// 执行 ip 命令
	cmd := exec.Command("bash", "-c", "ip a | grep inet | grep -v inet6 | awk -F 'inet ' '{print $2}' | awk -F '/' '{print $1}' | grep 10")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	client.IP = strings.TrimSpace(string(out))

	// 执行 df 命令
	cmd = exec.Command("bash", "-c", "df -hl")
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	client.DiskInfo = strings.TrimSpace(string(out))

	// 执行 nvidia-smi 命令
	cmd = exec.Command("bash", "-c", "nvidia-smi")
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	client.Driver = strings.TrimSpace(string(out))

	// 执行 lscpu 命令
	cmd = exec.Command("bash", "-c", "lscpu | grep 'Model name'")
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	client.CPU = strings.TrimSpace(string(out))

	// 执行 lspci 命令
	cmd = exec.Command("bash", "-c", "lspci | grep NVI | grep VGA")
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	client.Card = strings.TrimSpace(string(out))

	// 执行 cat 命令
	cmd = exec.Command("bash", "-c", "cat /etc/lsb-release")
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	client.SystemInfo = strings.TrimSpace(string(out))

	// 执行 lsmem 命令
	cmd = exec.Command("bash", "-c", "lsmem")
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	client.Mem = strings.TrimSpace(string(out))

	return client
}
