package main

import (
	"fmt"
	"log"
	"net/http"
	"simplecache/simplecache"
)

var db = map[string]string{
	"Tom": "639",
	"zz":  "22",
	"Jay": "567",
}

func main() {
	simplecache.NewGroup("score", 2<<10, simplecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := simplecache.NewHTTPPool(addr)
	log.Println("simplecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
