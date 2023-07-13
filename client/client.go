package client

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
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
	client := &Client{}

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

	// 输出结果
	fmt.Printf("%+v\n", client)

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
