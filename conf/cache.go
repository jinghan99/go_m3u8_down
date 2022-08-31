package conf

import (
	"encoding/json"
	"github.com/xujiajun/nutsdb"
)

//Cache 内存
//var Cache *cache.Cache //定义全局变量
//Ndb 持久化内存 db
var Ndb *nutsdb.DB

func init() {
	// 设置超时时间和清理时间
	//Cache = cache.New(5*time.Minute, 60*time.Second)
	initNdb()
}
func initNdb() {
	opt := nutsdb.DefaultOptions
	opt.Dir = "./resource/nutsdb"
	Ndb, _ = nutsdb.Open(opt)
}
func NdbPut(key string, value []byte) error {
	err := Ndb.Update(func(tx *nutsdb.Tx) error {
		//过期时间 为7天
		expiration := uint32(7 * 27 * 60 * 60)
		if err := tx.Put("", []byte(key), value, expiration); err != nil {
			return err
		}
		return nil
	})
	return err
}

func NdbDel(key string) {
	Ndb.Update(func(tx *nutsdb.Tx) error {
		err := tx.Delete("", []byte(key))
		return err
	})
}

func NdbGet(key string) ([]byte, error) {
	var value *nutsdb.Entry
	err := Ndb.View(func(tx *nutsdb.Tx) error {
		get, err := tx.Get("", []byte(key))
		if err == nil {
			value = get
			return nil
		}
		return err
	})
	return value.Value, err
}

func NdbPutAny(key string, v any) error {
	//先转json
	jsonByte, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = Ndb.Update(func(tx *nutsdb.Tx) error {
		//过期时间 为7天
		expiration := uint32(7 * 27 * 60 * 60)
		if err := tx.Put("", []byte(key), jsonByte, expiration); err != nil {
			return err
		}
		return nil
	})
	return err
}

func NdbGetAny(key string, v any) error {
	err := Ndb.Update(func(tx *nutsdb.Tx) error {
		get, err := tx.Get("", []byte(key))
		if err == nil {
			err = json.Unmarshal(get.Value, v)
			return err
		}
		return nil
	})
	return err
}
