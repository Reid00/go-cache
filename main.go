package main

import (
	"fmt"
	"log"
	"net/http"

	"cache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "555",
}

func main() {
	cache.NewGroup("score", 1<<12, cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := cache.NewHTTPPool(addr)

	log.Println("gocache is running at: ", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
	log.Println("this is good")	
}
