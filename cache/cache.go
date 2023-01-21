package cache

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"ziwex/db"
	"ziwex/types/jsonResponse"
	"ziwex/utils"
)

func createKey(route string, metadata string) string {
	return base64.StdEncoding.EncodeToString(append([]byte(route), []byte(metadata)...))
}

func Store(route string, metadata string, data interface{}) {
	ctx, cancel := utils.GetRedisContext()
	defer cancel()

	d, _ := json.Marshal(data)
	key := createKey(route, metadata)
	fmt.Println("stored with key: ", key)

	//TODO: balance time out
	_ = db.Redis.SetNX(ctx, key, d, time.Minute*120)
}

func Get(route string, metadata string) (interface{}, error) {
	key := createKey(route, metadata)
	ctx, cancel := utils.GetRedisContext()
	defer cancel()
	str, err := db.Redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	d := jsonResponse.Json{}

	jsonErr := json.Unmarshal([]byte(str), &d)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return d, nil
}
