package cache

import (
	"fmt"
	"time"
	"ziwex/db"
	"ziwex/utils"
)

const indexExpTime = cacheExpTime + (time.Second * 10)

type Index struct {
	IndexType    string
	IndexSubType string
	Index        string
	Value        string
}

const (
	ProductIndex        = "product"
	ProductDataSubIndex = "data"
)

func createIndex(index *Index) {
	ctx, cancel := utils.GetRedisContext()
	defer cancel()

	key := fmt.Sprintf("c:%s:%s:%s", index.IndexType, index.IndexSubType, index.Index)
	_ = db.Redis.Set(ctx, key, index.Value, indexExpTime)
}
