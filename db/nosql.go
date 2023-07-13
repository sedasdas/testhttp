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
func SaveClient(db *badger.DB, key string, value client.Client) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), jsonData)
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

// 获取所有 Client
func GetAllClients(db *badger.DB) ([]client.Client, error) {
	clients := make([]client.Client, 0)
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek([]byte("")); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			client, err := GetClient(db, string(key))
			if err != nil {
				return err
			}
			clients = append(clients, *client)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return clients, nil
}

// 更新 Client
func UpdateClient(db *badger.DB, key string, value client.Client) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), jsonData)
	})
}

// 删除 Client
func DeleteClient(db *badger.DB, key string) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}
