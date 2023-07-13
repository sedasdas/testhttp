package db

import (
	"bytes"
	"encoding/gob"
	"testhttp/client"

	badger "github.com/dgraph-io/badger/v3"
)

// 打开数据库
func OpenDB(path string) (*badger.DB, error) {
	return badger.Open(badger.DefaultOptions(path))
}

// 关闭数据库
func CloseDB(db *badger.DB) error {
	return db.Close()
}

// 将 Client 结构体转换为字节数组
func EncodeClient(client *client.Client) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// 将 client 编码为字节序列
	err := enc.Encode(client)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 将字节数组转换为 Client 结构体
func DecodeClient(data []byte) (*client.Client, error) {
	var client client.Client

	// 将字节序列解码为 client 结构体
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&client)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

// 保存 Client 到数据库
func SaveClient(db *badger.DB, client *client.Client) error {
	data, err := EncodeClient(client)
	if err != nil {
		return err
	}
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(client.IP), data)
	})
}

// 从数据库中获取 Client
func GetClient(db *badger.DB, ip string) (*client.Client, error) {
	var result *client.Client
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(ip))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			result, err = DecodeClient(val)
			return err
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
