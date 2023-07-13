package model

import "time"

type ClientInfo struct {
	LocalIP     string    `json:"local_ip"`
	SystemInfo  string    `json:"system_info"`
	LastUpdated time.Time `json:"last_updated"`
	Status      string    `json:"status"`
}
