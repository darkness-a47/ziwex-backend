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
	ProductIndex           = "product"
	ProductDataSubIndex    = "data"
	ProductSummerySubIndex = "summery"
)

func createIndex(index *Index) {
	ctx, cancel := utils.GetRedisContext()
	defer cancel()

	key := fmt.Sprintf("ci:%s:%s:%s", index.IndexType, index.IndexSubType, index.Index)
	_ = db.Redis.SetEX(ctx, key, index.Value, indexExpTime)

	ctx2, cancel2 := utils.GetRedisContext()
	defer cancel2()

	key2 := fmt.Sprintf("ci:%s:%s:set", index.IndexType, index.IndexSubType)
	_ = db.Redis.SAdd(ctx2, key2, key)
}

func getInvalidateIndexAll(index *Index) ([]string, error) {
	ctx, cancel := utils.GetRedisContext()
	defer cancel()

	key := fmt.Sprintf("ci:%s:%s:set", index.IndexType, index.IndexSubType)
	set, err := db.Redis.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	r := make([]string, 0)
	for _, v := range set {
		ctx, cancel := utils.GetRedisContext()
		defer cancel()

		value, err := db.Redis.GetDel(ctx, v).Result()
		if err == nil {
			r = append(r, value)
		}
	}

	ctxRem, cancelRem := utils.GetRedisContext()
	defer cancelRem()

	_ = db.Redis.SRem(ctxRem, key)

	return r, nil
}
