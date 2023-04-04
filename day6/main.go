package main

import (
	"flag"
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

func createGroup() *simplecache.Group {
	return simplecache.NewGroup("score", 2<<10, simplecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, simp *simplecache.Group) {
	peers := simplecache.NewHTTPPool(addr)
	peers.Set(addrs...)
	simp.RegisterPeers(peers)
	log.Println("simplecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, simp *simplecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := simp.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("fonted server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "simpleCache server port")
	flag.BoolVar(&api, "api", false, "start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	simp := createGroup()
	if api {
		go startAPIServer(apiAddr, simp)
	}
	startCacheServer(addrMap[port], []string(addrs), simp)
}
