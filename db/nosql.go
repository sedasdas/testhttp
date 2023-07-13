package db

import (
	"encoding/json"
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

// 保存 Client 到数据库
func SaveClient(db *badger.DB, key string, value []byte) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

// 从数据库中获取 Client
func GetClient(db *badger.DB, key string) (*client.Client, error) {
	var result *client.Client
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			err = json.Unmarshal(val, &result)
			return err
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
