package main

import (
	"log"
	"time"

	"github.com/allegro/bigcache/v2"
)

func main() {

	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		return
	}
	entry, err := cache.Get("my-unique-key")
	if err != nil {
		return
	}
	log.Println(entry)
}
