package cache

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"ziwex/db"
	"ziwex/utils"
)

const cacheExpTime = time.Minute * 120

func Store(route string, metadata string, data interface{}, index *Index) {
	ctx, cancel := utils.GetRedisContext()
	defer cancel()

	d, _ := json.Marshal(data)
	key := createKey(route, metadata)
	fmt.Println("stored with key: ", key)

	//TODO: balance time out
	_ = db.Redis.SetNX(ctx, key, d, cacheExpTime)

	if index != nil {
		index.Value = key
		createIndex(index)
	}
}

func Get(route string, metadata string) ([]byte, error) {
	key := createKey(route, metadata)
	ctx, cancel := utils.GetRedisContext()
	defer cancel()
	str, err := db.Redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(str), nil
}

func createKey(route string, metadata string) string {
	return base64.StdEncoding.EncodeToString(append([]byte(route), []byte(metadata)...))
}
