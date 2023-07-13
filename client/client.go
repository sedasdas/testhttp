package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
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

func StartClient() {

	// 每隔一分钟发送客户端信息
	tick := time.Tick(time.Minute)
	for range tick {
		client, _ := GetClientInfo()
		// 将客户端信息发送到服务端的新增 client 的 API
		err := addClient(&client)
		if err != nil {
			log.Printf("Failed to add client: %v", err)
		}
	}
}

func addClient(client *Client) error {
	// 将客户端信息转换为 JSON 格式
	clientJSON, err := json.Marshal(client)
	if err != nil {
		return err
	}

	// 发送 HTTP POST 请求到服务端的新增 client 的 API
	resp, err := http.Post("http://localhost:8080/client", "application/json", bytes.NewBuffer(clientJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to add client: %s", resp.Status)
	}

	return nil
}

// 获取 client 信息
func GetClientInfo() (Client, error) {
	client := Client{}

	// 执行 ip 命令
	cmd := exec.Command("bash", "-c", "ip a | grep inet | grep -v inet6 | awk -F 'inet ' '{print $2}' | awk -F '/' '{print $1}' | grep 10")
	out, err := cmd.Output()

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
	cmd = exec.Command("bash", "-c", "lspci | grepNVidia | grep VGA")
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

	return client, nil
}
