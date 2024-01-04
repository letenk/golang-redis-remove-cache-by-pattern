package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	// Open connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	batchSize := int64(100)
	pattern := "*:articles:*"

	// search key per batch
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = rdb.Scan(ctx, cursor, pattern, batchSize).Result()
		if err != nil {
			panic(err)
		}

		// Use pipeline to remove cache
		pipe := rdb.Pipeline()
		for _, key := range keys {
			pipe.Del(ctx, key)
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			fmt.Println("Error in pipeline execution:", err)
		}

		if cursor == 0 {
			break
		}
	}

	fmt.Println("remove cache successfully")
}
